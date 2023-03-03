package appcore_utils

import (
	"github.com/spf13/viper"
)

// Configurations wraps all the config variables required by the service
type Configurations struct {
	//Gin Mode
	GinIsReleaseMode bool

	//observability
	ObserveIsActive     bool
	ObserveOTLPEndpoint string
	ObserveInsecureMode string

	//database
	PostgresConnString string

	//Redis
	RedisUrl  string
	RedisPass string

	//Message Broker
	MemphisHost     string
	MemphisUsername string
	MemphisToken    string

	//Storage
	MinioURL        string
	MinioSSL        bool
	MinioAccessKey  string
	MinioSecretKey  string
	MinioBucketName string
}

// NewConfigurations returns a new Configuration object
func NewConfigurations() *Configurations {
	viper.AutomaticEnv()
	viper.SetDefault("GIN_IS_RELEASE_MODE", false)
	viper.SetDefault("OBSERVE_IS_ACTIVE", false)
	viper.SetDefault("OBSERVE_OTLP_ENDPOINT", "localhost:4317")
	viper.SetDefault("OBSERVE_INSECURE_MODE", "false")
	viper.SetDefault("POSTGRES_URL", "host=127.0.0.1 user=root password=example dbname=osp port=5432 sslmode=disable TimeZone=Asia/Bangkok")
	viper.SetDefault("REDIS_URL", "localhost:6379")
	viper.SetDefault("REDIS_PASS", "password123")
	viper.SetDefault("MEMPHIS_HOST", "localhost:9000")
	viper.SetDefault("MEMPHIS_USERNAME", "root")
	viper.SetDefault("MEMPHIS_TOKEN", "memphis")
	viper.SetDefault("MINIO_URL", "localhost:9010")
	viper.SetDefault("MINIO_SSL", false)
	viper.SetDefault("MINIO_ACCESS_KEY", "minioadmin")
	viper.SetDefault("MINIO_SECRET_KEY", "minioadmin")
	viper.SetDefault("MINIO_BUCKET_NAME", "miniobucket")

	configs := &Configurations{
		GinIsReleaseMode:    viper.GetBool("GIN_IS_RELEASE_MODE"),
		ObserveIsActive:     viper.GetBool("OBSERVE_IS_ACTIVE"),
		ObserveOTLPEndpoint: viper.GetString("OBSERVE_OTLP_ENDPOINT"),
		ObserveInsecureMode: viper.GetString("OBSERVE_INSECURE_MODE"),
		PostgresConnString:  viper.GetString("POSTGRES_URL"),
		RedisUrl:            viper.GetString("REDIS_URL"),
		RedisPass:           viper.GetString("REDIS_PASS"),
		MemphisHost:         viper.GetString("MEMPHIS_HOST"),
		MemphisUsername:     viper.GetString("MEMPHIS_USERNAME"),
		MemphisToken:        viper.GetString("MEMPHIS_TOKEN"),
		MinioURL:            viper.GetString("MINIO_URL"),
		MinioSSL:            viper.GetBool("MINIO_SSL"),
		MinioAccessKey:      viper.GetString("MINIO_ACCESS_KEY"),
		MinioSecretKey:      viper.GetString("MINIO_SECRET_KEY"),
		MinioBucketName:     viper.GetString("MINIO_BUCKET_NAME"),
	}
	return configs
}
