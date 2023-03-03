package appcore_router

import (
	"net/http"

	"github.com/sirupsen/logrus"
	requestid "github.com/sumit-tembe/gin-requestid"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Router struct {
	*gin.Engine

	logger *logrus.Logger
}

func New(logger *logrus.Logger) *Router {
	r := gin.New()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AddAllowHeaders("authorization")
	r.Use(cors.New(config))

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(requestid.RequestID(nil))

	/*r.Use(limit.NewRateLimiter(func(c *gin.Context) string {
		return c.ClientIP() // limit rate by client ip
	}, func(c *gin.Context) (*rate.Limiter, time.Duration) {
		return rate.NewLimiter(20, 1), 1 * time.Hour
	}, func(c *gin.Context) {
		c.AbortWithStatus(429) // handle exceed rate limit request
	}))*/
	return &Router{r, logger}
}

func (r *Router) ListenAndServe(ip, port string) *http.Server {
	r.logger.Info("listen api on ", ip+":"+port)
	s := &http.Server{
		Addr:    ip + ":" + port,
		Handler: r,
		// ReadTimeout:  10 * time.Second,
		// WriteTimeout: 10 * time.Second,
		// IdleTimeout:    1 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			r.logger.Fatalf("listen: %s\n", err)
		}
	}()
	return s
}
