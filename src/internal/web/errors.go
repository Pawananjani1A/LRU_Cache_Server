/*
	Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
	CreatedAt: 28 Mar 2024*/

package web

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorStruct struct {
	ctx          *gin.Context
	statusCode   int
	errorData    map[string]interface{}
	errorCode    string
	errorMessage string
}

func (e *ErrorStruct) Error() string {
	// implemented so that error can be used as golang error
	jm, _ := json.Marshal(map[string]interface{}{
		"errorCode":    e.errorCode,
		"errorData":    e.errorData,
		"errorMessage": e.errorMessage,
	})
	return string(jm)
}

func (e *ErrorStruct) GetStatusCode() int {
	return e.statusCode
}

func (e *ErrorStruct) GetErrorData() map[string]interface{} {
	if e.errorData == nil {
		return map[string]interface{}{}
	}
	return e.errorData
}

func (e *ErrorStruct) GetErrorCode() string {
	return e.errorCode
}

func (e *ErrorStruct) GetErrorMessage() string {
	return e.errorMessage
}

func NewHTTPBadRequestError(errorMessage string, errorData map[string]interface{}) *ErrorStruct {
	return NewHTTPError(http.StatusBadRequest, "BAD_REQUEST", errorMessage, errorData)
}

func NewHTTPNotFoundError(errorCode string, errorMessage string, errorData map[string]interface{}) *ErrorStruct {
	return NewHTTPError(http.StatusNotFound, errorCode, errorMessage, errorData)
}

func NewHTTPInternalServerError(errorCode string, errorMessage string, errorData map[string]interface{}) *ErrorStruct {
	return NewHTTPError(http.StatusInternalServerError, errorCode, errorMessage, errorData)
}

func NewHTTPUnprocessableEntityError(errorCode string, errorMessage string, errorData map[string]interface{}) *ErrorStruct {
	return NewHTTPError(http.StatusUnprocessableEntity, errorCode, errorMessage, errorData)
}

func NewHTTPFailedDependencyError(errorCode string, errorMessage string, errorData map[string]interface{}) *ErrorStruct {
	return NewHTTPError(http.StatusFailedDependency, errorCode, errorMessage, errorData)
}

func NewHTTPConflictError(errorCode string, errorMessage string, errorData map[string]interface{}) *ErrorStruct {
	return NewHTTPError(http.StatusConflict, errorCode, errorMessage, errorData)
}

func NewInvalidQueryParamsError(errorMessage string, errorData map[string]interface{}) *ErrorStruct {
	return NewHTTPError(http.StatusBadRequest, "INVALID_QUERY_PARAMS", errorMessage, errorData)
}

func NewHTTPError(statusCode int, errorCode string, errorMessage string, errorData map[string]interface{}) *ErrorStruct {
	return &ErrorStruct{statusCode: statusCode, errorCode: errorCode, errorMessage: errorMessage, errorData: errorData}
}

func NewError(errorCode string, errorMessage string, errorData map[string]interface{}) *ErrorStruct {
	return &ErrorStruct{errorCode: errorCode, errorMessage: errorMessage, errorData: errorData}
}
