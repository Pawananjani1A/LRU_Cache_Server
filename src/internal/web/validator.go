/*
   Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
   CreatedAt: 28 Mar 2024*/

package web

import (
	"github.com/gin-gonic/gin"
)

const (
	regExpEmail  string = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regExpMobile string = `^[6-9]\d{9}$`
	jsonBody     string = "jsonBody"
)

func GetJsonBody[BodyType any](ctx *gin.Context) BodyType {
	return ctx.MustGet(jsonBody).(BodyType)
}

func ValidateJsonBody[BodyType any](ctx *gin.Context) *ErrorStruct {
	var body BodyType
	var errS *ErrorStruct
	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		if err.Error() == "EOF" {
			if ctx.Request.Method != "GET" {
				errS = NewHTTPBadRequestError("no input received in payload", map[string]interface{}{})
			}
		} else {
			errS = NewHTTPBadRequestError(err.Error(), map[string]interface{}{})
		}
	}
	ctx.Set(jsonBody, body)
	return errS
}
