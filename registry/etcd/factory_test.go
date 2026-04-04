package etcd

import (
	"testing"

	conf "github.com/chnxq/xkitpkg/conf/v1"
	"github.com/stretchr/testify/assert"
)

func TestNewEtcdRegistry(t *testing.T) {
	cfg := conf.Registry{
		Etcd: &conf.Registry_Etcd{
			Endpoints: []string{"127.0.0.1:2379"},
		},
	}

	reg, err := NewRegistry(&cfg)
	assert.Nil(t, err)
	assert.NotNil(t, reg)
}
