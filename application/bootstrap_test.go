package application

import (
	"context"
	"testing"

	kit "github.com/chnxq/XGoKit"
	"github.com/chnxq/xkitpkg/conf/v1"
	"github.com/stretchr/testify/assert"
)

func initApp(ctx *Context) (*kit.App, func(), error) {
	app := NewApp(ctx)
	return app, func() {
	}, nil
}

func TestBootstrapWithNameVersion(t *testing.T) {
	serviceName := "test"
	version := "v0.0.1"

	ctx := NewContext(context.Background(), &conf.AppInfo{
		Project: "",
		AppId:   serviceName,
		Version: version,
	})

	err := RunApp(ctx, initApp)
	assert.Nil(t, err)
}

func TestNewInstanceId(t *testing.T) {
	instanceId := NewInstanceId("gowind-test-service", "1.0.0", "127.0.0.1", "8000")
	t.Logf("InstanceId: %s", instanceId)
}
