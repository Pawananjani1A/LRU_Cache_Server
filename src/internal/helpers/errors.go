/*
   Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
   CreatedAt: 28 Mar 2024*/

package helpers

import (
	"errors"
	"lruCache/poc/src/constants"
	"lruCache/poc/src/internal/web"
	"net/http"
)

func HandleError(err error) *web.ErrorStruct {
	var errStruct *web.ErrorStruct
	if errors.As(err, &errStruct) {
		return errStruct
	}
	errInterface := constants.ErrorsMap[err.Error()]
	if errInterface == nil {
		errInterface = web.NewHTTPError(http.StatusInternalServerError, "unhandled_error", err.Error(), nil)
	}
	return errInterface
}
