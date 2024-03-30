package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

import (
	"encoding/json"
)

type MultipleToolConfiguration []ToolConfiguration
type ToolConfiguration struct {
	Period                 string
	Component              string
	Hostname               string
	Username               string
	Password               string
	UseCompression         bool
	OutputPath             string
	MaxMessagesPerSearch   int
	MaxAuditlinesPerSearch int
	SkipUDS                bool
	Queuesizes             struct {
		MessageSearch      int
		Auditlog           int
		UDSearchAttributes int
		FileWriters        int
	}
	ParallelCount struct {
		MessageSearcher    int
		UDSearchAttributes int
		AuditlogReader     int
	}
}

type RuntimeConfiguration struct {
	BeVerbose          bool
	NoTrimForAuditLogs bool
}

///////////// global declarations

var wgSearchers sync.WaitGroup
var wgAuditers sync.WaitGroup
var wgEnrichers sync.WaitGroup
var wgWriters sync.WaitGroup

var runtimeConfig RuntimeConfiguration

func main() {

	configFile := flag.String("config", "./config.json", "Path to configuration file in JSON format")
	validateConfig := flag.Bool("check", false, "Validate config file only [no run]")
	overwritePreviousRun := flag.Bool("overwrite", false, "Overwrite export if it exists")
	verbose := flag.Bool("verbose", false, "Print verbose messages during processing")
	noTrimAuditLogs := flag.Bool("notrim", false, "Disable audit log message trimming to 4K bytes")

	flag.Parse()

	runtimeConfig.BeVerbose = *verbose
	runtimeConfig.NoTrimForAuditLogs = *noTrimAuditLogs

	configuration := readConfig(*configFile)
	if *validateConfig {
		printConfig(configuration)
		os.Exit(0)
	}

	for _, conf := range configuration {
		if strings.HasPrefix(conf.Component, "af.") {
			runner_java(conf, overwritePreviousRun)
		} else {
			fmt.Println("--------------------")
			fmt.Printf("Only component type AF is supported now \n")
			fmt.Printf("Skipping %s \n", conf.Component)
			fmt.Println("--------------------")
		}
	}

	fmt.Println("--------------------")
	fmt.Println("--------------------")
	fmt.Println("Export fully done")
}

func readConfig(configFile string) MultipleToolConfiguration {

	if configFile == "" {
		panic("Config file not specified")
	}

	fmt.Println("Reading config file:", configFile)

	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic("Config file not found")
	}

	var config MultipleToolConfiguration

	err = json.Unmarshal(file, &config)
	if err != nil {
		fmt.Println("Error parsing config file:", err)
		panic("Cannot unmarshal JSON")
	}

	for i := range config {

		conf := &config[i]

		// check config
		if conf.Hostname == "" {
			panic("Hostname is empty")
		}

		// check config
		if conf.Username == "" || conf.Password == "" {
			panic("Username / password are empty")
		}

		// check config
		if conf.Component == "" {
			panic("Component is required for correct work")
		}

		// baked-in defaults
		if conf.Period == "" {
			conf.Period = "DAILY"
		}

		if conf.MaxMessagesPerSearch == 0 {
			conf.MaxMessagesPerSearch = 10000
		}

		if conf.MaxAuditlinesPerSearch == 0 {
			conf.MaxAuditlinesPerSearch = 10000
		}

		if conf.OutputPath == "" {
			conf.OutputPath = "/gogoxi"
		}

		if conf.Queuesizes.MessageSearch == 0 {
			conf.Queuesizes.MessageSearch = 500
		}

		if conf.Queuesizes.Auditlog == 0 {
			conf.Queuesizes.Auditlog = 100000
		}

		if conf.Queuesizes.UDSearchAttributes == 0 {
			conf.Queuesizes.UDSearchAttributes = 100000
		}

		if conf.Queuesizes.FileWriters == 0 {
			conf.Queuesizes.FileWriters = 1000
		}

		if conf.ParallelCount.MessageSearcher == 0 {
			conf.ParallelCount.MessageSearcher = 2
		}

		if conf.ParallelCount.AuditlogReader == 0 {
			conf.ParallelCount.AuditlogReader = 5
		}

		if conf.ParallelCount.UDSearchAttributes == 0 {
			conf.ParallelCount.UDSearchAttributes = 5
		}

	}

	return config
}

func printConfig(conf MultipleToolConfiguration) {

	jsonFile, _ := json.MarshalIndent(conf, "", "\t")

	print(string(jsonFile))

}

func generateFilename(typeName string, conf ToolConfiguration) (string, bool) {

	path := fmt.Sprintf("%s/%s", conf.OutputPath, conf.Component)
	err := os.MkdirAll(path, 0750)
	if err != nil {
		panic(fmt.Sprintf("Cannot create output directory: %s", path))
	}

	filename := fmt.Sprintf("%s/%s.%s.%s.txt", path, typeName, conf.Period, time.Now().Format("2006-01-02"))

	if conf.UseCompression {
		filename = filename + ".gz"
	}

	filestat, err := os.Stat(filename)
	if filestat == nil {
		return filename, false
	} else if filestat.Size() == 0 {
		return filename, false
	} else {
		return filename, true
	}
}
