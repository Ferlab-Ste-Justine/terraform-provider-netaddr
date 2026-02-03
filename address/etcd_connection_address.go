package address

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"
	"slices"
	"strings"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type ParseAddr func(string) ([]byte, error)
type PrettifyAddr func([]byte) string
type AddressIsGreater func([]byte, []byte) bool
type AddressIsLess func([]byte, []byte) bool
type IncrementAddress func([]byte) []byte

type AddrEtcdKeyPrefixes struct {
	DeletedAddress string
	HardcodedAddress string
	GeneratedAddress string
	Name string
}

func GenerateAddrEtcdKeyPrefixes(rangePrefix string) AddrEtcdKeyPrefixes {
	return AddrEtcdKeyPrefixes{
		DeletedAddress: rangePrefix + "data/address/deleted/",
		HardcodedAddress: rangePrefix + "data/address/hardcoded/",
		GeneratedAddress: rangePrefix + "data/address/generated/",
		Name: rangePrefix + "data/name/",
	}
}

func (conn *EtcdConnection) getNextAddress(prefix string) ([]byte, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conn.Timeout)*time.Second)
	defer cancel()

	addrRangeKeys := GenerateAddrRangeEtcdKeys(prefix)
	getRes, err := conn.Client.Get(ctx, addrRangeKeys.NextAddress)

	if err != nil {
		return []byte{}, 0, err
	}

	if len(getRes.Kvs) == 0 {
		return []byte{}, 0, errors.New(fmt.Sprintf("Error accessing next address for range with prefix '%s': Key not found", prefix))
	}

	return getRes.Kvs[0].Value, getRes.Kvs[0].Version, nil
}

func (conn *EtcdConnection) getDeletedAddress(prefix string) ([]byte, bool, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conn.Timeout)*time.Second)
	defer cancel()

	addrKeyPrefixes := GenerateAddrEtcdKeyPrefixes(prefix)
	
	getRes, err := conn.Client.Get(ctx, addrKeyPrefixes.DeletedAddress, clientv3.WithPrefix())
	if err != nil {
		return []byte{}, false, 0, err
	}

	if len(getRes.Kvs) == 0 {
		return []byte{}, false, 0, nil
	}

	return bytes.TrimPrefix(getRes.Kvs[0].Key, []byte(addrKeyPrefixes.DeletedAddress)), true, getRes.Kvs[0].Version, nil
}

func (conn *EtcdConnection) addressIsHardcoded(prefix string, address []byte) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conn.Timeout)*time.Second)
	defer cancel()

	addrKeyPrefixes := GenerateAddrEtcdKeyPrefixes(prefix)
	
	getRes, err := conn.Client.Get(ctx, addrKeyPrefixes.HardcodedAddress + string(address))
	if err != nil {
		return false, err
	}

	return len(getRes.Kvs) > 0, nil
}

func (conn *EtcdConnection) addressIsDeleted(prefix string, address []byte) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conn.Timeout)*time.Second)
	defer cancel()

	addrKeyPrefixes := GenerateAddrEtcdKeyPrefixes(prefix)
	
	getRes, err := conn.Client.Get(ctx, addrKeyPrefixes.DeletedAddress + string(address))
	if err != nil {
		return false, err
	}

	return len(getRes.Kvs) > 0, nil
}

type AddressListEntry struct {
	Name string
	Address []byte
}

func (conn *EtcdConnection) getAddressListWithRetries(prefix string, retries int) ([]AddressListEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conn.Timeout)*time.Second)
	defer cancel()

	addrKeyPrefixes := GenerateAddrEtcdKeyPrefixes(prefix)

	getRes, err := conn.Client.Get(ctx, addrKeyPrefixes.Name, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
	if err != nil {
		if !shouldRetry(err, retries) {
			return []AddressListEntry{}, err
		}

		time.Sleep(100 * time.Millisecond)
		return conn.getAddressListWithRetries(prefix, retries - 1)
	}

	listing := make([]AddressListEntry, len(getRes.Kvs))
	for idx, val := range getRes.Kvs {
		listing[idx] = AddressListEntry{strings.TrimPrefix(string(val.Key), addrKeyPrefixes.Name), val.Value}
	}

	return listing, nil
}

