//go:build beta
// +build beta

/*
   Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
   CreatedAt: 28 Mar 2024*/

package config

import (
	"fmt"
	GLogger "lruCache/poc/lib/logger"
	"lruCache/poc/src/constants"
	"lruCache/poc/src/modules/database"

	"github.com/gin-gonic/gin"
)

const AppEnv = constants.EnvBeta
const GinMode = gin.DebugMode
const LogLevel = GLogger.DEBUG
const DBType = database.MONGO
const MaxWorkerPoolSizeBulkMigrationsProcessor = 150
const MaxWorkerChanBufferSizeMigrationProcessor = 10
const (
	TestReferralCode constants.ReferralCode = "pawananjanikumar"
)

var BaseRouterPath = fmt.Sprintf("/%s/%s/%s", string(constants.BasePrefix), AppEnv, string(constants.ModuleBasePrefix))

/*
	START :: HARD INTERVENTION CHECKS
*/

/*
	END :: HARD INTERVENTION CHECKS
*/
