// Package conf provides access to the configuration protocol buffer types.
// This is a compatibility package that exports the v1 types.
package conf

import pb "github.com/chnxq/xkitpkg/conf/v1"

// Re-export all types from v1
type ServerConfig = pb.ServerConfig
type Server = pb.Server
type Client = pb.Client
type Data = pb.Data
type Tracer = pb.Tracer
type Logger = pb.Logger
type Registry = pb.Registry
type OSS = pb.OSS
type Notification = pb.Notification
type AppInfo = pb.AppInfo
type TLS = pb.TLS
type TLS_Config = pb.TLS_Config
type Authentication = pb.Authentication
type Authorization = pb.Authorization
type RemoteConfig = pb.RemoteConfig
type Middleware = pb.Middleware
type Script = pb.Script