func (conn *EtcdConnection) GetAddressList(prefix string) ([]AddressListEntry, error) {
	return conn.getAddressListWithRetries(prefix, conn.Retries)
}

/*
  check before transaction:
    - address is within the range
  check during transaction:
    - address doesn't exist in hardcoded/
	- address doesn't exist in generated/
	- name doesn't exist in names/
  transaction:
    - Insert address in hardcoded/
	- Insert name in names/
*/
func (conn *EtcdConnection) createHardcodedAddressWithRetries(prefix string, name string, address []byte, prettify PrettifyAddr, retries int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conn.Timeout)*time.Second)
	defer cancel()

	addrRange, addrRangeExists, err := conn.getAddrRangeWithRetries(prefix, 0)
	if err != nil {
		if !shouldRetry(err, retries) {
			return err
		}

		time.Sleep(100 * time.Millisecond)
		return conn.createHardcodedAddressWithRetries(prefix, name, address, prettify, retries - 1)
	}
	if !addrRangeExists {
		return errors.New(fmt.Sprintf("Error created hardcoded address '%s': Range not found", prettify(address)))
	}

	if !AddressWithinBoundaries(address, addrRange.FirstAddress, addrRange.LastAddress) {
		return errors.New(fmt.Sprintf("Error created hardcoded address '%s': Ip is outside of range boundaries", prettify(address)))
	}

	isDeleted, isDeletedErr := conn.addressIsDeleted(prefix, address)
	if isDeletedErr != nil {
		if !shouldRetry(isDeletedErr, retries) {
			return isDeletedErr
		}

		time.Sleep(100 * time.Millisecond)
		return conn.createHardcodedAddressWithRetries(prefix, name, address, prettify, retries - 1)
	}

	addrKeyPrefixes := GenerateAddrEtcdKeyPrefixes(prefix)

	if isDeleted {
		tx := conn.Client.Txn(ctx).If(
			clientv3.Compare(clientv3.Version(addrKeyPrefixes.DeletedAddress + string(address)), ">", 0),
			clientv3.Compare(clientv3.Version(addrKeyPrefixes.HardcodedAddress + string(address)), "=", 0),
			clientv3.Compare(clientv3.Version(addrKeyPrefixes.GeneratedAddress + string(address)), "=", 0),
			clientv3.Compare(clientv3.Version(addrKeyPrefixes.Name + name), "=", 0),
		).Then(
			clientv3.OpDelete(addrKeyPrefixes.DeletedAddress + string(address)),
			clientv3.OpPut(addrKeyPrefixes.HardcodedAddress + string(address), name),
			clientv3.OpPut(addrKeyPrefixes.Name + name, string(address)),
		)

		resp, txErr := tx.Commit()
		if txErr != nil {
			if !shouldRetry(txErr, retries) {
				return txErr
			}

			time.Sleep(100 * time.Millisecond)
			return conn.createHardcodedAddressWithRetries(prefix, name, address, prettify, retries - 1)
		}

		if !resp.Succeeded {
			return errors.New(fmt.Sprintf("Failed to create hardcoded address '%s': Either address or name is already in use", prettify(address)))
		}

		return nil
	}

	
	tx := conn.Client.Txn(ctx).If(
		clientv3.Compare(clientv3.Version(addrKeyPrefixes.HardcodedAddress + string(address)), "=", 0),
		clientv3.Compare(clientv3.Version(addrKeyPrefixes.GeneratedAddress + string(address)), "=", 0),
		clientv3.Compare(clientv3.Version(addrKeyPrefixes.Name + name), "=", 0),
	).Then(
		clientv3.OpPut(addrKeyPrefixes.HardcodedAddress + string(address), name),
		clientv3.OpPut(addrKeyPrefixes.Name + name, string(address)),
	)

	resp, txErr := tx.Commit()
	if txErr != nil {
		if !shouldRetry(txErr, retries) {
			return txErr
		}

		time.Sleep(100 * time.Millisecond)
		return conn.createHardcodedAddressWithRetries(prefix, name, address, prettify, retries - 1)
	}

	if !resp.Succeeded {
		return errors.New(fmt.Sprintf("Failed to create hardcoded address '%s': Either address or name is already in use", prettify(address)))
	}

	return nil
}

