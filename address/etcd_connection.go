package address

import (
	"google.golang.org/grpc/codes"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdConnection struct {
	Client  *clientv3.Client
	Timeout int
	Retries int
	Strict  bool
}

func shouldRetry(err error, retries int) bool {
	etcdErr, ok := err.(rpctypes.EtcdError)
	if !ok {
		return false
	}
	
	if etcdErr.Code() != codes.Unavailable || retries <= 0 {
		return false
	}

	return true
}