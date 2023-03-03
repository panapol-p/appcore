package appcore

import (
	"context"
	"flag"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/memphisdev/memphis.go"
	"github.com/minio/minio-go"
	"github.com/panapol-p/appcore/appcore_handler"
	"github.com/panapol-p/appcore/appcore_internal/appcore_model"
	"github.com/panapol-p/appcore/appcore_router"
	"github.com/panapol-p/appcore/appcore_utils"
	"github.com/sirupsen/logrus"
	requestID "github.com/sumit-tembe/gin-requestid"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
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

	restAPIServer *http.Server
}

func init() {
	appcore_model.Port = new(string)
	appcore_model.Port = flag.String("port", "", "your service port")

	appcore_model.GrpcPort = new(string)
	appcore_model.GrpcPort = flag.String("grpcPort", "", "your service grpc port")

	appcore_model.IP = new(string)
	appcore_model.IP = flag.String("ip", "", "your service ip")

	flag.Parse()
	if *appcore_model.Port == "" {
		*appcore_model.Port = "8000"
	}

	if *appcore_model.GrpcPort == "" {
		*appcore_model.GrpcPort = "1" + *appcore_model.Port
	}

	if *appcore_model.IP == "" {
		*appcore_model.IP = "0.0.0.0"
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
	s.restAPIServer = initGinAPI(s.Config, s.ApiHandler, s.Logger)

	//check GRPC
	if s.ApiHandler.Module.GrpcServer() != nil {
		initGRPC(s.ApiHandler, s.Logger)
	}
}

func (s *Service) Stop() {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//shutdown api
	s.Logger.Info("shutdown API.")
	if err := s.restAPIServer.Shutdown(timeoutCtx); err != nil {
		s.Logger.Error(err.Error())
	}

	//shutdown grpc
	if s.ApiHandler.Module.GrpcServer() != nil {
		s.Logger.Info("shutdown GRPC.")
		s.ApiHandler.Module.GrpcServer().GracefulStop()
	}

	//shutdown db
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

	//shutdown message broker
	if s.MessageBroker != nil {
		s.Logger.Info("disconnect MessageBroker")
		s.MessageBroker.Close()
	}

	//shutdown cache
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
	return r.ListenAndServe(*appcore_model.IP, *appcore_model.Port)
}

func initGRPC(h *appcore_handler.ApiHandler, logger *logrus.Logger) {
	logger.Info("listen grpc on ", *appcore_model.IP+":"+*appcore_model.GrpcPort)
	grpcListener, err := net.Listen("tcp", *appcore_model.IP+":"+*appcore_model.GrpcPort)
	if err != nil {
		logger.Fatalf("GRPC: failed to listen: %v", err)
	}

	s := h.Module.GrpcServer()

	// Register gRPC health check service
	healthCheck := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthCheck)

	go func() {
		err = s.Serve(grpcListener)
		if err != nil {
			logger.Fatalf("GRPC: failed to listen: %v", err)
		}

	}()
}
