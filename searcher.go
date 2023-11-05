package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

import (
	"encoding/xml"
)

type SearchResults struct {
	AdapterFrameworkData XIAdapterMessage
	Metadata             ExportMetadata
}

func MessageSearcher(conf ToolConfiguration, input <-chan XIOverviewDetailsPeriodEntry, uds chan<- SearchResults) {
	defer wgSearchers.Done()

	/////////////////////
	/// step 1: http request

	for entry := range input {

		EndTime := entry.End
		found := int32(0)
		atomic.AddInt32(&statistics.MaxXIMessageLines, entry.Total)

		for {
			// specially for continuation
			response := searchMessages(conf, entry.Start, entry.Interface, entry.ReceiverComponent, entry.SenderComponent, EndTime)

			for _, message := range response.Response.List.AdapterFrameworkData {
				uds <- SearchResults{
					AdapterFrameworkData: message,
					Metadata: ExportMetadata{
						Component:   entry.Component,
						Period:      entry.Period,
						PeriodStart: entry.Start,
						PeriodEnd:   entry.End,
						Extracted:   time.Now(),
					},
				}
				found++
			}

			EndTime = response.Response.ContinuationDate
			if EndTime == "" {
				// break out of the loop
				break
			}
			if runtimeConfig.BeVerbose {
				fmt.Println("Looping over messages with continuation date [", EndTime, "] for ", entry.Interface, entry.ReceiverComponent, entry.SenderComponent)
			}
		}

		if found < entry.Total {
			if runtimeConfig.BeVerbose {
				fmt.Printf("Found [%d] out of expected [%d] messages\n", found, entry.Total)
			}
		}
	}
}

func searchMessages(conf ToolConfiguration, Start string, Interface string, ReceiverComponent string, SenderComponent string, EndTime string) XIgetMessageListResponse {
	requestTemplate := fmt.Sprintf(`<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:urn="urn:AdapterMessageMonitoringVi" xmlns:urn1="urn:com.sap.aii.mdt.server.adapterframework.ws" xmlns:urn2="urn:com.sap.aii.mdt.api.data" xmlns:lang="urn:java/lang">
	   <soapenv:Header/>
	   <soapenv:Body>
	      <urn:getMessageList>
	         <urn:filter>
	           <urn1:archive>false</urn1:archive>
	           <urn1:dateType>0</urn1:dateType>
	           <urn1:fromTime>%s</urn1:fromTime>
	           <urn1:interface>
	              <urn2:name>%s</urn2:name>
	            </urn1:interface>
	           <urn1:nodeId>0</urn1:nodeId>
	            <urn1:onlyFaultyMessages>false</urn1:onlyFaultyMessages>
				<urn1:receiverName>%s</urn1:receiverName>
	           <urn1:retries>0</urn1:retries>
	            <urn1:retryInterval>0</urn1:retryInterval>
				<urn1:senderName>%s</urn1:senderName>
	           <urn1:timesFailed>0</urn1:timesFailed>
	           <urn1:toTime>%s</urn1:toTime>
	           <urn1:wasEdited>false</urn1:wasEdited>
               <urn1:returnLogLocations>true</urn1:returnLogLocations>
			   <urn1:onlyLogLocationsWithPayload>true</urn1:onlyLogLocationsWithPayload>	            
	        </urn:filter>
	         <urn:maxMessages>%d</urn:maxMessages>
	      </urn:getMessageList>
	   </soapenv:Body>
	</soapenv:Envelope>`, Start, Interface, ReceiverComponent, SenderComponent, EndTime, conf.MaxMessagesPerSearch)

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

	response := searchResults.Body.GetMessageListResponse

	return response

}
