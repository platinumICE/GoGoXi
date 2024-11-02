package main

import (
	"fmt"
	"os"
	"time"
)

func runner_java(configuration ToolConfiguration, overwritePreviousRun *bool) {
	defer func() {
		// clear stats
		statistics = Statistics{}
	}()

	// prepare HTTP
	init_http(configuration)

	fmt.Println("--------------------")
	fmt.Printf("Started processing component [%s]\n", configuration.Component)

	outputfileAudit, exists1 := generateFilename("auditlog", configuration)
	outputfileXI, exists2 := generateFilename("ximessage", configuration)

	if checkForOverwrite := !*overwritePreviousRun; checkForOverwrite {
		if exists1 || exists2 {
			fmt.Printf("Output files exist, skipping component for the day\n")
			fmt.Println("--------------------")
			return
		}
	}

	//-------------

	startTime := time.Now()

	messageSearch := make(chan XIOverviewDetailsPeriodEntry, configuration.Queuesizes.MessageSearch)
	auditLogsChannel := make(chan UnifiedExportFormat, configuration.Queuesizes.Auditlog)
	UDSAttributesChannel := make(chan SearchResults, configuration.Queuesizes.UDSearchAttributes)

	auditLogWriterChannel := make(chan UnifiedExportFormat, configuration.Queuesizes.FileWriters)
	ximessageWriterChannel := make(chan UnifiedExportFormat, configuration.Queuesizes.FileWriters)

	// Statistics output
	statisticsTicker := time.NewTicker(time.Second)
	defer statisticsTicker.Stop()
	go StatisticsManager(statisticsTicker, messageSearch, auditLogsChannel, UDSAttributesChannel, auditLogWriterChannel, ximessageWriterChannel)

	wgWriters.Add(2)

	go ExportWriter(outputfileAudit, auditLogWriterChannel, configuration.UseCompression, &statistics.AuditlogLines)
	go ExportWriter(outputfileXI, ximessageWriterChannel, configuration.UseCompression, &statistics.XIMessageLines)

	for i := 0; i < configuration.ParallelCount.AuditlogReader; i++ {
		wgAuditers.Add(1)
		go AuditlogReader(configuration, auditLogsChannel, auditLogWriterChannel)
	}

	for i := 0; i < configuration.ParallelCount.UDSearchAttributes; i++ {
		wgEnrichers.Add(1)
		go udsEnricher(configuration, UDSAttributesChannel, ximessageWriterChannel, auditLogsChannel)
	}

	for i := 0; i < configuration.ParallelCount.MessageSearcher; i++ {
		wgSearchers.Add(1)
		go MessageSearcher(configuration, messageSearch, UDSAttributesChannel)
	}

	// let's get this party started
	go OverviewLoader(configuration, messageSearch)

	wgSearchers.Wait()
	close(UDSAttributesChannel) // no more enricher tasks

	// wait for enrichers to finish
	wgEnrichers.Wait()
	close(auditLogsChannel)       // no more audit log tasks
	close(ximessageWriterChannel) // no more files to write

	// wait for auditors to finish
	wgAuditers.Wait()
	close(auditLogWriterChannel) // no more file to write

	// wait for writers to finish
	wgWriters.Wait()

	if statistics.AuditlogLines == 0 && statistics.XIMessageLines == 0 {
		os.Remove(outputfileAudit)
		os.Remove(outputfileXI)
	}

	//---------------------
	elapsed := time.Now().Sub(startTime)

	fmt.Println("--------------------")
	fmt.Printf("Export results for [%s]:\n", configuration.Component)
	fmt.Printf("  XI Messages     : %d / %d\n", statistics.XIMessageLines, statistics.MaxXIMessageLines)
	fmt.Printf("  Auditlog lines  : %d\n", statistics.AuditlogLines)
	fmt.Printf("  Time taken      : %v\n", elapsed)
	fmt.Println("--------------------")

}
