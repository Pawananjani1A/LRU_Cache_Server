/*
   Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
   CreatedAt: 28 Mar 2024*/

package constants

import "lruCache/poc/src/modules/communications"

const (
	CommEventNameAddedToGroupSuccess communications.CommunicationEventID = "COMM_EVENT_ADDED_TO_GROUP_SUCCESS"
	CommEventNameAddedToGroupFailure communications.CommunicationEventID = "COMM_EVENT_ADDED_TO_GROUP_FAILURE"
)

var CommunicationTemplates = communications.CommunicationTemplates{
	CommEventNameAddedToGroupSuccess: []communications.CommunicationEventType{communications.CommEventTypeSSE},
	CommEventNameAddedToGroupFailure: []communications.CommunicationEventType{communications.CommEventTypeSSE},
}

var CommSSEEventNameMap = communications.CommEventIDSSEEventMap{
	CommEventNameAddedToGroupSuccess: "sse_GroupSpends_UserModule_AddedToGroupSuccess",
	CommEventNameAddedToGroupFailure: "sse_GroupSpends_UserModule_AddedToGroupFailure",
}
