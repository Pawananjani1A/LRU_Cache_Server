/*
	Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
	CreatedAt: 28 Mar 2024*/

package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const success = "SUCCESS"
const created = "CREATED"
const internal_server_error = "INTERNAL_SERVER_ERROR"

type JSONResponse struct {
	ctx          *gin.Context
	statusCode   int
	responseBody *struct {
		Code    string      `json:"code,omitempty"`
		Message string      `json:"message,omitempty"`
		Data    interface{} `json:"data,omitempty"`
	}
}

func (j *JSONResponse) GetStatusCode() int {
	return j.statusCode
}

func (j *JSONResponse) GetResponseBody() interface{} {
	if j.responseBody == nil {
		return map[string]interface{}{}
	}
	return j.responseBody
}

func NewHTTPSuccessResponse(ctx *gin.Context, message string, jsonData interface{}) *JSONResponse {
	return NewJSONResponse(ctx, http.StatusOK, success, message, jsonData)
}

func NewHTTPFailureResponse(ctx *gin.Context, message string, jsonData interface{}) *JSONResponse {
	return NewJSONResponse(ctx, http.StatusInternalServerError, internal_server_error, message, jsonData)
}

func NewHTTPCreatedResponse(ctx *gin.Context, message string, jsonData interface{}) *JSONResponse {
	return NewJSONResponse(ctx, http.StatusCreated, created, message, jsonData)
}

func NewHTTPAcceptedResponse(ctx *gin.Context, message string, jsonData interface{}) *JSONResponse {
	return NewJSONResponse(ctx, http.StatusAccepted, success, message, jsonData)
}

func NewJSONResponse(ctx *gin.Context, statusCode int, jsonCode string, jsonMessage string, jsonData interface{}) *JSONResponse {
	return &JSONResponse{
		ctx:        ctx,
		statusCode: statusCode,
		responseBody: &struct {
			Code    string      `json:"code,omitempty"`
			Message string      `json:"message,omitempty"`
			Data    interface{} `json:"data,omitempty"`
		}{Code: jsonCode, Message: jsonMessage, Data: jsonData},
	}
}

func (j *JSONResponse) JSON() {
	j.ctx.JSON(j.statusCode, j.responseBody)
}
