package tracer

import (
	"context"

	conf "github.com/chnxq/xkitpkg/conf/v1"
)

// ReloadTracerProvider rebuilds and replaces the global tracer provider from runtime config.
// It returns false when the required config is incomplete and no reload is attempted.
func ReloadTracerProvider(ctx context.Context, cfg *conf.Tracer, appInfo *conf.AppInfo) (bool, error) {
	if cfg == nil || appInfo == nil {
		return false, nil
	}
	if _, _, err := NewTracerProviderWithShutdown(ctx, cfg, appInfo); err != nil {
		return false, err
	}
	return true, nil
}
