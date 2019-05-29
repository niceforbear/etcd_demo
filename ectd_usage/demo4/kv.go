package main

import (
	"context"
	"fmt"
	"github.com/etcd-io/etcd/clientv3"
	"github.com/etcd-io/etcd/mvcc/mvccpb"
	"time"
)

func main() {
	var (
		config clientv3.Config
		client * clientv3.Client
		err error
		kv clientv3.KV
		getResp *clientv3.GetResponse
		delResp *clientv3.DeleteResponse
		//idx int
		kvpair mvccpb.KeyValue
	)

	// config
	config = clientv3.Config{
		Endpoints: []string{"0.0.0.0:2379"},
		DialTimeout: 5 * time.Second,
	}

	// connect
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	// kv-pair object
	kv = clientv3.NewKV(client)

	if getResp, err = kv.Get(context.TODO(), "/foo/bar/", clientv3.WithPrefix()); err != nil {
		fmt.Println(err)
	} else {
		fmt.Print(getResp.Kvs)
	}

	if delResp, err = kv.Delete(context.TODO(), "/foo/bar/a", clientv3.WithPrevKV()); err != nil {
		return
	}
	if len(delResp.PrevKvs) != 0 {
		for _, kvpair = range delResp.PrevKvs {
			fmt.Println(string(kvpair.Key), string(kvpair.Value), kvpair.Version, kvpair.Lease)
		}
	}

	if delResp, err = kv.Delete(context.TODO(), "/foo/bar/a", clientv3.WithFromKey(), clientv3.WithLimit(2)); err != nil {
		return
	}
}
