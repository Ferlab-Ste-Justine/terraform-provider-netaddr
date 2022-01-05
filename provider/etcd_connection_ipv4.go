package provider

import (
	"context"
	"errors"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Ipv4EtcdKeyPrefixes struct {
	DeletedAddress string
	HardcodedAddress string
	GeneratedAddress string
	Name string
}

func GenerateIpv4EtcdKeyPrefixes(networkPrefix string) Ipv4EtcdKeyPrefixes {
	return Ipv4EtcdKeyPrefixes{
		DeletedAddress: networkPrefix + "data/ipv4/address/deleted/",
		HardcodedAddress: networkPrefix + "data/ipv4/address/hardcoded/",
		GeneratedAddress: networkPrefix + "data/ipv4/address/generated/",
		Name: networkPrefix + "data/ipv4/name/",
	}
}

func (conn *EtcdConnection) getNextIpv4Address(prefix string) ([]byte, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conn.Timeout)*time.Second)
	defer cancel()

	netKeys := GenerateNetworkEtcdKeys(prefix)
	getRes, err := conn.Client.Get(ctx, netKeys.NextIpv4)

	if err != nil {
		return []byte{}, 0, err
	}

	if len(getRes.Kvs) == 0 {
		return []byte{}, 0, errors.New(fmt.Sprintf("Error accessing next ipv4 for network with prefix '%s': Key not found", prefix))
	}

	return getRes.Kvs[0].Value, getRes.Kvs[0].Version, nil
}

func (conn *EtcdConnection) getDeletedIpv4Address(prefix string) ([]byte, bool, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conn.Timeout)*time.Second)
	defer cancel()

	ipv4KeyPrefixes := GenerateIpv4EtcdKeyPrefixes(prefix)
	
	getRes, err := conn.Client.Get(ctx, ipv4KeyPrefixes.DeletedAddress, clientv3.WithPrefix())
	if err != nil {
		return []byte{}, false, 0, err
	}

	if len(getRes.Kvs) == 0 {
		return []byte{}, false, 0, nil
	}

	return getRes.Kvs[0].Key, true, getRes.Kvs[0].Version, nil
}

func (conn *EtcdConnection) addressIsHardcoded(prefix string, address []byte) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conn.Timeout)*time.Second)
	defer cancel()

	ipv4KeyPrefixes := GenerateIpv4EtcdKeyPrefixes(prefix)
	
	getRes, err := conn.Client.Get(ctx, ipv4KeyPrefixes.HardcodedAddress + string(address))
	if err != nil {
		return false, err
	}

	return len(getRes.Kvs) > 0, nil
}

/*
  check before transaction:
    - Ip is within the range
  check during transaction:
    - address doesn't exist in hardcoded/
	- address doesn't exist in generated/
	- name doesn't exist in names/
  transaction:
    - Insert address in hardcoded/
	- Insert name in names/
*/
func (conn *EtcdConnection) createHardcodedIpv4AddressWithRetries(prefix string, name string, address []byte, retries int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conn.Timeout)*time.Second)
	defer cancel()

	network, netExists, err := conn.getNetworkInfoWithRetries(prefix, 0)
	if err != nil {
		if !shouldRetry(err, retries) {
			return err
		}

		time.Sleep(100 * time.Millisecond)
		return conn.createHardcodedIpv4AddressWithRetries(prefix, name, address, retries - 1)
	}
	if !netExists {
		return errors.New(fmt.Sprintf("Error created hardcoded ipv4 '%s': Network not found", Ipv4BytesToString(address)))
	}

	if !AddressWithinBoundaries(address, network.Ipv4.First, network.Ipv4.Last) {
		return errors.New(fmt.Sprintf("Error created hardcoded ipv4 '%s': Ip is outside of address boundaries", Ipv4BytesToString(address)))
	}

	ipv4KeyPrefixes := GenerateIpv4EtcdKeyPrefixes(prefix)
	tx := conn.Client.Txn(ctx).If(
		clientv3.Compare(clientv3.Version(ipv4KeyPrefixes.HardcodedAddress + string(address)), "=", 0),
		clientv3.Compare(clientv3.Version(ipv4KeyPrefixes.GeneratedAddress + string(address)), "=", 0),
		clientv3.Compare(clientv3.Version(ipv4KeyPrefixes.Name + name), "=", 0),
	).Then(
		clientv3.OpPut(ipv4KeyPrefixes.HardcodedAddress + string(address), name),
		clientv3.OpPut(ipv4KeyPrefixes.Name + name, string(address)),
	)

	resp, txErr := tx.Commit()
	if txErr != nil {
		if !shouldRetry(txErr, retries) {
			return txErr
		}

		time.Sleep(100 * time.Millisecond)
		return conn.createHardcodedIpv4AddressWithRetries(prefix, name, address, retries - 1)
	}

	if !resp.Succeeded {
		return errors.New(fmt.Sprintf("Failed to create hardcoded ipv4 '%s': Either address or name is already in use", Ipv4BytesToString(address)))
	}

	return nil
}

