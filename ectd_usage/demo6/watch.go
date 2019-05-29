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
		watcher clientv3.Watcher
		watchRespChan <-chan clientv3.WatchResponse
		watchResp *clientv3.WatchResponse
		event *clientv3.Event
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

	kv = clientv3.NewKV(client)

	// simulate kv change
	go func() {
		for {
			kv.Put(context.TODO(), "/f/b", "test")

			kv.Delete(context.TODO(), "/f/b")

			time.Sleep(1 * time.Second)
		}
	}()

	watcher = clientv3.NewWatcher(client)
	watchRespChan = watcher.Watch(context.TODO(), "/f/b")

	for watchResp = range watchRespChan {
		for _, event = range watchResp.Events {
			switch  event.Type {
			case mvccpb.PUT:
				fmt.Println("modify:", string(event.Kv.Value), "revision:", string(event.Kv.CreateRevision), event.Kv.ModRevision))
			case mvccpb.DELETE:
				fmt.Println("delete:", string(event.Kv.ModRevision))
			}
		}
	}

	ctx, cancelFunc := context.WithCancel(context.TODO())
	time.AfterFunc(5*time.Second, func(){
		cancelFunc()
	})

	fmt.Print(ctx.Err())
}
