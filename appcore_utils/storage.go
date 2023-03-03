package appcore_utils

import (
	"github.com/minio/minio-go"
)

func NewStorage(configs *Configurations) *minio.Client {
	minioClient, err := minio.New(configs.MinioURL, configs.MinioAccessKey, configs.MinioSecretKey, configs.MinioSSL)
	if err != nil {
		panic("cannot connect storage (minio) >> " + err.Error())
	}
	return minioClient
}
