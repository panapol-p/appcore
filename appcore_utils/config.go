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
	CockroachConnString string

	//Redis
	RedisUrl  string
	RedisPass string

	//Message Broker
	MemphisHost     string
	MemphisUsername string
	MemphisToken    string
}

// NewConfigurations returns a new Configuration object
func NewConfigurations() *Configurations {
	viper.AutomaticEnv()
	viper.SetDefault("GIN_IS_RELEASE_MODE", false)
	viper.SetDefault("OBSERVE_IS_ACTIVE", false)
	viper.SetDefault("OBSERVE_OTLP_ENDPOINT", "localhost:4317")
	viper.SetDefault("OBSERVE_INSECURE_MODE", "false")
	viper.SetDefault("COCKROACH_URL", "host=127.0.0.1 user=root password=example dbname=osp port=5432 sslmode=disable TimeZone=Asia/Bangkok")
	viper.SetDefault("REDIS_URL", "localhost:6379")
	viper.SetDefault("REDIS_PASS", "password123")
	viper.SetDefault("MEMPHIS_HOST", "localhost:9000")
	viper.SetDefault("MEMPHIS_USERNAME", "root")
	viper.SetDefault("MEMPHIS_TOKEN", "memphis")

	configs := &Configurations{
		GinIsReleaseMode:    viper.GetBool("GIN_IS_RELEASE_MODE"),
		ObserveIsActive:     viper.GetBool("OBSERVE_IS_ACTIVE"),
		ObserveOTLPEndpoint: viper.GetString("OBSERVE_OTLP_ENDPOINT"),
		ObserveInsecureMode: viper.GetString("OBSERVE_INSECURE_MODE"),
		CockroachConnString: viper.GetString("COCKROACH_URL"),
		RedisUrl:            viper.GetString("REDIS_URL"),
		RedisPass:           viper.GetString("REDIS_PASS"),
		MemphisHost:         viper.GetString("MEMPHIS_HOST"),
		MemphisUsername:     viper.GetString("MEMPHIS_USERNAME"),
		MemphisToken:        viper.GetString("MEMPHIS_TOKEN"),
	}
	return configs
}
