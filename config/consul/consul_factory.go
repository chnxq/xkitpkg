package consul

import (
	"github.com/chnxq/XGoKit/config"
	"github.com/chnxq/XGoKit/libs/config/consul"
	"github.com/chnxq/XGoKit/log"
	"github.com/chnxq/xkitpkg/conf/v1"
	"github.com/hashicorp/consul/api"
)

func ConfigFactory(cfg *conf.RemoteConfig) (config.Source, error) {
	consulClient, err := api.NewClient(
		&api.Config{
			Address:    cfg.GetConsul().GetAddress(),
			Scheme:     cfg.GetConsul().GetScheme(),
			PathPrefix: "",
			Datacenter: "",
			Transport:  nil,
			HttpClient: nil,
			HttpAuth:   nil,
			WaitTime:   0,
			Token:      "",
			TokenFile:  "",
			Namespace:  "",
			Partition:  "",
			TLSConfig:  api.TLSConfig{},
		},
	)
	if err != nil {
		log.Fatal("create remote config consul client failed: ", err)
	}
	source, err := consul.New(consulClient, consul.WithPath("app/cart/configs/"))
	if err != nil {
		log.Fatal("create remote config consul source failed: ", err)
	}
	return source, nil
}
