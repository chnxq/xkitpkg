package application

import (
	kit "github.com/chnxq/XGoKit"
)

// InitAppFunc 应用初始化函数类型
type InitAppFunc func(ctx *AppCtx) (app *kit.App, cleanup func(), err error)
