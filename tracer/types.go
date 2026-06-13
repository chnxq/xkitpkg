package tracer

type Type string

const (
	Std      Type = "std"
	OtlpHttp Type = "otlp-http"
	OtlpGrpc Type = "otlp-grpc"
	Aliyun   Type = "aliyun"
	Tencent  Type = "tencent"
)
