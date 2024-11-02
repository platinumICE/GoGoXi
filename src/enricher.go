package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

import (
	"encoding/xml"
)

func udsEnricher(conf ToolConfiguration, input <-chan SearchResults, fileWriter chan<- UnifiedExportFormat, audit chan<- UnifiedExportFormat) {
	defer wgEnrichers.Done()

	/////////////////////
	/// step 1: http request

	for message := range input {

		if message.AdapterFrameworkData.MessageKey == "" {
			fmt.Printf("WARNING: udsEnricher received empty MessageKey!\n")
			continue
		}

		if conf.SkipUDS == false {

			requestTemplate := fmt.Sprintf(`<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:urn="urn:AdapterMessageMonitoringVi">
	   <soapenv:Header/>
	   <soapenv:Body>
	      <urn:getUserDefinedSearchAttributes>
	         <urn:messageKey>%s</urn:messageKey>
	         <urn:archive>false</urn:archive>
	      </urn:getUserDefinedSearchAttributes>
	   </soapenv:Body>
	</soapenv:Envelope>`, message.AdapterFrameworkData.MessageKey)

			req, err := http.NewRequest("POST", fmt.Sprintf("%s/AdapterMessageMonitoring/basic?style=document", conf.Hostname), strings.NewReader(requestTemplate))
			req.SetBasicAuth(conf.Username, conf.Password)
			req.Header.Set("Content-Type", "text/xml; charset=utf-8")
			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}

			switch resp.StatusCode {
			case 200:
				// noop
			case 401:
				fmt.Printf("HTTP 401: incorrect password for user %s\n", conf.Username)
				os.Exit(2)
			case 403:
				fmt.Printf("HTTP 403: incorrect password for user %s\n", conf.Username)
				os.Exit(2)
			default:
				fmt.Printf("HTTP %s: cannot read overview for host %s\n", resp.Status, conf.Hostname)
				os.Exit(3)
			}

			var buf bytes.Buffer
			_, _ = io.Copy(&buf, resp.Body)
			resp.Body.Close()

			searchResults := new(XIEnvelop)
			err = xml.Unmarshal(buf.Bytes(), &searchResults)
			buf.Reset()

			if err != nil {
				fmt.Printf("Please verify that host [%s], username [%s] and password are correct\n", conf.Hostname, conf.Username)
				fmt.Printf("HTTP call returned: %s\n", err.Error())
				os.Exit(3)
			}

			export_xi_message := FormatXIMessageForExport(conf, message, searchResults.Body.GetUDSAttributesResponse.Response.BusinessAttributes)

			audit <- export_xi_message
			fileWriter <- export_xi_message

		} else {
			emptyUDS := new([]XIBusinessAttribute)
			export_xi_message := FormatXIMessageForExport(conf, message, *emptyUDS)

			audit <- export_xi_message
			fileWriter <- export_xi_message
		}

	}
}

