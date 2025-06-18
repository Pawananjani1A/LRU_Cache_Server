/*
   Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
   CreatedAt: 28 Mar 2024*/

package constants

import (
	"lruCache/poc/src/internal/web"
)

const (
	//ErrMaxCacheSizeLimitReached
	ErrMaxCacheSizeLimitReached = "MAX_CACHE_SIZE_LIMIT_REACHED"
)

var ErrorsMap = map[string]*web.ErrorStruct{
	ErrMaxCacheSizeLimitReached: web.NewHTTPConflictError(ErrMaxCacheSizeLimitReached, "maximum group size already reached. please create a new group.", nil),
}
