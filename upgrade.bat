::指定起始文件夹
set DIR=%cd%

go get all
go mod tidy

cd %DIR%\conf
go get all
go mod tidy

cd %DIR%\config
go get all
go mod tidy

cd %DIR%\logger
go get all
go mod tidy

cd %DIR%\registry
go get all
go mod tidy

cd %DIR%\tracer
go get all
go mod tidy

cd %DIR%\transport
go get all
go mod tidy

