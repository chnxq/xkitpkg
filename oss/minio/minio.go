package minio

import (
	"github.com/chnxq/xkitpkg/conf"
	"github.com/chnxq/xkitpkg/logger/log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewClient(conf *conf.OSS) *minio.Client {
	impl, err := minio.New(conf.Minio.Endpoint,
		&minio.Options{
			Creds:  credentials.NewStaticV4(conf.Minio.AccessKey, conf.Minio.SecretKey, conf.Minio.Token),
			Secure: conf.Minio.UseSsl,
		},
	)
	if err != nil {
		log.Fatal("failed opening connection to minio", err)
		return nil
	}

	return impl
}