func (conn *EtcdConnection) CreateHardcodedAddress(prefix string, name string, address []byte, prettify PrettifyAddr) error {
	return conn.createHardcodedAddressWithRetries(prefix, name, address, prettify, conn.Retries)
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
func (conn *EtcdConnection) deleteHardcodedAddressWithRetries(prefix string, name string, address []byte, prettify PrettifyAddr, addrIsLess AddressIsLess, retries int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conn.Timeout)*time.Second)
	defer cancel()

	addrKeyPrefixes := GenerateAddrEtcdKeyPrefixes(prefix)

	nextAddr, nextAddrVer, err := conn.getNextAddress(prefix)
	if err != nil {
		if !shouldRetry(err, retries) {
			return err
		}

		time.Sleep(100 * time.Millisecond)
		return conn.deleteHardcodedAddressWithRetries(prefix, name, address, prettify, addrIsLess, retries - 1)
	}

	if !addrIsLess(address, nextAddr) {
		addrRangeKeys := GenerateAddrRangeEtcdKeys(prefix)
		tx := conn.Client.Txn(ctx).If(
			clientv3.Compare(clientv3.Version(addrRangeKeys.NextAddress), "=", nextAddrVer),
			clientv3.Compare(clientv3.Version(addrKeyPrefixes.HardcodedAddress + string(address)), ">", 0),
			clientv3.Compare(clientv3.Version(addrKeyPrefixes.Name + name), ">", 0),
		).Then(
			clientv3.OpDelete(addrKeyPrefixes.HardcodedAddress + string(address)),
			clientv3.OpDelete(addrKeyPrefixes.Name + name),
		)
	
		resp, txErr := tx.Commit()
		if txErr != nil {
			if !shouldRetry(txErr, retries) {
				return txErr
			}
	
			time.Sleep(100 * time.Millisecond)
			return conn.deleteHardcodedAddressWithRetries(prefix, name, address, prettify, addrIsLess, retries - 1)
		}
	
		if !resp.Succeeded {
			if retries <= 0 {
				return errors.New(fmt.Sprintf("Failed to delete hardcoded address '%s': Address or name have not been assigned.", prettify(address)))
			}
			return conn.deleteHardcodedAddressWithRetries(prefix, name, address, prettify, addrIsLess, retries - 1)
		}

		return nil
	}

	tx := conn.Client.Txn(ctx).If(
		clientv3.Compare(clientv3.Version(addrKeyPrefixes.DeletedAddress + string(address)), "=", 0),
		clientv3.Compare(clientv3.Version(addrKeyPrefixes.HardcodedAddress + string(address)), ">", 0),
		clientv3.Compare(clientv3.Version(addrKeyPrefixes.Name + name), ">", 0),
	).Then(
		clientv3.OpDelete(addrKeyPrefixes.HardcodedAddress + string(address)),
		clientv3.OpDelete(addrKeyPrefixes.Name + name),
		clientv3.OpPut(addrKeyPrefixes.DeletedAddress  + string(address), name),
	)

	resp, txErr := tx.Commit()
	if txErr != nil {
		if !shouldRetry(txErr, retries) {
			return txErr
		}

		time.Sleep(100 * time.Millisecond)
		return conn.deleteHardcodedAddressWithRetries(prefix, name, address, prettify, addrIsLess, retries - 1)
	}

	if !resp.Succeeded {
		return errors.New(fmt.Sprintf("Failed to delete hardcoded address '%s': Either address or name have not been assigned or address was already deleted", prettify(address)))
	}

	return nil
}

func (conn *EtcdConnection) DeleteHardcodedAddress(prefix string, name string, address []byte, prettify PrettifyAddr, addrIsLess AddressIsLess) error {
	return conn.deleteHardcodedAddressWithRetries(prefix, name, address, prettify, addrIsLess, conn.Retries)
}

