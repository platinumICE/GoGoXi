# GoGoXi
GoGoXi SAP PO message and audit log extraction for ELK

## Description

This tool allows for extraction of XI messages and audit log entries from Adapter Framework of SAP Process Orcherstration using HTTP API calls only. All messages which fall into specified Period (DAILY, MONTHLY or YEARLY) will be extracted, if these messages are still persisted, Only XI headers, UDS attributes and audit log entries are exported, message payload is not requested. Resulting exports are saved in specified local file system folder. Export files should be used for import to ELK stack for processing (*link will be published later*).

Tool should be compatible with SAP Netweaver 7.35 or older, tested with SAP Netweaver 7.50.

Overall tool architecture and ELK usage scenario is describled on SAP Blogs (*link will be published later*).

## Configuration file format


Configuration file is represented in JSON file format. Multiple SAP PO instances may be maintained in one configuration file as an array. Each entry will be processed one after another from top to bottom. 

Minimal configuration file is presented below:

	[{
	    "Period": "MONTHLY",
	    "Component": "af.poq.sappoq",
	    "Hostname": "https://poq.local.network:50001",
	    "Password": "EgiEdPVD1tuvMf40Yvne",
	    "Username": "USERNAME"
	}]

| Parameter | Type | Default | Comment |
| --- | --- | --- | --- | 
|	Period                 | string | DAILY | Extraction period, can be set among SAP PO standard values: DAILY, MONTHLY, YEARLY |
|	Component              | string |  --- | AF Component name as displayed in SAP PO |
|	Hostname               | string |  --- | Hostname must be specified complete with schema (http:// or https://) and port number |
|	Username               | string |  --- | Username to use for authentication |
|	Password               | string |  --- | Password to use for authentication |
|	UseCompression         | bool |  false | Enables generating result files with GZIP compression applied (filenames will be adapted to *.txt.gz if enabled) |
|	OutputPath             | string |  /gogoxi | Specifies target directory to output result files |
|	MaxMessagesPerSearch   | int |  10000 | Maximum number of messages to search using AdapterMessageMonitoringVi service. No need to change in production runs. |
|	MaxAuditlinesPerSearch | int |  10000 | Maximum number of audit messages to get using AdapterMessageMonitoringVi service per message ID. No need to change in production runs. |
|	Queuesizes / MessageSearch      | int | 500 | Internal buffer size. Stores overview lines for later search. |
|	Queuesizes / UDSearchAttributes | int | 100000 | Internal buffer size. Stores messages to augment with User-Defined Search attributes. |
|	Queuesizes / Auditlog           | int | 100000 | Internal buffer size. Stores messages to extract audit log entries. |
|	Queuesizes / FileWriters        | int | 1000 | Internal buffer size. Stores lines for file output to disk. |
|	ParallelCount /	MessageSearcher    | int | 2 | Defines number of parallel threads to search messages based on Overview data. |
|	ParallelCount /	UDSearchAttributes | int | 5 | Defines number of parallel threads to augment messages with UDS attributes. |
|	ParallelCount /	AuditlogReader     | int | 5 | Defines number of parallel threads to read audit log entries. |

## Output directory structure

If compression is not enabled:

	<OUTPUT DIRECTORY> /					 		<-- matches "OutputFolder" configuration
		<ADAPTER FRAMEWORK COMPONENT> /      		<-- matches "Component" configuration
			ximessage.<PERIOD>.<CURRENT DATE>.txt   <-- newline delimited JSON file, one line corresponds to XI message header augmented by UDS attributes
			auditlog.<PERIOD>.<CURRENT DATE>.txt    <-- newline delimited JSON file, one line corresponds to audit log message augmented with XI message header and UDS attributes

If compression is enabled, filenames will be \*.txt.gz


## Runtime options

	Usage of gogoxi_exporter.exe:
	  -check
	        Validate config file only [no run]. Shows full configuration file with defaults applied
	  -config string
	        Path to configuration file in JSON format (default "./config.json")
	  -notrim
	        Disable audit log message trimming to 4K bytes	        
	  -overwrite
	        Overwrite export if it exists
	  -verbose
	        Print verbose messages during processing

Tool will prevent accidental output file overwrite. This allows for continuation if one config entry fails (for example, in case of incorrect or expired password).

## File format (ximessage)

File format is newline-delimited JSON. Each line represents one XI message in following structure (example):


	{
	    "ximessage": {
	        "cancelable": false,
	        "connectionName": "SOAP_http://sap.com/xi/XI/System",
	        "direction": "OUTBOUND",
	        "editable": false,
	        "endTime": "2023-11-08T11:09:13.377+03:00",
	        "endpoint": "<local>",
	        "headers": {
	            "original": "content-length=5450\nhttp=POST\ncontent-type=multipart/related; boundary=SAP_11271d92-7e0e-11ee-9105-00000c9d89ea_END; type=\"text/xml\"; start=\"<soap-E1F49308B5D91EEE9FC1C22063DFC044@sap.com>\"\n",
	            "contentLength": 5450
	        },
	        "interface": {
	            "name": "PurchaseRequest_Out",
	            "namespace": "urn:company:Purchase"
	        },
	        "isPersistent": true,
	        "messageID": "e1f49308-b5d9-1eee-9fc1-c22063df6044",
	        "messageKey": "e1f49308-b5d9-1eee-9fc1-c22063df6044\\OUTBOUND\\878415850\\EO\\0\\",
	        "messageType": "Send",
	        "nodeId": "878415850",
	        "persistUntil": "2024-05-06T11:08:57.744+03:00",
	        "protocol": "XI",
	        "qualityOfService": "EO",
	        "receiverName": "EXTERNAL_p",
	        "receiverParty": {
	            "agency": "http://sap.com/xi/XI",
	            "name": "",
	            "schema": "XIParty"
	        },
	        "referenceID": "",
	        "restartable": false,
	        "retries": 10,
	        "retryInterval": 300,
	        "scheduleTime": "2023-11-08T11:08:57.748+03:00",
	        "senderName": "ERP_S4HANA_P",
	        "senderParty": {
	            "agency": "http://sap.com/xi/XI",
	            "name": "",
	            "schema": "XIParty"
	        },
	        "sequenceNumber": 0,
	        "serializationContext": "",
	        "serviceDefinition": "",
	        "softwareComponent": "",
	        "startTime": "2023-11-08T11:08:57.744+03:00",
	        "status": "success",
	        "timesFailed": 0,
	        "transport": "Loopback",
	        "version": "0",
	        "wasEdited": false,
	        "scenarioIdentifier": "dir://ICO/f10d5af248d83dbe9ac62b827e8e223a",
	        "parentID": "",
	        "duration": 15.633,
	        "size": 5450,
	        "messagePriority": 1,
	        "rootID": "",
	        "sequenceID": "",
	        "passportTID": "e78c57187e0d11eeb07100000c9d89ea",
	        "logLocations": [
	            "Receiver JSON Request",
	            "AM"
	        ],
	        "UDS": {
	            "RequestNumber": [
	                "0048588891"
	            ]
	        },
	        "UDSKeys": [
	            "RequestNumber"
	        ]
	    },
	    "gogoxi": {
	        "component": "af.pop.sappop",
	        "extractedOn": "2023-11-08T12:02:25.7824767+03:00"
	    }
	}

Full JSON structure is described in *export_format.go*. Notable remarks:

1. "*ximessage.nodeId*" is String
2. "*ximessage.duration*" is recalculated in seconds, originally in milliseconds
3. "*ximessage.retryInterval*" is recalculated in seconds, originally in milliseconds
4. "*ximessage.persistUntil*" is cleared when originally is equal to 1970-01-01
5. "*ximessage.UDS*" lists user-defined search attributes as object with UDS attribute name as object key and array of values as object value. UDS attribute values are always represented as arrays even if contains only one value.
6. "*ximessage.UDSKeys*" lists all available UDS attribute names as array



## File format (auditlog)

File format is newline-delimited JSON. Each line represents one auditlog entry in following structure (example):


	{
	    "ximessage": {
	        "cancelable": false,
	        "connectionName": "SOAP_http://sap.com/xi/XI/System",
	        "direction": "OUTBOUND",
	        "editable": false,
	        "endTime": "2023-11-08T11:09:13.377+03:00",
	        "endpoint": "<local>",
	        "headers": {
	            "original": "content-length=5450\nhttp=POST\ncontent-type=multipart/related; boundary=SAP_11271d92-7e0e-11ee-9105-00000c9d89ea_END; type=\"text/xml\"; start=\"<soap-E1F49308B5D91EEE9FC1C22063DFC044@sap.com>\"\n",
	            "contentLength": 5450
	        },
	        "interface": {
	            "name": "PurchaseRequest_Out",
	            "namespace": "urn:company:Purchase"
	        },
	        "isPersistent": true,
	        "messageID": "e1f49308-b5d9-1eee-9fc1-c22063df6044",
	        "messageKey": "e1f49308-b5d9-1eee-9fc1-c22063df6044\\OUTBOUND\\878415850\\EO\\0\\",
	        "messageType": "Send",
	        "nodeId": "878415850",
	        "persistUntil": "2024-05-06T11:08:57.744+03:00",
	        "protocol": "XI",
	        "qualityOfService": "EO",
	        "receiverName": "EXTERNAL_p",
	        "receiverParty": {
	            "agency": "http://sap.com/xi/XI",
	            "name": "",
	            "schema": "XIParty"
	        },
	        "referenceID": "",
	        "restartable": false,
	        "retries": 10,
	        "retryInterval": 300,
	        "scheduleTime": "2023-11-08T11:08:57.748+03:00",
	        "senderName": "ERP_S4HANA_P",
	        "senderParty": {
	            "agency": "http://sap.com/xi/XI",
	            "name": "",
	            "schema": "XIParty"
	        },
	        "sequenceNumber": 0,
	        "serializationContext": "",
	        "serviceDefinition": "",
	        "softwareComponent": "",
	        "startTime": "2023-11-08T11:08:57.744+03:00",
	        "status": "success",
	        "timesFailed": 0,
	        "transport": "Loopback",
	        "version": "0",
	        "wasEdited": false,
	        "scenarioIdentifier": "dir://ICO/f10d5af248d83dbe9ac62b827e8e223a",
	        "parentID": "",
	        "duration": 15.633,
	        "size": 5450,
	        "messagePriority": 1,
	        "rootID": "",
	        "sequenceID": "",
	        "passportTID": "e78c57187e0d11eeb07100000c9d89ea",
	        "logLocations": [
	            "Receiver JSON Request",
	            "AM"
	        ],
	        "UDS": {
	            "RequestNumber": [
	                "0048588891"
	            ]
	        },
	        "UDSKeys": [
	            "RequestNumber"
	        ]
	    },
	    "auditlog": {
	        "messageID": "e1f49308-b5d9-1eee-9fc1-c22063df6044",
	        "timestamp": "2023-11-08T11:09:35.869+03:00",
	        "status": "SUCCESS",
	        "localizedText": "Message status set to DLNG",
	        "textKey": "STATUS_SET_SUCCESS",
	        "textKeyParams": {
	            "0": "DLNG"
	        }
	    },
	    "gogoxi": {
	        "component": "af.pop.sappop",
	        "extractedOn": "2023-11-08T12:02:25.7824767+03:00"
	    }
	}

Full JSON structure is described in *export_format.go*. Notable remarks:

1. "*ximessage*" object is identical to XI message export format so all the remarks listed for ximessage file format apply here
2. "*gogoxi*" object is identical to XI message export format
3. "*auditlog.messageID*" refers to XI message ID and is identical to "*ximessage.messageID*"
4. "*auditlog.status*" is decoded from one letter code to full word
- W -> WARNING
- E -> ERROR
- S -> SUCCESS
5. "*auditlog.localizedText*", "*auditlog.textKey*", "*auditlog.textKeyParams*" are cut to 4096 characters unless "*-notrim*" runtime option is applied.
6. "*auditlog.textCutFrom*" attribute shows sum of original length in characters of attributes "*auditlog.localizedText*", "*auditlog.textKey*", "*auditlog.textKeyParams*". If none of these attributes are modified (cut in length) then "*auditlog.textCutFrom*" is missing
7. "*auditlog.textKey*" is cleared if it is equal to "*auditlog.localizedText*"