func (conn *EtcdConnection) CreateHardcodedIpv4Address(prefix string, name string, address []byte) error {
	return conn.createHardcodedIpv4AddressWithRetries(prefix, name, address, conn.Retries)
}

/*
  if address greater than or equal to next address:
    check during transaction:
      - next address version has not changed
	  - address exists in hardcoded/
	  - name exists in name/
    transaction:
      - delete address from hardcoded/
	  - delete name from name/
  if address less than next address:\
    check during transaction:
	  - address does not exist in deleted/
	  - address exists in hardcoded/
	  - name exists in name/
    transaction:
      - delete address from hardcoded/
	  - delete name from name/
	  - add address to deleted/
*/
func (conn *EtcdConnection) deleteHardcodedIpv4AddressWithRetries(prefix string, name string, address []byte, retries int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conn.Timeout)*time.Second)
	defer cancel()

	ipv4KeyPrefixes := GenerateIpv4EtcdKeyPrefixes(prefix)

	nextAddr, nextAddrVer, err := conn.getNextIpv4Address(prefix)
	if err != nil {
		if !shouldRetry(err, retries) {
			return err
		}

		time.Sleep(100 * time.Millisecond)
		return conn.deleteHardcodedIpv4AddressWithRetries(prefix, name, address, retries - 1)
	}

	if !AddressLessThan(address, nextAddr) {
		netKeys := GenerateNetworkEtcdKeys(prefix)
		tx := conn.Client.Txn(ctx).If(
			clientv3.Compare(clientv3.Version(netKeys.NextIpv4), "=", nextAddrVer),
			clientv3.Compare(clientv3.Version(ipv4KeyPrefixes.HardcodedAddress + string(address)), ">", 0),
			clientv3.Compare(clientv3.Version(ipv4KeyPrefixes.Name + name), ">", 0),
		).Then(
			clientv3.OpDelete(ipv4KeyPrefixes.HardcodedAddress + string(address)),
			clientv3.OpDelete(ipv4KeyPrefixes.Name + name),
		)
	
		resp, txErr := tx.Commit()
		if txErr != nil {
			if !shouldRetry(txErr, retries) {
				return txErr
			}
	
			time.Sleep(100 * time.Millisecond)
			return conn.deleteHardcodedIpv4AddressWithRetries(prefix, name, address, retries - 1)
		}
	
		if !resp.Succeeded {
			if retries <= 0 {
				return errors.New(fmt.Sprintf("Failed to delete hardcoded ipv4 '%s': Address or name have not been assigned.", Ipv4BytesToString(address)))
			}
			return conn.deleteHardcodedIpv4AddressWithRetries(prefix, name, address, retries - 1)
		}
	}

	tx := conn.Client.Txn(ctx).If(
		clientv3.Compare(clientv3.Version(ipv4KeyPrefixes.DeletedAddress + string(address)), "=", 0),
		clientv3.Compare(clientv3.Version(ipv4KeyPrefixes.HardcodedAddress + string(address)), ">", 0),
		clientv3.Compare(clientv3.Version(ipv4KeyPrefixes.Name + name), ">", 0),
	).Then(
		clientv3.OpDelete(ipv4KeyPrefixes.HardcodedAddress + string(address)),
		clientv3.OpDelete(ipv4KeyPrefixes.Name + name),
		clientv3.OpPut(ipv4KeyPrefixes.DeletedAddress  + string(address), name),
	)

	resp, txErr := tx.Commit()
	if txErr != nil {
		if !shouldRetry(txErr, retries) {
			return txErr
		}

		time.Sleep(100 * time.Millisecond)
		return conn.deleteHardcodedIpv4AddressWithRetries(prefix, name, address, retries - 1)
	}

	if !resp.Succeeded {
		return errors.New(fmt.Sprintf("Failed to delete hardcoded ipv4 '%s': Either address or name have not been assigned or address was already deleted", Ipv4BytesToString(address)))
	}

	return nil
}

func (conn *EtcdConnection) DeleteHardcodedIpv4Address(prefix string, name string, address []byte) error {
	return conn.deleteHardcodedIpv4AddressWithRetries(prefix, name, address, conn.Retries)
}

