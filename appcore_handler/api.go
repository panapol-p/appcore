package appcore_handler

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/memphisdev/memphis.go"
	"github.com/minio/minio-go"
	inf "github.com/panapol-p/appcore/appcore_internal/appcore_interface"
	"github.com/panapol-p/appcore/appcore_internal/appcore_model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"gorm.io/gorm"
)

type ApiHandler struct {
	ServiceName string
	Version     string
	Module      inf.Module

	DB            *gorm.DB
	Cache         *redis.Client
	MessageBroker *memphis.Conn
	Storage       *minio.Client
}

func NewHandler(serviceName, version string, module inf.Module, db *gorm.DB, cache *redis.Client, messageBroker *memphis.Conn, storage *minio.Client) *ApiHandler {
	return &ApiHandler{
		ServiceName:   serviceName,
		Version:       version,
		Module:        module,
		DB:            db,
		Cache:         cache,
		MessageBroker: messageBroker,
		Storage:       storage,
	}
}

func (h *ApiHandler) HealthCheck(c *gin.Context) {
	ctx := context.Background()

	isError := false
	errorMessage := ""

	defer func() {
		if isError {
			c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"service": h.ServiceName,
				"version": h.Version,
				"error":   errorMessage,
			})
		} else {
			c.JSON(http.StatusOK, map[string]interface{}{
				"service": h.ServiceName,
				"message": "pong",
				"version": h.Version,
			})
		}
	}()

	if h.DB != nil {
		sql, err := h.DB.DB()
		if err != nil {
			isError = true
			errorMessage = "cannot ping to database service"
			return
		}
		err = sql.Ping()
		if err != nil {
			isError = true
			errorMessage = "cannot ping to database service"
			return
		}
	}

	if h.Cache != nil {
		status := h.Cache.Ping(ctx)
		if status.Err() != nil {
			isError = true
			errorMessage = "cannot ping to cache service"
			return
		}
	}

	if h.MessageBroker != nil {
		if !h.MessageBroker.IsConnected() {
			isError = true
			errorMessage = "cannot ping to Message Broker service"
			return
		}
	}

	if h.Module.GrpcServer() != nil {

		//healthcheck
		// grpc healthcheck
		// Create gRPC client for health check
		healthConn, err := grpc.Dial(*appcore_model.IP+":"+*appcore_model.GrpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Println(1, err.Error())
			isError = true
			errorMessage = "cannot ping to GRPC service"
			return
		}
		healthClient := grpc_health_v1.NewHealthClient(healthConn)

		// Check gRPC server health status
		healthStatus, err := healthClient.Check(c, &grpc_health_v1.HealthCheckRequest{})
		if err != nil {
			log.Println(2, err.Error())
			isError = true
			errorMessage = "cannot ping to GRPC service"
			return
		}
		if healthStatus.Status != grpc_health_v1.HealthCheckResponse_SERVING {
			log.Println(3, err.Error())
			isError = true
			errorMessage = "cannot ping to GRPC service"
			return
		}
	}
}
