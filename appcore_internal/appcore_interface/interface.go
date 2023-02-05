package inf

import (
	"github.com/panapol-p/appcore/appcore_router"
	"github.com/panapol-p/appcore/appcore_utils"
	"gorm.io/gorm"
)

type Handler struct {
	Store          *gorm.DB
	Version        string
	CircuitBreaker appcore_utils.CircuitBreaker
}

type Module interface {
	ModuleAPI(r *appcore_router.Router)
}
