package main

import "time"

type UnifiedExportFormat struct {
	XIMessage ExportFormatXIMessage `json:"ximessage"`
	AuditLog  *ExportFormatAuditLog `json:"auditlog,omitempty"`
	Metadata  ExportMetadata        `json:"gogoxi"`
}

type ExportMetadata struct {
	Component   string    `json:"component"`
	Period      string    `json:"-"`
	PeriodStart string    `json:"-"`
	PeriodEnd   string    `json:"-"`
	Extracted   time.Time `json:"extractedOn"`
}

type ExportXIInterface struct {
	Name              string `json:"name"`
	Namespace         string `json:"namespace"`
	SenderParty       string `json:"senderParty,omitempty"`
	SenderComponent   string `json:"senderComponent,omitempty"`
	ReceiverParty     string `json:"receiverParty,omitempty"`
	ReceiverComponent string `json:"receiverComponent,omitempty"`
}

type ExportXIMessageParty struct {
	Agency string `json:"agency"`
	Name   string `json:"name"`
	Schema string `json:"schema"`
}

type ExportFormatXIMessage struct {
	ApplicationComponent string `json:"applicationComponent,omitempty"`
	BusinessMessage      *bool  `json:"businessMessage,omitempty"`
	Cancelable           *bool  `json:"cancelable,omitempty"`
	ConnectionName       string `json:"connectionName"`
	Credential           string `json:"credential,omitempty"`
	Direction            string `json:"direction"`
	Editable             *bool  `json:"editable,omitempty"`
	EndTime              string `json:"endTime,omitempty"`
	Endpoint             string `json:"endpoint"`
	ErrorCategory        string `json:"errorCategory,omitempty"`
	ErrorCode            string `json:"errorCode,omitempty"`
	// little enhancement
	Headers struct {
		Original      string `json:"original"`
		ContentLength int64  `json:"contentLength"`
	} `json:"headers"`
	Interface            *ExportXIInterface    `json:"interface,omitempty"`
	IsPersistent         bool                  `json:"isPersistent"`
	MessageID            string                `json:"messageID"`
	MessageKey           string                `json:"messageKey"`
	MessageType          string                `json:"messageType"`
	NodeId               string                `json:"nodeId"` // math makes no sense for node id, so force string
	PersistUntil         string                `json:"persistUntil,omitempty"`
	Protocol             string                `json:"protocol"`
	QualityOfService     string                `json:"qualityOfService"`
	ReceiverInterface    *ExportXIInterface    `json:"receiverInterface,omitempty"`
	ReceiverName         string                `json:"receiverName"`
	ReceiverParty        *ExportXIMessageParty `json:"receiverParty,omitempty"`
	ReferenceID          string                `json:"referenceID"`
	Restartable          *bool                 `json:"restartable,omitempty"`
	Retries              int                   `json:"retries"`
	RetryInterval        float32               `json:"retryInterval"`
	ScheduleTime         string                `json:"scheduleTime,omitempty"`
	SenderInterface      *ExportXIInterface    `json:"senderInterface,omitempty"`
	SenderName           string                `json:"senderName"`
	SenderParty          *ExportXIMessageParty `json:"senderParty,omitempty"`
	SequenceNumber       int                   `json:"sequenceNumber"`
	SerializationContext string                `json:"serializationContext"`
	ServiceDefinition    string                `json:"serviceDefinition"`
	SoftwareComponent    string                `json:"softwareComponent"`
	StartTime            string                `json:"startTime"`
	Status               string                `json:"status"`
	TimesFailed          int                   `json:"timesFailed"`
	Transport            string                `json:"transport"`
	ValidUntil           string                `json:"validUntil,omitempty"`
	Version              string                `json:"version"`
	WasEdited            *bool                 `json:"wasEdited,omitempty"`
	// BusinessAttributes
	// payloadPermissionWarning
	ErrorLabel         int                 `json:"errorLabel,omitempty"`
	ScenarioIdentifier string              `json:"scenarioIdentifier"`
	ParentID           string              `json:"parentID"`
	Duration           float32             `json:"duration"`
	Size               int                 `json:"size"`
	MessagePriority    int                 `json:"messagePriority"`
	RootID             string              `json:"rootID"`
	SequenceID         string              `json:"sequenceID"`
	Passport           string              `json:"-"`
	PassportTID        string              `json:"passportTID"`
	LogLocations       []string            `json:"logLocations"`
	UDS                map[string][]string `json:"UDS,omitempty"`
	UDSKeys            []string            `json:"UDSKeys,omitempty"`
}

type ExportFormatAuditLog struct {
	MessageID         string            `json:"messageID"`
	Timestamp         string            `json:"timestamp"`
	Status            string            `json:"status"`
	LocalizedText     string            `json:"localizedText"`
	TextKey           string            `json:"textKey,omitempty"`
	TextKeyParams     map[string]string `json:"textKeyParams,omitempty"`
	TextLengthCutFrom int               `json:"textCutFrom,omitempty"`
}