/* 
	if deleted/ has addresses:
	  get an address from deleted/
	  check during transaction:
	    - picked address is present in deleted/
		- name is absent from name/ for all relevant prefixes
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
		- name is absent from name/ for all relevant prefixes
	  transaction:
	    - add picked address to generated/
		- set next assignable address to picked address + 1
		- add name to name/
*/
func (conn *EtcdConnection) createGeneratedAddressWithRetries(prefix string, mutExclPrefixes []string, name string, addrIsGreater AddressIsGreater, incAddr IncrementAddress, retries int) ([]byte, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conn.Timeout)*time.Second)
	defer cancel()

	addrKeyPrefixes := GenerateAddrEtcdKeyPrefixes(prefix)
	addrRangeKeys := GenerateAddrRangeEtcdKeys(prefix)

	nameNoPresent := []clientv3.Cmp{}
	for _, mutExclPrefix := range mutExclPrefixes{
		addrKeyMutExclPrefixes := GenerateAddrEtcdKeyPrefixes(mutExclPrefix)
		nameNoPresent = append(nameNoPresent, clientv3.Compare(clientv3.Version(addrKeyMutExclPrefixes.Name + name), "=", 0))
	}

	deletedAddr, deletedAddrExists, _, deletedAddrErr := conn.getDeletedAddress(prefix)
	if deletedAddrErr != nil {
		if !shouldRetry(deletedAddrErr, retries) {
			return []byte{}, false, deletedAddrErr
		}

		time.Sleep(100 * time.Millisecond)
		return conn.createGeneratedAddressWithRetries(prefix, mutExclPrefixes, name, addrIsGreater, incAddr, retries - 1)
	}

	if deletedAddrExists {
		tx := conn.Client.Txn(ctx).If(
			slices.Concat(
				[]clientv3.Cmp{
					clientv3.Compare(clientv3.Version(addrKeyPrefixes.DeletedAddress + string(deletedAddr)), ">", 0),
				},
				nameNoPresent,
			)...
		).Then(
			clientv3.OpDelete(addrKeyPrefixes.DeletedAddress + string(deletedAddr)),
			clientv3.OpPut(addrKeyPrefixes.GeneratedAddress  + string(deletedAddr), name),
			clientv3.OpPut(addrKeyPrefixes.Name + name, string(deletedAddr)),
		)
	
		resp, txErr := tx.Commit()
		if txErr != nil {
			if !shouldRetry(txErr, retries) {
				return []byte{}, false, txErr
			}
	
			time.Sleep(100 * time.Millisecond)
			return conn.createGeneratedAddressWithRetries(prefix, mutExclPrefixes, name, addrIsGreater, incAddr, retries - 1)
		}
	
		if !resp.Succeeded {
			if retries <= 0 {
				return []byte{}, false, errors.New("Failed to create generated address: Selected name has already been assigned")
			}

			return conn.createGeneratedAddressWithRetries(prefix, mutExclPrefixes, name, addrIsGreater, incAddr, retries - 1)
		}
		
		return deletedAddr, false, nil
	}

	addrRange, addrRangeExists, addrRangeErr := conn.getAddrRangeWithRetries(prefix, 0)
	if addrRangeErr != nil {
		if !shouldRetry(addrRangeErr, retries) {
			return []byte{}, false, addrRangeErr
		}

		time.Sleep(100 * time.Millisecond)
		return conn.createGeneratedAddressWithRetries(prefix, mutExclPrefixes, name, addrIsGreater, incAddr, retries - 1)
	}
	if !addrRangeExists {
		return []byte{}, false, errors.New("Error creating generated address: Range does not exist")
	}

	nextAddr, nextAddrVer, nextAddrErr := conn.getNextAddress(prefix)
	if nextAddrErr != nil {
		if !shouldRetry(nextAddrErr, retries) {
			return []byte{}, false, nextAddrErr
		}

		time.Sleep(100 * time.Millisecond)
		return conn.createGeneratedAddressWithRetries(prefix, mutExclPrefixes, name, addrIsGreater, incAddr, retries - 1)
	}

	if addrIsGreater(nextAddr, addrRange.LastAddress) {
		//Range is full
		return []byte{}, true, nil
	}

	isHardcoded, isHarcodedErr := conn.addressIsHardcoded(prefix, nextAddr)
	if isHarcodedErr != nil {
		if !shouldRetry(isHarcodedErr, retries) {
			return []byte{}, false, isHarcodedErr
		}

		time.Sleep(100 * time.Millisecond)
		return conn.createGeneratedAddressWithRetries(prefix, mutExclPrefixes, name, addrIsGreater, incAddr, retries - 1)
	}

	for isHardcoded {
		nextAddr = incAddr(nextAddr)

		if addrIsGreater(nextAddr, addrRange.LastAddress) {
			//Range is full
			return []byte{}, true, nil
		}

		isHardcoded, isHarcodedErr = conn.addressIsHardcoded(prefix, nextAddr)
		if isHarcodedErr != nil {
			if !shouldRetry(isHarcodedErr, retries) {
				return []byte{}, false, isHarcodedErr
			}
	
			time.Sleep(100 * time.Millisecond)
			return conn.createGeneratedAddressWithRetries(prefix, mutExclPrefixes, name, addrIsGreater, incAddr, retries - 1)
		}
	}

	tx := conn.Client.Txn(ctx).If(
		slices.Concat(
			[]clientv3.Cmp{
				clientv3.Compare(clientv3.Version(addrRangeKeys.NextAddress), "=", nextAddrVer),
				clientv3.Compare(clientv3.Version(addrKeyPrefixes.HardcodedAddress + string(nextAddr)), "=", 0),
			},
			nameNoPresent,
		)...
	).Then(
		clientv3.OpPut(addrKeyPrefixes.GeneratedAddress + string(nextAddr), name),
		clientv3.OpPut(addrRangeKeys.NextAddress, string(incAddr(nextAddr))),
		clientv3.OpPut(addrKeyPrefixes.Name + name, string(nextAddr)),
	)

	resp, txErr := tx.Commit()
	if txErr != nil {
		if !shouldRetry(txErr, retries) {
			return []byte{}, false, txErr
		}

		time.Sleep(100 * time.Millisecond)
		return conn.createGeneratedAddressWithRetries(prefix, mutExclPrefixes, name, addrIsGreater, incAddr, retries - 1)
	}

	if !resp.Succeeded {
		if retries <= 0 {
			return []byte{}, false, errors.New("Failed to create generated address: Selected name has already been assigned")
		}

		return conn.createGeneratedAddressWithRetries(prefix, mutExclPrefixes, name, addrIsGreater, incAddr, retries - 1)
	}

	return nextAddr, false, nil
}

