package consul

import (
	"testing"

	"github.com/chnxq/xkitpkg/conf/v1"
	"github.com/stretchr/testify/assert"
)

func TestNewConsulRegistry(t *testing.T) {
	cfg := conf.Registry{
		Consul: &conf.Registry_Consul{
			Scheme:      "http",
			Address:     "localhost:8500",
			HealthCheck: false,
		},
	}

	reg, err := NewRegistry(&cfg)
	assert.Nil(t, err)
	assert.NotNil(t, reg)
}
