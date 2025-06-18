/*
   Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
   CreatedAt: 28 Mar 2024*/

package web

import "github.com/gin-gonic/gin"

func SendResponse(ctx *gin.Context, response *JSONResponse, errorInterface *ErrorStruct) {
	responseCode := buildResponseCode(response, errorInterface)
	responseBody := buildResponseBody(response, errorInterface)
	ctx.JSON(responseCode, responseBody)
}

func BuildResponse(ctx *gin.Context, response *JSONResponse, errorInterface *ErrorStruct) (int, interface{}) {
	return buildResponseCode(response, errorInterface), buildResponseBody(response, errorInterface)
}

func buildResponseCode(data *JSONResponse, responseErr *ErrorStruct) int {
	if responseErr == nil {
		return data.GetStatusCode()
	}
	return responseErr.GetStatusCode()
}

func buildResponseBody(data *JSONResponse, responseErr *ErrorStruct) interface{} {
	if responseErr == nil {
		return data.GetResponseBody()
	}
	return BuildErrorResponse(responseErr)
}

func BuildErrorResponse(responseErr *ErrorStruct) interface{} {
	return map[string]interface{}{
		"errorCode":    responseErr.GetErrorCode(),
		"errorMessage": responseErr.GetErrorMessage(),
		"errorData":    responseErr.GetErrorData(),
	}
}
