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
		lease clientv3.Lease
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseId clientv3.LeaseID
		putResp *clientv3.PutResponse
		kv clientv3.KV
		getResp clientv3.GetResponse
		keepResp *clientv3.LeaseKeepAliveResponse
		keepRespChan <-chan *clientv3.LeaseKeepAliveResponse
		ctx context.Context
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

	// apply a new lease
	lease = clientv3.NewLease(client)

	// apply 10s
	if leaseGrantResp, err = lease.Grant(context.TODO(), 10); err != nil {
		return
	}

	// put kv with lease
	// get leaseID
	leaseId = leaseGrantResp.ID

	ctx, _ = context.WithTimeout(context.TODO(), 5* time.Second)

	// auto add lease
	//if keepRespChan, err = lease.KeepAlive(context.TODO(), leaseId); err != nil {
	if keepRespChan, err = lease.KeepAlive(ctx, leaseId); err != nil {
		return
	}

	go func() {
		for {
			select {
			case keepResp = <- keepRespChan:
				if keepRespChan == nil {
					fmt.Println("lease invalid")
					break
				} else {
					fmt.Println("get auto lease resp:", keepResp.ID)
				}

			}
		}
	}()



	// get kv api subset
	kv = clientv3.NewKV(client)

	if putResp, err = kv.Put(context.TODO(), "/foo/bar", "abc", clientv3.WithLease(leaseId)); err != nil {
		return
	}

	fmt.Println("Write KEY success:", putResp.Header.Revision)

	// sensor out of date
	for {
		if getResp, err = kv.Get(context.TODO(), "/foo/bar"); err != nil {
			fmt.Println(err)
			return
		}

		if getResp.Count == 0 {
			fmt.Println("out of date")
			break
		}

		fmt.Println("not out of date:", getResp.Kvs)

		time.Sleep(time.Second * 2)
	}
}
