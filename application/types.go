package application

import (
	kit "github.com/chnxq/XGoKit"
)

// InitAppFunc 应用初始化函数类型
type InitAppFunc func(ctx *Context) (app *kit.App, cleanup func(), err error)
