/*
   Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
   CreatedAt: 28 Mar 2024*/

package constants

type GroupRequestType string

type GroupRequestFailureCode string

const (
	GroupRequestValidationFailed GroupRequestFailureCode = "VALIDATION_FAILED"
)

var GroupFailureReasonMap = map[GroupRequestFailureCode]string{
	GroupRequestValidationFailed: "group request validation failed",
}
