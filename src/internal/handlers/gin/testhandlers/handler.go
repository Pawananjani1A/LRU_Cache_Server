/*
   Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
   CreatedAt: 28 Mar 2024*/

package testhandlers

import (
	"lruCache/poc/src/internal/web"

	"github.com/gin-gonic/gin"
)

func CheckRouterPath(ctx *gin.Context) (*web.JSONResponse, *web.ErrorStruct) {
	return web.NewHTTPSuccessResponse(ctx, "successfully called CheckRouterPath", map[string]interface{}{}), nil
}
