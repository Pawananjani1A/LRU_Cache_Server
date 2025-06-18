/*
   Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
   CreatedAt: 28 Mar 2024*/

package communications

type CommSSEEventData struct {
	TimeStamp      int64       `json:"ts,omitempty"`
	EventName      string      `json:"name,omitempty"`
	PayloadVersion string      `json:"payloadVersion,omitempty"`
	EventData      interface{} `json:"data,omitempty"`
}

type CommSSEEventPayload struct {
	PackageName    string             `json:"packageName,omitempty"`
	AppUserID      string             `json:"appUserId,omitempty"`
	PayloadVersion string             `json:"payloadVersion,omitempty"`
	Profile        string             `json:"profile,omitempty"`
	Events         []CommSSEEventData `json:"events,omitempty"`
}