/* 
	if deleted/ has addresses:
	  get an address from deleted/
	  check during transaction:
	    - picked address is present in deleted/
		- name is absent from name/
	  transaction:
	    - Remove picked address from deleted/
		- Add picked address to generated/
		- Add name to name/
	if deleted/ has no address:
	  get next assignable address
	  increment next address until has address not present in hardcoded/ is found
	  check during transaction:
	    - next assignable address has the same version
		- picked address is absent from hardcoded/
		- name is absent from name/
	  transaction:
	    - add picked address to generated/
		- set next assignable address to picked address + 1
		- add name to name/
*/
func (conn *EtcdConnection) createGeneratedIpv4AddressWithRetries(prefix string, name string, retries int) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conn.Timeout)*time.Second)
	defer cancel()

	ipv4KeyPrefixes := GenerateIpv4EtcdKeyPrefixes(prefix)
	netKeys := GenerateNetworkEtcdKeys(prefix)

	deletedAddr, deletedAddrExists, _, deletedAddrErr := conn.getDeletedIpv4Address(prefix)
	if deletedAddrErr != nil {
		if !shouldRetry(deletedAddrErr, retries) {
			return []byte{}, deletedAddrErr
		}

		time.Sleep(100 * time.Millisecond)
		return conn.createGeneratedIpv4AddressWithRetries(prefix, name, retries - 1)
	}

	if deletedAddrExists {
		tx := conn.Client.Txn(ctx).If(
			clientv3.Compare(clientv3.Version(ipv4KeyPrefixes.DeletedAddress + string(deletedAddr)), ">", 0),
			clientv3.Compare(clientv3.Version(ipv4KeyPrefixes.Name + name), "=", 0),
		).Then(
			clientv3.OpDelete(ipv4KeyPrefixes.DeletedAddress + string(deletedAddr)),
			clientv3.OpPut(ipv4KeyPrefixes.GeneratedAddress  + string(deletedAddr), name),
			clientv3.OpPut(ipv4KeyPrefixes.Name + name, string(deletedAddr)),
		)
	
		resp, txErr := tx.Commit()
		if txErr != nil {
			if !shouldRetry(txErr, retries) {
				return []byte{}, txErr
			}
	
			time.Sleep(100 * time.Millisecond)
			return conn.createGeneratedIpv4AddressWithRetries(prefix, name, retries - 1)
		}
	
		if !resp.Succeeded {
			if retries <= 0 {
				return []byte{}, errors.New("Failed to create generated ipv4: Selected name has already been assigned")
			}

			return conn.createGeneratedIpv4AddressWithRetries(prefix, name, retries - 1)
		}
		
		return deletedAddr, nil
	}

	network, networkExists, networkErr := conn.getNetworkInfoWithRetries(prefix, 0)
	if networkErr != nil {
		if !shouldRetry(networkErr, retries) {
			return []byte{}, networkErr
		}

		time.Sleep(100 * time.Millisecond)
		return conn.createGeneratedIpv4AddressWithRetries(prefix, name, retries - 1)
	}
	if !networkExists {
		return []byte{}, errors.New("Error creating generated ipv4: Network does not exist")
	}

	nextAddr, nextAddrVer, nextAddrErr := conn.getNextIpv4Address(prefix)
	if nextAddrErr != nil {
		if !shouldRetry(nextAddrErr, retries) {
			return []byte{}, nextAddrErr
		}

		time.Sleep(100 * time.Millisecond)
		return conn.createGeneratedIpv4AddressWithRetries(prefix, name, retries - 1)
	}

	if AddressGreaterThan(nextAddr, network.Ipv4.Last) {
		return []byte{}, errors.New("Error creating generated ipv4: Network ran out of addresses")
	}

	isHardcoded, isHarcodedErr := conn.addressIsHardcoded(prefix, nextAddr)
	if isHarcodedErr != nil {
		if !shouldRetry(isHarcodedErr, retries) {
			return []byte{}, isHarcodedErr
		}

		time.Sleep(100 * time.Millisecond)
		return conn.createGeneratedIpv4AddressWithRetries(prefix, name, retries - 1)
	}

	for isHardcoded {
		nextAddr := IncAddressBy1(nextAddr)

		if AddressGreaterThan(nextAddr, network.Ipv4.Last) {
			return []byte{}, errors.New("Error creating generated ipv4: Network ran out of addresses")
		}

		isHardcoded, isHarcodedErr = conn.addressIsHardcoded(prefix, nextAddr)
		if isHarcodedErr != nil {
			if !shouldRetry(isHarcodedErr, retries) {
				return []byte{}, isHarcodedErr
			}
	
			time.Sleep(100 * time.Millisecond)
			return conn.createGeneratedIpv4AddressWithRetries(prefix, name, retries - 1)
		}
	}

	tx := conn.Client.Txn(ctx).If(
		clientv3.Compare(clientv3.Version(netKeys.NextIpv4), "=", nextAddrVer),
		clientv3.Compare(clientv3.Version(ipv4KeyPrefixes.HardcodedAddress + string(nextAddr)), "=", 0),
		clientv3.Compare(clientv3.Version(ipv4KeyPrefixes.Name + name), "=", 0),
	).Then(
		clientv3.OpDelete(ipv4KeyPrefixes.GeneratedAddress + string(nextAddr)),
		clientv3.OpPut(netKeys.NextIpv4, string(IncAddressBy1(nextAddr))),
		clientv3.OpPut(ipv4KeyPrefixes.Name + name, string(nextAddr)),
	)

	resp, txErr := tx.Commit()
	if txErr != nil {
		if !shouldRetry(txErr, retries) {
			return []byte{}, txErr
		}

		time.Sleep(100 * time.Millisecond)
		return conn.createGeneratedIpv4AddressWithRetries(prefix, name, retries - 1)
	}

	if !resp.Succeeded {
		if retries <= 0 {
			return []byte{}, errors.New("Failed to create generated ipv4: Selected name has already been assigned")
		}

		return conn.createGeneratedIpv4AddressWithRetries(prefix, name, retries - 1)
	}

	return nextAddr, nil
}

