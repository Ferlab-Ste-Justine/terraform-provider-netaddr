package address

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type AddrRangeKeyspace struct {
	Type               string
	FirstAddress       []byte
	LastAddress        []byte
	NextAddress        []byte
	Names              []AddressListEntry
	GeneratedAddresses []AddressListEntry
	HardcodedAddresses []AddressListEntry
	DeletedAddresses   []AddressListEntry
}

func (conn *EtcdConnection) getKeyspaceAddrListWithRetries(addrPrefix string, retries int) ([]AddressListEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conn.Timeout)*time.Second)
	defer cancel()

	getRes, err := conn.Client.Get(ctx, addrPrefix, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
	if err != nil {
		if !shouldRetry(err, retries) {
			return []AddressListEntry{}, err
		}

		time.Sleep(100 * time.Millisecond)
		return conn.getKeyspaceAddrListWithRetries(addrPrefix, retries - 1)
	}

	listing := make([]AddressListEntry, len(getRes.Kvs))
	for idx, val := range getRes.Kvs {
		address, _ := bytes.CutPrefix(val.Key, []byte(addrPrefix))
		listing[idx] = AddressListEntry{string(val.Value), address}
	}

	return listing, nil
}

func (conn *EtcdConnection) GetKeyspaceAddrList(addrPrefix string) ([]AddressListEntry, error) {
	return conn.getKeyspaceAddrListWithRetries(addrPrefix, conn.Retries)
}

func (conn *EtcdConnection) GetAddrRangeKeyspace(prefix string) (AddrRangeKeyspace, error) {
	addrRange, addrRangeExists, addrRangeErr := conn.GetAddrRange(prefix)
	if !addrRangeExists {
		return AddrRangeKeyspace{}, errors.New(fmt.Sprintf("Error retrieving address range at prefix '%s': Range does not exist", prefix))
	}
	if addrRangeErr != nil {
		return AddrRangeKeyspace{}, addrRangeErr
	}

	nextAddress, _, nextAddressErr := conn.getNextAddress(prefix)
	if nextAddressErr != nil {
		return AddrRangeKeyspace{}, nextAddressErr
	}

	namesList, namesListErr := conn.GetAddressList(prefix)
	if namesListErr != nil {
		return AddrRangeKeyspace{}, namesListErr
	}

	addrKeyPrefixes := GenerateAddrEtcdKeyPrefixes(prefix)

	generatedList, generatedListErr := conn.GetKeyspaceAddrList(addrKeyPrefixes.GeneratedAddress)
	if generatedListErr != nil {
		return AddrRangeKeyspace{}, generatedListErr
	}
	
	hardcodedList, hardcodedListErr := conn.GetKeyspaceAddrList(addrKeyPrefixes.HardcodedAddress)
	if hardcodedListErr != nil {
		return AddrRangeKeyspace{}, hardcodedListErr
	}
	
	deletedList, deletedListErr := conn.GetKeyspaceAddrList(addrKeyPrefixes.DeletedAddress)
	if deletedListErr != nil {
		return AddrRangeKeyspace{}, deletedListErr
	}

	return AddrRangeKeyspace{
		Type: addrRange.Type,
		FirstAddress: addrRange.FirstAddress,
		LastAddress: addrRange.LastAddress,
		NextAddress: nextAddress,
		Names: namesList,
		GeneratedAddresses: generatedList,
		HardcodedAddresses: hardcodedList,
		DeletedAddresses: deletedList,
	}, nil
}