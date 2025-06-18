/*
   Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
   CreatedAt: 28 Mar 2024*/

package repositorymodels

import "time"

type ConsumerEventModel struct {
	EventID      string      `json:"eventId,omitempty" bson:"eventId,omitempty"` // X2MPW
	AppUserID    string      `json:"appUserId,omitempty" bson:"appUserId,omitempty"`
	EventType    string      `json:"eventType,omitempty" bson:"eventType,omitempty"`
	Payload      interface{} `json:"payload,omitempty" bson:"payload,omitempty"`
	Status       string      `json:"status,omitempty" bson:"status,omitempty"`
	ErrorCode    *string     `json:"errorCode,omitempty" bson:"errorCode,omitempty"`
	ErrorMessage *string     `json:"errorMessage,omitempty" bson:"errorMessage,omitempty"`
	CreatedAt    *time.Time  `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	CreatedBy    *string     `json:"createdBy,omitempty" bson:"createdBy,omitempty"`
	UpdatedAt    *time.Time  `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	DocVersion   int32       `json:"_v,omitempty" bson:"_v,omitempty"`
}
