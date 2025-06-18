/*
   Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
   CreatedAt: 28 Mar 2024*/

package communications

var communicationsInstance ICommunications

func setCommunicationsInstance(commInstance ICommunications) {
	communicationsInstance = commInstance
}

func getCommunicationsInstance() ICommunications {
	return communicationsInstance
}

type communications struct {
	appName                string
	communicationTemplates *CommunicationTemplates

	sseEventIdMap            *CommEventIDSSEEventMap
	sseEventDestinationTopic string
	sseEventBuilder          CommSSEEventBuilder

	globalAlertEventIdMap       *CommEventIDGlobalAlertMap
	globalAlertDestinationTopic string
	globalAlertBuilder          CommGlobalAlertBuilder

	auditLogEventIdMap       *CommEventIDAuditLogMap
	auditLogDestinationTopic string
	auditLogBuilder          CommAuditLogBuilder
}

func NewCommunications(
	appName string,
	communicationTemplates *CommunicationTemplates,

	sseEventIdMap *CommEventIDSSEEventMap,
	sseEventDestinationTopic string,
	sseEventBuilder CommSSEEventBuilder,

	globalAlertEventIdMap *CommEventIDGlobalAlertMap,
	globalAlertDestinationTopic string,
	globalAlertBuilder CommGlobalAlertBuilder,

	auditLogEventIdMap *CommEventIDAuditLogMap,
	auditLogDestinationTopic string,
	auditLogBuilder CommAuditLogBuilder,
) ICommunications {
	var commInstance = &communications{
		appName:                appName,
		communicationTemplates: communicationTemplates,

		sseEventIdMap:            sseEventIdMap,
		sseEventDestinationTopic: sseEventDestinationTopic,
		sseEventBuilder:          sseEventBuilder,

		globalAlertEventIdMap:       globalAlertEventIdMap,
		globalAlertDestinationTopic: globalAlertDestinationTopic,
		globalAlertBuilder:          globalAlertBuilder,

		auditLogEventIdMap:       auditLogEventIdMap,
		auditLogDestinationTopic: auditLogDestinationTopic,
		auditLogBuilder:          auditLogBuilder,
	}
	setCommunicationsInstance(commInstance)
	return commInstance

}

type ICommunications interface {
	getAppName() string
	getCommunicationTemplates() *CommunicationTemplates

	getSSEEventIdMap() *CommEventIDSSEEventMap
	getSSEEventDestination() string
	getSSEEventBuilder() CommSSEEventBuilder

	getGlobalAlertEventIdMap() *CommEventIDGlobalAlertMap
	getGlobalAlertDestination() string
	getGlobalAlertBuilder() CommGlobalAlertBuilder

	getAuditLogEventIdMap() *CommEventIDAuditLogMap
	getAuditLogDestination() string
	getAuditLogBuilder() CommAuditLogBuilder
}

func (cem *communications) getAppName() string {
	return cem.appName
}
func (cem *communications) getCommunicationTemplates() *CommunicationTemplates {
	return cem.communicationTemplates
}

func (cem *communications) getSSEEventIdMap() *CommEventIDSSEEventMap {
	return cem.sseEventIdMap
}
func (cem *communications) getSSEEventDestination() string {
	return cem.sseEventDestinationTopic
}
func (cem *communications) getSSEEventBuilder() CommSSEEventBuilder {
	return cem.sseEventBuilder
}

func (cem *communications) getGlobalAlertEventIdMap() *CommEventIDGlobalAlertMap {
	return cem.globalAlertEventIdMap
}
func (cem *communications) getGlobalAlertDestination() string {
	return cem.globalAlertDestinationTopic
}
func (cem *communications) getGlobalAlertBuilder() CommGlobalAlertBuilder {
	return cem.globalAlertBuilder
}

func (cem *communications) getAuditLogDestination() string {
	return cem.auditLogDestinationTopic
}
func (cem *communications) getAuditLogEventIdMap() *CommEventIDAuditLogMap {
	return cem.auditLogEventIdMap
}
func (cem *communications) getAuditLogBuilder() CommAuditLogBuilder {
	return cem.auditLogBuilder
}
