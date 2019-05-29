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
		putResp *clientv3.PutResponse
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

	if putResp, err = kv.Put(context.TODO(), "/foo/bar", "hello", clientv3.WithPrevKV()); err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("Revision:", putResp.Header.Revision)
		if putResp.PrevKv != nil {
			fmt.Println("PV:", string(putResp.PrevKv.Value))
		}
	}

}
