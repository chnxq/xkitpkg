package nacos

import (
	"testing"

	conf "github.com/chnxq/xkitpkg/conf/v1"
	"github.com/stretchr/testify/assert"
)

func TestNewNacosRegistry(t *testing.T) {
	cfg := conf.Registry{
		Nacos: &conf.Registry_Nacos{
			Address: "127.0.0.1",
			Port:    8848,
		},
	}

	reg, err := NewRegistry(&cfg)
	assert.Nil(t, err)
	assert.NotNil(t, reg)
}
