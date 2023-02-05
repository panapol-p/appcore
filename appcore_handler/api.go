package appcore_handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/memphisdev/memphis.go"
	inf "github.com/panapol-p/appcore/appcore_internal/appcore_interface"
	"gorm.io/gorm"
)

type ApiHandler struct {
	ServiceName string
	Version     string
	Module      inf.Module

	DB            *gorm.DB
	Cache         *redis.Client
	MessageBroker *memphis.Conn
}

func NewHandler(serviceName, version string, module inf.Module, db *gorm.DB, cache *redis.Client, messageBroker *memphis.Conn) *ApiHandler {
	return &ApiHandler{
		ServiceName:   serviceName,
		Version:       version,
		Module:        module,
		DB:            db,
		Cache:         cache,
		MessageBroker: messageBroker,
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

}