func (conn *EtcdConnection) CreateGeneratedAddress(prefix string, name string, addrIsGreater AddressIsGreater, incAddr IncrementAddress) ([]byte, error) {
	addr, full, err := conn.createGeneratedAddressWithRetries(prefix, []string{prefix}, name, addrIsGreater, incAddr, conn.Retries)
	if err != nil {
		return addr, err
	}
	if full {
		return addr, errors.New("Error creating generated address: Address range ran out of addresses")
	}

	return addr, nil
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
func (conn *EtcdConnection) deleteGeneratedAddressWithRetries(prefix string, name string, address []byte, prettify PrettifyAddr, retries int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conn.Timeout)*time.Second)
	defer cancel()

	addrKeyPrefixes := GenerateAddrEtcdKeyPrefixes(prefix)

	tx := conn.Client.Txn(ctx).If(
		clientv3.Compare(clientv3.Version(addrKeyPrefixes.GeneratedAddress + string(address)), ">", 0),
		clientv3.Compare(clientv3.Version(addrKeyPrefixes.Name + name), ">", 0),
	).Then(
		clientv3.OpDelete(addrKeyPrefixes.GeneratedAddress + string(address)),
		clientv3.OpDelete(addrKeyPrefixes.Name + name),
		clientv3.OpPut(addrKeyPrefixes.DeletedAddress  + string(address), name),
	)

	resp, txErr := tx.Commit()
	if txErr != nil {
		if !shouldRetry(txErr, retries) {
			return txErr
		}

		time.Sleep(100 * time.Millisecond)
		return conn.deleteGeneratedAddressWithRetries(prefix, name, address, prettify, retries - 1)
	}

	if !resp.Succeeded {
		return errors.New(fmt.Sprintf("Failed to delete generated address '%s': Either address or name have not been assigned or address was already deleted", prettify(address)))
	}

	return nil
}

func (conn *EtcdConnection) DeleteGeneratedAddress(prefix string, name string, address []byte, prettify PrettifyAddr) error {
	return conn.deleteGeneratedAddressWithRetries(prefix, name, address, prettify, conn.Retries)
}

