package main

import (
	"context"
	"fmt"
	"github.com/etcd-io/etcd/clientv3"
	"time"
)

func main() {
	var (
		config clientv3.Config
		client * clientv3.Client
		err error
		kv clientv3.KV
		getResp *clientv3.GetResponse
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

	if getResp, err = kv.Get(context.TODO(), "/foo/bar", clientv3.WithCountOnly()); err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println(getResp.Kvs, getResp.Count)
	}

}
