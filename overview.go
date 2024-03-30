package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

import (
	"encoding/xml"
)

type XIComponents struct {
	XMLName    xml.Name `xml:"MessageStatisticsQueryResults"`
	Components []string `xml:"XIComponents>Component"`
}

type XIOverviewPeriods struct {
	XMLName xml.Name `xml:"MessageStatisticsQueryResults"`
	Period  []struct {
		Type     string
		Interval []struct {
			Begin string
			End   string
		}
	} `xml:"Periods>Period"`
}

type XIOverviewDetailsXML struct {
	XMLName     xml.Name `xml:"MessageStatisticsQueryResults"`
	PeriodEntry []struct {
		Entry []string `xml:"Entry"`
	} `xml:"Data>DataRows>Row"`
}

type XIOverviewDetailsPeriodEntry struct {
	Component         string
	Period            string
	Start             string
	End               string
	SenderComponent   string
	ReceiverComponent string
	Interface         string
	Total             int32
}

func OverviewLoader(conf ToolConfiguration, output chan<- XIOverviewDetailsPeriodEntry) {
	pointOfNoReturn := false
	defer func() {
		if r := recover(); r != nil {
			if pointOfNoReturn {
				panic(r)
			}

			fmt.Println("Error caught on Overview call")
			fmt.Println(r)
		}
	}()
	defer close(output)

	/////////////////////
	/// step 1: http request

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/mdt/messageoverviewqueryservlet?component=%s&view=SR_ENTRY_OVERVIEW_XPI", conf.Hostname, conf.Component), nil)
	req.SetBasicAuth(conf.Username, conf.Password)
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

	pointOfNoReturn = true

	overviewPeriods := new(XIOverviewPeriods)
	err = xml.Unmarshal(responseBytes, overviewPeriods)

	if err != nil {
		fmt.Printf("Please verify that host [%s], username [%s] and password are correct\n", conf.Hostname, conf.Username)
		fmt.Printf("HTTP call returned: %s\n", err.Error())
		os.Exit(2)
	}

	if len(overviewPeriods.Period) == 0 {
		fmt.Printf("No period info was found on the host %s\n", conf.Hostname)
		fmt.Printf("Please verify that component [%s] is set correctly\n", conf.Component)
		os.Exit(1)
	}

	/////////////////////
	/// step 2: http request

	for _, period := range overviewPeriods.Period {
		if period.Type != conf.Period {
			continue
		}

		for _, interval := range period.Interval {

			req, err = http.NewRequest("GET", fmt.Sprintf("%s/mdt/messageoverviewqueryservlet?component=%s&view=SR_ENTRY_OVERVIEW_XPI&begin=%s&end=%s", conf.Hostname, conf.Component, interval.Begin, interval.End), nil)
			req.SetBasicAuth(conf.Username, conf.Password)
			resp, err = client.Do(req)
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

			responseBytes, err = io.ReadAll(resp.Body)
			resp.Body.Close()
			overviewDetailsXML := new(XIOverviewDetailsXML)
			err = xml.Unmarshal(responseBytes, overviewDetailsXML)

			for _, entry := range overviewDetailsXML.PeriodEntry {

				output <- XIOverviewDetailsPeriodEntry{
					Component:         conf.Component,
					Period:            period.Type,
					Start:             interval.Begin,
					End:               interval.End,
					SenderComponent:   overviewParseStr(entry.Entry[0]),
					ReceiverComponent: overviewParseStr(entry.Entry[1]),
					Interface:         overviewParseStr(entry.Entry[2]),
					Total:             overviewParseInt(entry.Entry[3], entry.Entry[4], entry.Entry[5], entry.Entry[6]),
				}
			}
		}
	}
}

func overviewParseStr(input string) string {
	if input == "-" {
		return ""
	}
	return input
}

func overviewParseInt(input1 string, input2 string, input3 string, input4 string) int32 {
	var result int32
	if input1 != "-" {
		i, _ := strconv.ParseInt(input1, 10, 32)
		result += int32(i)
	}
	if input2 != "-" {
		i, _ := strconv.ParseInt(input2, 10, 32)
		result += int32(i)
	}
	if input3 != "-" {
		i, _ := strconv.ParseInt(input3, 10, 32)
		result += int32(i)
	}
	if input4 != "-" {
		i, _ := strconv.ParseInt(input4, 10, 32)
		result += int32(i)
	}

	return result
}
