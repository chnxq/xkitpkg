::指定起始文件夹
set DIR=%cd%

go get all
go mod tidy

cd %DIR%\app
go get all
go mod tidy

cd %DIR%\cache
go get all
go mod tidy

cd %DIR%\conf
go get all
go mod tidy

cd %DIR%\config
go get all
go mod tidy

cd %DIR%\config\consul
go get all
go mod tidy

cd %DIR%\config\etcd
go get all
go mod tidy

cd %DIR%\config\kubernetes
go get all
go mod tidy

cd %DIR%\config\nacos
go get all
go mod tidy


cd %DIR%\logger
go get all
go mod tidy

cd %DIR%\logger\aliyun
go get all
go mod tidy

cd %DIR%\logger\fluentd
go get all
go mod tidy

cd %DIR%\logger\logrus
go get all
go mod tidy

cd %DIR%\logger\tencent
go get all
go mod tidy

cd %DIR%\logger\zerolog
go get all
go mod tidy

cd %DIR%\logger\zap
go get all
go mod tidy


cd %DIR%\registry
go get all
go mod tidy

cd %DIR%\registry\consul
go get all
go mod tidy

cd %DIR%\registry\etcd
go get all
go mod tidy

cd %DIR%\registry\nacos
go get all
go mod tidy

cd %DIR%\tracer
go get all
go mod tidy

cd %DIR%\transport
go get all
go mod tidy

cd %DIR%\oss\minio
go get all
go mod tidy

cd %DIR%\server_utils
go get all
go mod tidy


