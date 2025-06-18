/*
   Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
   CreatedAt: 28 Mar 2024*/

package gin

import (
	"fmt"
	"lruCache/poc/src/config"
	"lruCache/poc/src/constants"
	"strings"

	"github.com/gin-gonic/gin"
)

func CommunicationsMiddleware(c *gin.Context) {
	//c.Next()
	fullPath := c.FullPath()
	baseRoutePath := strings.Replace(fullPath, config.BaseRouterPath, "", -1)
	splits := strings.SplitN(baseRoutePath, "/", 4)
	fmt.Println(splits)
	routeGroup := splits[2]
	switch routeGroup {
	case string(constants.RGGroup):
		routePath := splits[3]
		switch routePath {
		case string(constants.RPUpdateGroup):
			break
		}

	}
}
