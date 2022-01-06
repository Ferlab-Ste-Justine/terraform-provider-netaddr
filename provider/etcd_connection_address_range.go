package provider

import (
	"context"
	"errors"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type AddressRange struct {
	Type         string
	FirstAddress []byte
	LastAddress  []byte
}

type AddrRangeEtcdKeys struct {
	Type         string
	FirstAddress string
	LastAddress  string
	NextAddress  string
}

func GenerateAddrRangeEtcdKeys(rangePrefix string) AddrRangeEtcdKeys {
	return AddrRangeEtcdKeys{
		Type: rangePrefix + "info/type",
		FirstAddress: rangePrefix + "info/firstaddr",
		LastAddress: rangePrefix + "info/lastaddr",
		NextAddress: rangePrefix + "data/nextaddr",
	}
}

func (conn *EtcdConnection) createAddrRangeWithRetries(prefix string, addrRange AddressRange, retries int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conn.Timeout)*time.Second)
	defer cancel()

	rangeKeys := GenerateAddrRangeEtcdKeys(prefix)
	tx := conn.Client.Txn(ctx).Then(
		clientv3.OpPut(rangeKeys.Type, string(addrRange.Type)),
		clientv3.OpPut(rangeKeys.FirstAddress, string(addrRange.FirstAddress)),
		clientv3.OpPut(rangeKeys.LastAddress, string(addrRange.LastAddress)),
		clientv3.OpPut(rangeKeys.NextAddress, string(addrRange.FirstAddress)),
	)

	resp, err := tx.Commit()
	if err != nil {
		if !shouldRetry(err, retries) {
			return err
		}

		time.Sleep(100 * time.Millisecond)
		return conn.createAddrRangeWithRetries(prefix, addrRange, retries - 1)
	}

	if !resp.Succeeded {
		return errors.New(fmt.Sprintf("Failed to address range '%s' for unforeseen reason: Transaction with no condition failed", prefix))
	}

	return nil
}

func (conn *EtcdConnection) CreateAddrRange(prefix string, addrRange AddressRange) error {
	return conn.createAddrRangeWithRetries(prefix, addrRange, conn.Retries)
}

func (conn *EtcdConnection) getAddrRangeWithRetries(prefix string, retries int) (AddressRange, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conn.Timeout)*time.Second)
	defer cancel()
	var addrRange AddressRange

	infoKeys := prefix + "info/"
	getRes, err := conn.Client.Get(ctx, infoKeys, clientv3.WithPrefix())

	if err != nil {
		if !shouldRetry(err, retries) {
			return addrRange, false, err
		}

		time.Sleep(100 * time.Millisecond)
		return conn.getAddrRangeWithRetries(prefix, retries - 1)
	}

	if len(getRes.Kvs) != 3 {
		return addrRange, false, nil
	}

	rangeKeys := GenerateAddrRangeEtcdKeys(prefix)
	for _, kv := range getRes.Kvs {
		switch string(kv.Key) {
		case (rangeKeys.Type):
			addrRange.Type = string(kv.Value)
		case (rangeKeys.FirstAddress):
			addrRange.FirstAddress = kv.Value
		case (rangeKeys.LastAddress):
			addrRange.LastAddress = kv.Value
		}
	}

	return addrRange, true, nil
}

func (conn *EtcdConnection) GetAddrRange(prefix string) (AddressRange, bool, error) {
	return conn.getAddrRangeWithRetries(prefix, conn.Retries)
}

func (conn *EtcdConnection) destroyAddrRangeWithRetries(prefix string, retries int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conn.Timeout)*time.Second)
	defer cancel()

	_, err := conn.Client.Delete(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		if !shouldRetry(err, retries) {
			return err
		}

		time.Sleep(100 * time.Millisecond)
		return conn.destroyAddrRangeWithRetries(prefix, retries - 1)
	}

	return nil
}

func (conn *EtcdConnection) DestroyAddrRange(prefix string) error {
	return conn.destroyAddrRangeWithRetries(prefix, conn.Retries)
}