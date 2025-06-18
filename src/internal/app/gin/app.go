/*
	Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
	CreatedAt: 28 Mar 2024*/

package app_gin

import (
	"context"
	GLogger "lruCache/poc/lib/logger"
	"lruCache/poc/src/config"
	apigin "lruCache/poc/src/internal/api/gin"
	ginmiddleware "lruCache/poc/src/internal/middleware/gin"

	"github.com/gin-gonic/gin"
)

func SetupServer(ctx context.Context, log *GLogger.LoggerService) *gin.Engine {
	log.Info(ctx, "setting_up_routes...")
	gin.SetMode(config.GinMode)
	engine := gin.New()
	engine.HandleMethodNotAllowed = true
	engine.Use(ginmiddleware.CorsMiddleware())
	engine.Use(ginmiddleware.OtelMiddleware(config.ServiceName))
	engine.Use(gin.Recovery())
	//engine.Use(ginmiddleware.CommunicationsMiddleware)
	baseRouterGroup := &apigin.RouterGroup{Router: engine.Group(config.BaseRouterPath), Log: log}
	baseRouterGroup.DefaultRoutes()
	baseRouterGroup.AddV1Routes()
	return engine
}
