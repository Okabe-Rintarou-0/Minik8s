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
	cli    *clientv3.Client
	//kv     clientv3.KV
	//err    error
)

func startClient() (err error) {
	cli, err = clientv3.New(config)
	return
}

func init() {
	config = clientv3.Config{
		Endpoints:            []string{"localhost:2379"},
		DialTimeout:          30 * time.Second,
		DialKeepAliveTimeout: 30 * time.Second,
	}
}

func checkAndStartClient() error {
	if cli == nil {
		err := startClient()
		if err != nil {
			log.Printf("[etcd] Connect to etcd failed, err:%v\n", err)
		}
	}
	return nil
}

func Put(key, value string) (err error) {
	if err = checkAndStartClient(); err != nil {
		return err
	}
	_, err = cli.Put(ctx, key, value)
	return err
}

func Get(key string) (value string, err error) {
	if err = checkAndStartClient(); err != nil {
		return "", err
	}
	var resp *clientv3.GetResponse
	resp, err = cli.Get(ctx, key)
	if err != nil {
		return "", err
	}

	if len(resp.Kvs) > 0 {
		return string(resp.Kvs[0].Value), nil
	} else {
		return "", nil
	}
}

func Delete(key string) (err error) {
	if err = checkAndStartClient(); err != nil {
		return err
	}
	_, err = cli.Delete(ctx, key)
	return err
}

func GetAll(keyPrefix string) (values []string, err error) {
	if err = checkAndStartClient(); err != nil {
		return nil, err
	}
	var resp *clientv3.GetResponse
	resp, err = cli.Get(ctx, keyPrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	for _, kv := range resp.Kvs {
		values = append(values, string(kv.Value))
	}
	return values, nil
}
