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