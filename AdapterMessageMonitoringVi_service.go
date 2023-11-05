package main

type XIEnvelop struct {
	Body Body `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
}

type Body struct {
	GetMessageListResponse   XIgetMessageListResponse                 `xml:"urn:AdapterMessageMonitoringVi getMessageListResponse"`
	GetLogEntriesResponse    XIgetLogEntriesResponse                  `xml:"urn:AdapterMessageMonitoringVi getLogEntriesResponse"`
	GetUDSAttributesResponse XIgetUserDefinedSearchAttributesResponse `xml:"urn:AdapterMessageMonitoringVi getUserDefinedSearchAttributesResponse"`
}

type XIgetMessageListResponse struct {
	Response struct {
		ContinuationDate string `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws date"`
		List             struct {
			AdapterFrameworkData []XIAdapterMessage `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws AdapterFrameworkData"`
		} `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws list"`
		Warning bool `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws warning"`
	} `xml:"urn:AdapterMessageMonitoringVi Response"`
}

type XIgetLogEntriesResponse struct {
	Response struct {
		AuditLogEntryData []XIAuditLogEntryData `xml:"urn:com.sap.aii.mdt.api.data AuditLogEntryData"`
	} `xml:"urn:AdapterMessageMonitoringVi Response"`
}

type XIgetUserDefinedSearchAttributesResponse struct {
	Response struct {
		BusinessAttributes []XIBusinessAttribute `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws BusinessAttribute"`
	} `xml:"urn:AdapterMessageMonitoringVi Response"`
}

type XIAuditLogEntryData struct {
	Timestamp string `xml:"urn:com.sap.aii.mdt.api.data timeStamp"`
	TextKey   string `xml:"urn:com.sap.aii.mdt.api.data textKey"`
	Params    struct {
		Strings []string `xml:"urn:java/lang String"`
	} `xml:"urn:com.sap.aii.mdt.api.data params"`
	Status        string `xml:"urn:com.sap.aii.mdt.api.data status"`
	LocalizedText string `xml:"urn:com.sap.aii.mdt.api.data localizedText"`
}

type XIAdapterMessage struct {
	ApplicationComponent string       `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws applicationComponent"`
	BusinessMessage      *XIBoolean   `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws businessMessage"`
	Cancelable           *XIBoolean   `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws cancelable"`
	ConnectionName       string       `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws connectionName"`
	Credential           string       `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws credential"`
	Direction            string       `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws direction"`
	Editable             *XIBoolean   `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws editable"`
	EndTime              string       `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws endTime"`
	Endpoint             string       `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws endpoint"`
	ErrorCategory        string       `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws errorCategory"`
	ErrorCode            string       `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws errorCode"`
	Headers              string       `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws headers"`
	Interface            *XIInterface `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws interface"`
	IsPersistent         bool         `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws isPersistent"`
	MessageID            string       `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws messageID"`
	MessageKey           string       `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws messageKey"`
	MessageType          string       `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws messageType"`
	NodeId               int          `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws nodeId"`
	PersistUntil         string       `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws persistUntil"`
	Protocol             string       `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws protocol"`
	QualityOfService     string       `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws qualityOfService"`
	// ReceiverInterface    *XIInterface    `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws receiverInterface"`
	// sadly: SAP PO does not store this data...
	ReceiverName  string          `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws receiverName"`
	ReceiverParty *XIMessageParty `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws receiverParty"`
	ReferenceID   string          `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws referenceID"`
	Restartable   *XIBoolean      `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws restartable"`
	Retries       int             `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws retries"`
	RetryInterval int             `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws retryInterval"`
	ScheduleTime  string          `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws scheduleTime"`
	// SenderInterface      *XIInterface    `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws senderInterface"`
	// sadly: SAP PO does not store this data...
	SenderName           string          `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws senderName"`
	SenderParty          *XIMessageParty `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws senderParty"`
	SequenceNumber       int             `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws sequenceNumber"`
	SerializationContext string          `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws serializationContext"`
	ServiceDefinition    string          `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws serviceDefinition"`
	SoftwareComponent    string          `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws SoftwareComponent"`
	StartTime            string          `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws startTime"`
	Status               string          `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws status"`
	TimesFailed          int             `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws timesFailed"`
	Transport            string          `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws transport"`
	ValidUntil           string          `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws validUntil"`
	Version              string          `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws version"`
	WasEdited            *XIBoolean      `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws wasEdited"`
	// BusinessAttributes // skip this -- very strange, might be "Content data"
	// payloadPermissionWarning // skip this -- uninteresting
	ErrorLabel         int            `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws errorLabel"`
	ScenarioIdentifier string         `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws scenarioIdentifier"`
	ParentID           string         `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws parentID"`
	Duration           *XIDuration    `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws duration"`
	Size               int            `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws size"`
	MessagePriority    int            `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws messagePriority"`
	RootID             string         `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws rootID"`
	SequenceID         string         `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws sequenceID"`
	Passport           string         `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws passport"`
	PassportTID        string         `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws passportTID"`
	LogLocations       XILogLocations `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws logLocations"`
}

type XIBusinessAttribute struct {
	Name  string `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws name"`
	Value string `xml:"urn:com.sap.aii.mdt.server.adapterframework.ws value"`
}

type XIInterface struct {
	Name              string `xml:"urn:com.sap.aii.mdt.api.data name"`
	Namespace         string `xml:"urn:com.sap.aii.mdt.api.data namespace"`
	SenderParty       string `xml:"urn:com.sap.aii.mdt.api.data senderParty"`
	SenderComponent   string `xml:"urn:com.sap.aii.mdt.api.data senderComponent"`
	ReceiverParty     string `xml:"urn:com.sap.aii.mdt.api.data receiverParty"`
	ReceiverComponent string `xml:"urn:com.sap.aii.mdt.api.data receiverComponent"`
}

type XIBoolean struct {
	Value bool `xml:"urn:com.sap.aii.mdt.api.data value"`
}

type XIMessageParty struct {
	Agency string `xml:"urn:com.sap.aii.mdt.api.data agency"`
	Name   string `xml:"urn:com.sap.aii.mdt.api.data name"`
	Schema string `xml:"urn:com.sap.aii.mdt.api.data schema"`
}

type XIDuration struct {
	Value int64 `xml:"urn:com.sap.aii.mdt.api.data duration"`
}

type XILogLocations struct {
	String []string `xml:"urn:java/lang String"`
}
