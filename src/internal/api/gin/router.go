/*
	Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
	CreatedAt: 28 Mar 2024*/

package api_gin

import (
	GLogger "lruCache/poc/lib/logger"
	"lruCache/poc/src/constants"
	routemodels "lruCache/poc/src/internal/api/gin/models"
	"lruCache/poc/src/internal/handlers/gin/cachehandlers"
	"lruCache/poc/src/internal/handlers/gin/testhandlers"
	ginmiddlewares "lruCache/poc/src/internal/middleware/gin"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type RouterGroup struct {
	Router *gin.RouterGroup
	Log    *GLogger.LoggerService
}

func (rg *RouterGroup) DefaultRoutes() {
	defaultRoutes := rg.Router.Group(string(constants.RGDocs))
	defaultRoutes.GET("/", ginSwagger.WrapHandler(swaggerFiles.Handler))
	defaultRoutes.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"message": "successfully_called_health_endpoint",
		})
		ctx.Next()
	})
}

func (rg *RouterGroup) AddV1Routes() {
	v1Routes := rg.Router.Group(string(constants.RV1))
	internalRoutes := v1Routes.Group(string(constants.RGInternalGroup))
	internalRoutes.GET(string(constants.RGCheckRouterPath), ginmiddlewares.ServeEndpoint[routemodels.TestModel](testhandlers.CheckRouterPath))

	//Cache Routes
	cacheRoutes := v1Routes.Group(string(constants.RGCache))
	cacheRoutes.POST(string(constants.RPGet), ginmiddlewares.ServeEndpoint[routemodels.GetModel](cachehandlers.GetValue))
	cacheRoutes.POST(string(constants.RPSet), ginmiddlewares.ServeEndpoint[routemodels.SetModel](cachehandlers.SetValue))
}
