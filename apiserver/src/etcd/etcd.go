package etcd

import (
	"context"
	"go.etcd.io/etcd/clientv3"
	"log"
	"time"
)

var (
	ctx    = context.Background()
	config clientv3.Config
	//cli    *clientv3.Client
	//kv     clientv3.KV
	//err    error
)

func init() {
	config = clientv3.Config{
		Endpoints:            []string{"localhost:2379"},
		DialTimeout:          30 * time.Second,
		DialKeepAliveTimeout: 30 * time.Second,
	}

	cli, err := clientv3.New(config)
	defer cli.Close()
	if err != nil {
		log.Printf("[etcd] connect to etcd failed, err:%v\n", err)
	}
	log.Printf("[etcd] connect to etcd success\n")
}

func Put(key, value string) {
	cli, err := clientv3.New(config)
	defer cli.Close()
	if err != nil {
		log.Printf("[etcd] connect to etcd failed, err:%v\n", err)
	}

	if _, err := cli.Put(ctx, key, value); err != nil {
		log.Fatal(err)
	}
}

func Get(key string) string {
	cli, err := clientv3.New(config)
	defer cli.Close()
	if err != nil {
		log.Printf("[etcd] connect to etcd failed, err:%v\n", err)
	}
	resp, err := cli.Get(ctx, key)
	if err != nil {
		log.Fatal(err)
	}

	if len(resp.Kvs) > 0 {
		return string(resp.Kvs[0].Value)
	} else {
		return ""
	}
}

func Delete(key string) {
	cli, err := clientv3.New(config)
	defer cli.Close()
	if err != nil {
		log.Printf("[etcd] connect to etcd failed, err:%v\n", err)
	}
	_, err = cli.Delete(ctx, key)
	if err != nil {
		log.Fatal(err)
	}

	//log.Printf("Delete keys: %v\n", resp.Deleted)

}
