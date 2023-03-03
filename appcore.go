package appcore

import (
	"context"
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/memphisdev/memphis.go"
	"github.com/minio/minio-go"
	"github.com/panapol-p/appcore/appcore_handler"
	"github.com/panapol-p/appcore/appcore_router"
	"github.com/panapol-p/appcore/appcore_utils"
	"github.com/sirupsen/logrus"
	requestID "github.com/sumit-tembe/gin-requestid"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"gorm.io/gorm"
)

type Service struct {
	ServiceName   string
	Version       string
	ApiHandler    *appcore_handler.ApiHandler
	Config        *appcore_utils.Configurations
	Logger        *logrus.Logger
	DB            *gorm.DB
	Cache         *redis.Client
	MessageBroker *memphis.Conn
	Storage       *minio.Client

	server *http.Server
}

var port *string
var ip *string

func init() {
	port = new(string)
	port = flag.String("port", "", "your service port")

	ip = new(string)
	ip = flag.String("ip", "", "your service ip")

	flag.Parse()
	if *port == "" {
		*port = "8000"
	}

	if *ip == "" {
		*ip = "0.0.0.0"
	}
}

func NewService(serviceName, version string, apiHandler *appcore_handler.ApiHandler, config *appcore_utils.Configurations, logger *logrus.Logger) *Service {

	logger.Info(">>>>> service : " + serviceName + " " + version + " <<<<<")

	if config.ObserveIsActive {
		traceCleaner := appcore_utils.InitTracer(serviceName, config.ObserveOTLPEndpoint, config.ObserveInsecureMode)
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := traceCleaner(timeoutCtx); err != nil {
			logger.Error("observability -> ", err.Error())
			os.Exit(1)
		}
	}

	return &Service{
		ServiceName:   serviceName,
		Version:       version,
		ApiHandler:    apiHandler,
		Config:        config,
		Logger:        logger,
		DB:            apiHandler.DB,
		Cache:         apiHandler.Cache,
		MessageBroker: apiHandler.MessageBroker,
		Storage:       apiHandler.Storage,
	}
}

func (s *Service) Run() {
	s.server = initGinAPI(s.Config, s.ApiHandler, s.Logger)
}

func (s *Service) Stop() {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(timeoutCtx); err != nil {
		s.Logger.Error(err.Error())
	}

	if s.DB != nil {
		s.Logger.Info("disconnect DB")
		sqlDB, err := s.DB.DB()
		if err != nil {
			s.Logger.Error(err.Error())
		} else {
			err = sqlDB.Close()
			if err != nil {
				s.Logger.Error(err.Error())
			}
		}
	}

	if s.MessageBroker != nil {
		s.Logger.Info("disconnect MessageBroker")
		s.MessageBroker.Close()
	}

	if s.Cache != nil {
		s.Logger.Info("disconnect Cache")
		err := s.Cache.Close()
		if err != nil {
			s.Logger.Error(err.Error())
		}
	}

	s.Logger.Info("server was shutdown")
}

func initGinAPI(configs *appcore_utils.Configurations, h *appcore_handler.ApiHandler, logger *logrus.Logger) *http.Server {
	if configs.GinIsReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r := appcore_router.New(logger)

	if configs.ObserveIsActive {
		r.Use(otelgin.Middleware(h.ServiceName + "-" + h.Version))
	}

	r.Use(gin.LoggerWithConfig(requestID.GetLoggerConfig(nil, nil, nil)))

	v := r.Group("/api")
	v.GET("/ping", h.HealthCheck)

	h.Module.ModuleAPI(r)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"service": h.ServiceName + " service " + h.Version,
			"code":    "PAGE_NOT_FOUND",
			"message": "Page not found"})
	})
	return r.ListenAndServe(*ip, *port)
}
