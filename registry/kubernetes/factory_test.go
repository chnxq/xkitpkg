package kubernetes

import (
	"testing"

	conf "github.com/chnxq/xkitpkg/conf/v1"
	"github.com/stretchr/testify/assert"
)

func TestNewKubernetesRegistry(t *testing.T) {
	var cfg conf.Registry
	reg, err := NewRegistry(&cfg)
	assert.Nil(t, err)
	assert.NotNil(t, reg)
}
