package provider

import (
	"context"
	"errors"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type NetworkEtcdKeys struct {
	FirstIpv4 string
	LastIpv4 string
	NextIpv4 string
	FirstIpv6 string
	LastIpv6 string
	NextIpv6 string
	FirstMac string
	LastMac string
	NextMac string
}

func GenerateNetworkEtcdKeys(networkPrefix string) NetworkEtcdKeys {
	return NetworkEtcdKeys{
		FirstIpv4: networkPrefix + "info/ipv4/first",
		LastIpv4: networkPrefix + "info/ipv4/last",
		NextIpv4: networkPrefix + "data/ipv4/next",
		FirstIpv6: networkPrefix + "info/ipv6/first",
		LastIpv6: networkPrefix + "info/ipv6/last",
		NextIpv6: networkPrefix + "data/ipv6/next",
		FirstMac: networkPrefix + "info/mac/first",
		LastMac: networkPrefix + "info/mac/last",
		NextMac: networkPrefix + "data/mac/next",
	}
}

func (conn *EtcdConnection) createNetworkWithRetries(prefix string, net Network, retries int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conn.Timeout)*time.Second)
	defer cancel()

	netKeys := GenerateNetworkEtcdKeys(prefix)
	tx := conn.Client.Txn(ctx).Then(
		clientv3.OpPut(netKeys.FirstIpv4, string(net.Ipv4.First)),
		clientv3.OpPut(netKeys.LastIpv4, string(net.Ipv4.Last)),
		clientv3.OpPut(netKeys.NextIpv4, string(net.Ipv4.First)),
		clientv3.OpPut(netKeys.FirstIpv6, string(net.Ipv6.First)),
		clientv3.OpPut(netKeys.LastIpv6, string(net.Ipv6.Last)),
		clientv3.OpPut(netKeys.NextIpv6, string(net.Ipv6.First)),
		clientv3.OpPut(netKeys.FirstMac, string(net.Mac.First)),
		clientv3.OpPut(netKeys.LastMac, string(net.Mac.Last)),
		clientv3.OpPut(netKeys.NextMac, string(net.Mac.First)),
	)

	resp, err := tx.Commit()
	if err != nil {
		if !shouldRetry(err, retries) {
			return err
		}

		time.Sleep(100 * time.Millisecond)
		return conn.createNetworkWithRetries(prefix, net, retries - 1)
	}

	if !resp.Succeeded {
		return errors.New(fmt.Sprintf("Failed to create network '%s' for unforeseen reason: Transaction with no condition failed", prefix))
	}

	return nil
}

func (conn *EtcdConnection) CreateNetwork(prefix string, net Network) error {
	return conn.createNetworkWithRetries(prefix, net, conn.Retries)
}

func (conn *EtcdConnection) getNetworkInfoWithRetries(prefix string, retries int) (Network, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conn.Timeout)*time.Second)
	defer cancel()
	var net Network

	infoKeys := prefix + "info/"
	getRes, err := conn.Client.Get(ctx, infoKeys, clientv3.WithPrefix())

	if err != nil {
		if !shouldRetry(err, retries) {
			return net, false, err
		}

		time.Sleep(100 * time.Millisecond)
		return conn.getNetworkInfoWithRetries(prefix, retries - 1)
	}

	if len(getRes.Kvs) != 6 {
		return net, false, nil
	}

	netKeys := GenerateNetworkEtcdKeys(prefix)
	for _, kv := range getRes.Kvs {
		switch string(kv.Key) {
		case (netKeys.FirstIpv4):
			net.Ipv4.First = kv.Value
		case (netKeys.LastIpv4):
			net.Ipv4.Last = kv.Value
		case (netKeys.FirstIpv6):
			net.Ipv6.First = kv.Value
		case (netKeys.LastIpv6):
			net.Ipv6.Last = kv.Value
		case (netKeys.FirstMac):
			net.Mac.First = kv.Value
		case (netKeys.LastMac):
			net.Mac.Last = kv.Value
		}
	}

	return net, true, nil
}

func (conn *EtcdConnection) GetNetworkInfo(prefix string) (Network, bool, error) {
	return conn.getNetworkInfoWithRetries(prefix, conn.Retries)
}

func (conn *EtcdConnection) destroyNetworkWithRetries(prefix string, retries int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conn.Timeout)*time.Second)
	defer cancel()

	_, err := conn.Client.Delete(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		if !shouldRetry(err, retries) {
			return err
		}

		time.Sleep(100 * time.Millisecond)
		return conn.destroyNetworkWithRetries(prefix, retries - 1)
	}

	return nil

}

func (conn *EtcdConnection) DestroyNetwork(prefix string) error {
	return conn.destroyNetworkWithRetries(prefix, conn.Retries)
}