func (conn *EtcdConnection) CreateGeneratedIpv4Address(prefix string, name string) ([]byte, error) {
	return conn.createGeneratedIpv4AddressWithRetries(prefix, name, conn.Retries)
}

/* 
	check during transaction:
	  - address is present in generated/
	  - name is present in name/
	transaction:
	  - remote address from generated/
	  - remove name from name/ 
	  - add address to deleted/
*/
func (conn *EtcdConnection) deleteGeneratedIpv4AddressWithRetries(prefix string, name string, address []byte, retries int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conn.Timeout)*time.Second)
	defer cancel()

	ipv4KeyPrefixes := GenerateIpv4EtcdKeyPrefixes(prefix)

	tx := conn.Client.Txn(ctx).If(
		clientv3.Compare(clientv3.Version(ipv4KeyPrefixes.GeneratedAddress + string(address)), ">", 0),
		clientv3.Compare(clientv3.Version(ipv4KeyPrefixes.Name + name), ">", 0),
	).Then(
		clientv3.OpDelete(ipv4KeyPrefixes.GeneratedAddress + string(address)),
		clientv3.OpDelete(ipv4KeyPrefixes.Name + name),
		clientv3.OpPut(ipv4KeyPrefixes.DeletedAddress  + string(address), name),
	)

	resp, txErr := tx.Commit()
	if txErr != nil {
		if !shouldRetry(txErr, retries) {
			return txErr
		}

		time.Sleep(100 * time.Millisecond)
		return conn.deleteGeneratedIpv4AddressWithRetries(prefix, name, address, retries - 1)
	}

	if !resp.Succeeded {
		return errors.New(fmt.Sprintf("Failed to delete generated ipv4 '%s': Either address or name have not been assigned or address was already deleted", Ipv4BytesToString(address)))
	}

	return nil
}

func (conn *EtcdConnection) DeleteGeneratedIpv4Address(prefix string, name string, address []byte) error {
	return conn.deleteGeneratedIpv4AddressWithRetries(prefix, name, address, conn.Retries)
}

func (conn *EtcdConnection) getIpv4AddressWithRetries(prefix string, name string, retries int) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conn.Timeout)*time.Second)
	defer cancel()

	ipv4KeyPrefixes := GenerateIpv4EtcdKeyPrefixes(prefix)

	getRes, err := conn.Client.Get(ctx, ipv4KeyPrefixes.Name + name)
	if err != nil {
		if !shouldRetry(err, retries) {
			return []byte{}, err
		}

		time.Sleep(100 * time.Millisecond)
		return conn.getIpv4AddressWithRetries(prefix, name, retries - 1)
	}

	if len(getRes.Kvs) == 0 {
		return []byte{}, errors.New(fmt.Sprintf("Error retrieving address of ipv4 '%s': Name not found", name))
	}

	return getRes.Kvs[0].Value, nil
}

func (conn *EtcdConnection) GetIpv4Address(prefix string, name string) ([]byte, error) {
	return conn.getIpv4AddressWithRetries(prefix, name, conn.Retries)
}