func (conn *EtcdConnection) findAddressWithRetries(prefix string, name string, retries int) ([]byte, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conn.Timeout)*time.Second)
	defer cancel()

	addrKeyPrefixes := GenerateAddrEtcdKeyPrefixes(prefix)

	getRes, err := conn.Client.Get(ctx, addrKeyPrefixes.Name + name)
	if err != nil {
		if !shouldRetry(err, retries) {
			return []byte{}, false, err
		}

		time.Sleep(100 * time.Millisecond)
		return conn.findAddressWithRetries(prefix, name, retries - 1)
	}

	if len(getRes.Kvs) == 0 {
		return []byte{}, false, nil
	}

	return getRes.Kvs[0].Value, true, nil
}

func (conn *EtcdConnection) FindAddress(prefix string, name string) ([]byte, bool, error) {
	return conn.findAddressWithRetries(prefix, name, conn.Retries)
}

func (conn *EtcdConnection) GetAddress(prefix string, name string) ([]byte, error) {
	addr, found, err := conn.FindAddress(prefix, name)
	
	if err == nil && (!found) {
		return []byte{}, errors.New(fmt.Sprintf("Error retrieving address with name '%s': Name not found", name))
	}

	return addr, err
}

func (conn *EtcdConnection) getAddressDetailsWithRetries(prefix string, name string, retries int) (bool, bool, []byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conn.Timeout)*time.Second)
	defer cancel()

	addrKeyPrefixes := GenerateAddrEtcdKeyPrefixes(prefix)

	getRes, err := conn.Client.Get(ctx, addrKeyPrefixes.Name + name)
	if err != nil {
		if !shouldRetry(err, retries) {
			return false, false, []byte{}, err
		}

		time.Sleep(100 * time.Millisecond)
		return conn.getAddressDetailsWithRetries(prefix, name, retries - 1)
	}

	if len(getRes.Kvs) == 0 {
		return false, false, []byte{}, nil
	}

	isHardcoded, isHardcodedErr := conn.addressIsHardcoded(prefix, getRes.Kvs[0].Value)
	if isHardcodedErr != nil {
		if !shouldRetry(isHardcodedErr, retries) {
			return false, false, []byte{}, isHardcodedErr
		}

		time.Sleep(100 * time.Millisecond)
		return conn.getAddressDetailsWithRetries(prefix, name, retries - 1)
	}

	return true, isHardcoded, getRes.Kvs[0].Value, nil
}

func (conn *EtcdConnection) GetAddressDetails(prefix string, name string) (bool, bool, []byte, error) {
	return conn.getAddressDetailsWithRetries(prefix, name, conn.Retries)
}

//Multi-range methods
func (conn *EtcdConnection) FindAddressRangeByBoundaries(prefixes []string, addr []byte) (string, AddressRange, bool, error) {
	for _, prefix := range prefixes {
		addrRange, addrRangeExists, addrRangeErr := conn.GetAddrRange(prefix)
		if addrRangeErr != nil {
			return "", AddressRange{}, false, addrRangeErr
		}
		if !addrRangeExists {
			return "", AddressRange{}, false, errors.New(fmt.Sprintf("Range with prefix '%s' does not exist", prefix))
		}

		if len(addr) != len(addrRange.FirstAddress) {
			return "", AddressRange{}, false, errors.New(fmt.Sprintf("Range with prefix '%s' does not match passed address format: Address type mismatch", prefix))
		}

		if AddressWithinBoundaries(addr, addrRange.FirstAddress, addrRange.LastAddress) {
			return prefix, addrRange, true, nil
		}
	}

	return "", AddressRange{}, false, nil
}

func (conn *EtcdConnection) FindAddressDetailsInRanges(prefixes []string, name string) (bool, bool, []byte, string, error) {
	for _, prefix := range prefixes {
		addrExists, addrIsHardcoded, addr, detailsErr := conn.GetAddressDetails(prefix, name)
		if detailsErr != nil {
			return false, false, []byte{}, "", detailsErr
		}

		if addrExists {
			return addrExists, addrIsHardcoded, addr, prefix, nil
		}
	}

	return false, false, []byte{}, "", nil
}