package intergrated

import (
	"context"

	"github.com/casbin/casbin/v3/persist"
	jwtv5 "github.com/golang-jwt/jwt/v5"

	"github.com/chnxq/xkitpkg/conf/v1"
	"github.com/chnxq/xkitpkg/middleware"
	"github.com/chnxq/xkitpkg/middleware/auth/casbin"
	"github.com/chnxq/xkitpkg/middleware/auth/jwt"
	"github.com/chnxq/xkitpkg/middleware/selector"
)

func NewHttpJwtMiddleware(cfg *conf.Middleware_Auth, adapter persist.Adapter) middleware.Middleware {
	jwtBuilder := selector.Server(
		jwt.Server(func(token *jwtv5.Token) (interface{}, error) {
			return []byte(cfg.GetKey()), nil
		}),
		casbin.Server(
			casbin.WithCasbinPolicy(adapter),
			casbin.WithSecurityUserCreator(casbin.NewSecurityUser),
		),
	)
	whiteList := cfg.GetWhiteList()
	operationMap := make(map[string]struct{})
	if whiteList != nil {
		prefix := whiteList.GetPrefix()
		if len(prefix) != 0 {
			jwtBuilder.Prefix(prefix...)
		}
		regex := whiteList.GetRegex()
		if len(regex) != 0 {
			jwtBuilder.Regex(regex...)
		}
		path := whiteList.GetPath()
		if len(path) != 0 {
			jwtBuilder.Path(path...)
		}
		match := whiteList.GetMatch()
		if len(match) != 0 {
			for _, item := range match {
				operationMap[item] = struct{}{}
			}
		}
	}
	jwtBuilder.Match(func(ctx context.Context, operation string) bool {
		if _, ok := operationMap[operation]; ok {
			return false
		}
		return true
	})
	return jwtBuilder.Build()
}

func NewGrpcJwtMiddleware(cfg *conf.Middleware_Auth) middleware.Middleware {
	jwtBuilder := selector.Server(
		jwt.Server(func(token *jwtv5.Token) (interface{}, error) {
			return []byte(cfg.GetKey()), nil
		}),
	)
	whiteList := cfg.GetWhiteList()
	operationMap := make(map[string]struct{})
	if whiteList != nil {
		prefix := whiteList.GetPrefix()
		if len(prefix) != 0 {
			jwtBuilder.Prefix(prefix...)
		}
		regex := whiteList.GetRegex()
		if len(regex) != 0 {
			jwtBuilder.Regex(regex...)
		}
		path := whiteList.GetPath()
		if len(path) != 0 {
			jwtBuilder.Path(path...)
		}
		match := whiteList.GetMatch()
		if len(match) != 0 {
			for _, item := range match {
				operationMap[item] = struct{}{}
			}
		}
	}
	jwtBuilder.Match(func(ctx context.Context, operation string) bool {
		if _, ok := operationMap[operation]; ok {
			return false
		}
		return true
	})
	return jwtBuilder.Build()
}
