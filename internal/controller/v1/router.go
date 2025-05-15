package v1

import (
	"log/slog"
	"net/http"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/P1xart/effective_mobile_service/docs"
	"github.com/P1xart/effective_mobile_service/internal/controller/v1/middleware"
	"github.com/P1xart/effective_mobile_service/internal/service"

	"github.com/gin-gonic/gin"
)

func NewRouter(log *slog.Logger, router *gin.Engine, services *service.Services) {
	router.Use(middleware.CORS(log))
	router.Use(middleware.Log(log))

	router.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))


	v1 := router.Group("/v1")

	newHumanRoutes(log, v1.Group("/human"), services.Human)
}
