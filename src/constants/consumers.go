/*
   Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
   CreatedAt: 28 Mar 2024*/

package constants

type ConsumerEventType string
type ConsumerEventStatus string

const (
	CEUserCategoryChange ConsumerEventType = "USER_CATEGORY_CHANGE"
)

const (
	CEStatusNotStarted ConsumerEventStatus = "NOT_STARTED"
	CEStatusInProgress ConsumerEventStatus = "IN_PROGRESS"
	CEStatusFailure    ConsumerEventStatus = "FAILURE"
	CEStatusSuccess    ConsumerEventStatus = "SUCCESS"
)

type ConsumerEventFailureCode string

const (
	CEFailureUserNotFound ConsumerEventFailureCode = "X2_CORE_USER_NOT_FOUND"
)

var ConsumerEventFailureMessageMap = map[ConsumerEventFailureCode]string{
	CEFailureUserNotFound: "user not found against appuserid",
}
