package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

import (
	"encoding/xml"
)

func AuditlogReader(conf ToolConfiguration, input <-chan UnifiedExportFormat, output chan<- UnifiedExportFormat) {
	defer wgAuditers.Done()

	for entry := range input {

		if entry.XIMessage.MessageKey == "" {
			fmt.Printf("WARNING: AuditlogReader received empty MessageKey!\n")
			continue
		}

		requestTemplate := fmt.Sprintf(`<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:urn="urn:AdapterMessageMonitoringVi" xmlns:urn1="urn:com.sap.aii.mdt.server.adapterframework.ws">
   <soapenv:Header/>
   <soapenv:Body>
      <urn:getLogEntries>
         <urn:messageKey>%s</urn:messageKey>
         <urn:archive>false</urn:archive>
         <!--Optional:-->
         <urn:maxResults>%d</urn:maxResults>
           <urn:locale>
            <!--Optional:-->
            <urn1:language>en</urn1:language>
       </urn:locale>
       <!--Optional:-->
         <urn:olderThan>3000-01-01T00:00:00+00:00</urn:olderThan>
      </urn:getLogEntries>
   </soapenv:Body>
</soapenv:Envelope>`, entry.XIMessage.MessageKey, conf.MaxAuditlinesPerSearch)

		req, err := http.NewRequest("POST", fmt.Sprintf("%s/AdapterMessageMonitoring/basic?style=document", conf.Hostname), strings.NewReader(requestTemplate))
		req.SetBasicAuth(conf.Username, conf.Password)
		req.Header.Set("Content-Type", "text/xml; charset=utf-8")
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

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

		responseBytes, err := io.ReadAll(resp.Body)

		searchResults := new(XIEnvelop)
		err = xml.Unmarshal(responseBytes, &searchResults)

		if err != nil {
			fmt.Printf("Please verify that host [%s], username [%s] and password are correct\n", conf.Hostname, conf.Username)
			fmt.Printf("HTTP call returned: %s\n", err.Error())
			os.Exit(3)
		}

		results := searchResults.Body.GetLogEntriesResponse.Response.AuditLogEntryData

		for _, log := range results {
			output <- FormatAuditLogForExport(entry, log)
		}

	}
}

func FormatAuditLogForExport(entry UnifiedExportFormat, log XIAuditLogEntryData) UnifiedExportFormat {
	export := ExportFormatAuditLog{
		MessageID:         entry.XIMessage.MessageID,
		Timestamp:         log.Timestamp,
		LocalizedText:     log.LocalizedText,
		TextKey:           log.TextKey,
		TextKeyParams:     make(map[string]string),
		Status:            decodeStatus(log.Status),
		TextLengthCutFrom: 0,
	}

	originalLength := 0
	cutDone := false

	for i, param := range log.Params.Strings {

		originalLength += len(param)
		if runtimeConfig.NoTrimForAuditLogs == false && len(param) > 4096 {
			// protection against very long audit log entries
			cutDone = true
			param = param[0:4096]
		}

		export.TextKeyParams[strconv.Itoa(i)] = param
	}

	originalLength += len(export.TextKey)
	if export.TextKey == export.LocalizedText {
		export.TextKey = ""
	}
	if runtimeConfig.NoTrimForAuditLogs == false && len(export.TextKey) > 4096 {
		// protection against very long audit log entries
		cutDone = true
		export.TextKey = export.TextKey[0:4096]
	}

	originalLength += len(export.LocalizedText)
	if runtimeConfig.NoTrimForAuditLogs == false && len(export.LocalizedText) > 4096 {
		// protection against very long audit log entries
		cutDone = true
		export.LocalizedText = export.LocalizedText[0:4096]
	}

	if cutDone {
		export.TextLengthCutFrom = originalLength
	}

	entry.AuditLog = &export

	return entry
}

func decodeStatus(in string) string {

	switch in {
	case "W":
		return "WARNING"
	case "E":
		return "ERROR"
	case "S":
		return "SUCCESS"
	default:
		return in
	}
}
