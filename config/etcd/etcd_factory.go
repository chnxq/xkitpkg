package etcd

import (
	"log"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"

	"github.com/chnxq/XGoKit/config"
	etcdcfg "github.com/chnxq/XGoKit/libs/config/etcd"

	"github.com/chnxq/xkitpkg/conf/v1"
)

func ConfigFactory(cfg *conf.RemoteConfig) (config.Source, error) {
	// create an etcd client
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.GetEtcd().GetEndpoints(),
		DialTimeout: cfg.GetEtcd().GetTimeout().AsDuration(),
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
	})
	if err != nil {
		log.Fatal(err)
	}

	// configure the source, "path" is required
	source, err := etcdcfg.New(client, etcdcfg.WithPath(cfg.GetEtcd().GetKey()), etcdcfg.WithPrefix(true))
	if err != nil {
		log.Fatalln(err)
	}
	return source, nil
}