func FormatXIMessageForExport(conf ToolConfiguration, message SearchResults, uds []XIBusinessAttribute) UnifiedExportFormat {

	entry := message.AdapterFrameworkData
	export := ExportFormatXIMessage{
		ApplicationComponent: entry.ApplicationComponent,
		ConnectionName:       entry.ConnectionName,
		Credential:           entry.Credential,
		Direction:            entry.Direction,
		EndTime:              entry.EndTime,
		Endpoint:             entry.Endpoint,
		ErrorCategory:        entry.ErrorCategory,
		ErrorCode:            entry.ErrorCode,
		IsPersistent:         entry.IsPersistent,
		MessageID:            entry.MessageID,
		MessageKey:           entry.MessageKey,
		MessageType:          entry.MessageType,
		NodeId:               strconv.Itoa(entry.NodeId),
		PersistUntil:         entry.PersistUntil,
		Protocol:             entry.Protocol,
		QualityOfService:     entry.QualityOfService,
		ReceiverName:         entry.ReceiverName,
		ReferenceID:          entry.ReferenceID,
		Retries:              entry.Retries,
		RetryInterval:        normalizedDuration(entry.RetryInterval), // convert to seconds
		ScheduleTime:         entry.ScheduleTime,
		SenderName:           entry.SenderName,
		SequenceNumber:       entry.SequenceNumber,
		SerializationContext: entry.SerializationContext,
		ServiceDefinition:    entry.ServiceDefinition,
		SoftwareComponent:    entry.SoftwareComponent,
		StartTime:            entry.StartTime,
		Status:               entry.Status,
		TimesFailed:          entry.TimesFailed,
		Transport:            entry.Transport,
		ValidUntil:           entry.ValidUntil,
		Version:              entry.Version,
		ErrorLabel:           entry.ErrorLabel,
		ScenarioIdentifier:   entry.ScenarioIdentifier,
		ParentID:             entry.ParentID,
		Size:                 entry.Size,
		MessagePriority:      entry.MessagePriority,
		RootID:               entry.RootID,
		SequenceID:           entry.SequenceID,
		Passport:             entry.Passport,
		PassportTID:          entry.PassportTID,
		LogLocations:         entry.LogLocations.String,
		UDS:                  make(map[string][]string),
	}

	export.Headers.Original = entry.Headers
	if o := export.Headers.Original; strings.HasPrefix(export.Headers.Original, "content-length=") {
		o = strings.SplitN(o, "\n", 2)[0]
		contentLength, _ := strings.CutPrefix(o, "content-length=")
		export.Headers.ContentLength, _ = strconv.ParseInt(contentLength, 10, 32)
	}

	if entry.BusinessMessage != nil {
		export.BusinessMessage = &entry.BusinessMessage.Value
	}

	if entry.Cancelable != nil {
		export.Cancelable = &entry.Cancelable.Value
	}

	if entry.Editable != nil {
		export.Editable = &entry.Editable.Value
	}

	if entry.Interface != nil {
		export.Interface = &ExportXIInterface{
			Name:              entry.Interface.Name,
			Namespace:         entry.Interface.Namespace,
			SenderParty:       entry.Interface.SenderParty,
			SenderComponent:   entry.Interface.SenderComponent,
			ReceiverParty:     entry.Interface.ReceiverParty,
			ReceiverComponent: entry.Interface.ReceiverComponent,
		}
	}

	// sadly: SAP PO does not store this data
	// if entry.ReceiverInterface != nil {
	// 	export.ReceiverInterface = &ExportXIInterface{
	// 		Name:              entry.ReceiverInterface.Name,
	// 		Namespace:         entry.ReceiverInterface.Namespace,
	// 		SenderParty:       entry.ReceiverInterface.SenderParty,
	// 		SenderComponent:   entry.ReceiverInterface.SenderComponent,
	// 		ReceiverParty:     entry.ReceiverInterface.ReceiverParty,
	// 		ReceiverComponent: entry.ReceiverInterface.ReceiverComponent,
	// 	}
	// }

	if entry.ReceiverParty != nil {
		export.ReceiverParty = &ExportXIMessageParty{
			Agency: entry.ReceiverParty.Agency,
			Name:   entry.ReceiverParty.Name,
			Schema: entry.ReceiverParty.Schema,
		}
	}

	if entry.Restartable != nil {
		export.Restartable = &entry.Restartable.Value
	}

	// sadly: SAP PO does not store this data
	// if entry.SenderInterface != nil {
	// 	export.SenderInterface = &ExportXIInterface{
	// 		Name:              entry.SenderInterface.Name,
	// 		Namespace:         entry.SenderInterface.Namespace,
	// 		SenderParty:       entry.SenderInterface.SenderParty,
	// 		SenderComponent:   entry.SenderInterface.SenderComponent,
	// 		ReceiverParty:     entry.SenderInterface.ReceiverParty,
	// 		ReceiverComponent: entry.SenderInterface.ReceiverComponent,
	// 	}
	// }

	if entry.SenderParty != nil {
		export.SenderParty = &ExportXIMessageParty{
			Agency: entry.SenderParty.Agency,
			Name:   entry.SenderParty.Name,
			Schema: entry.SenderParty.Schema,
		}
	}

	if entry.WasEdited != nil {
		export.WasEdited = &entry.WasEdited.Value
	}
	if entry.Duration != nil {
		export.Duration = normalizedDuration(int(entry.Duration.Value))
	} else {
		export.Duration = 0
	}

	// special handling: clear zero persist date
	persistUntil, err := time.Parse("2006-01-02T15:04:05.000-07:00", export.PersistUntil)
	if err == nil {
		if persistUntil.Unix() == 0 {
			export.PersistUntil = ""
		}
	}

	for _, attr := range uds {
		old, exists := export.UDS[attr.Name]
		if exists {
			new := append(old, attr.Value)
			export.UDS[attr.Name] = new
		} else {
			new := []string{attr.Value}
			export.UDS[attr.Name] = new
		}
	}

	for k, _ := range export.UDS {
		export.UDSKeys = append(export.UDSKeys, k)
	}

	return UnifiedExportFormat{
		XIMessage: export,
		Metadata:  message.Metadata,
	}

}

func normalizedDuration(in int) float32 {
	return float32(in) / 1000.0
}
