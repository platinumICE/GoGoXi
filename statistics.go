package main

import (
	"fmt"
	"time"
)

type Statistics struct {
	AuditlogLines     int32
	XIMessageLines    int32
	MaxXIMessageLines int32
}

var statistics Statistics

func StatisticsManager(ticker *time.Ticker, messageSearch chan XIOverviewDetailsPeriodEntry, auditlog chan UnifiedExportFormat, uds chan SearchResults, logwriter chan UnifiedExportFormat, ximessage chan UnifiedExportFormat) {
	// fmt.Printf("\033[s") // ANSI escape sequence for saving cursor position
	firstTime := true

	for range ticker.C {
		if !firstTime {
			fmt.Printf("\033[%dA", 10) // Move up 'n' lines
		}
		firstTime = false

		fmt.Printf("\033[K") // Clear the line
		fmt.Println("--------------------")
		fmt.Printf("\033[K") // Clear the line
		fmt.Println("Buffers:")
		fmt.Printf("\033[K") // Clear the line
		fmt.Printf("  Search    : %d / %d\n", len(messageSearch), cap(messageSearch))
		fmt.Printf("\033[K") // Clear the line
		fmt.Printf("  UDS       : %d / %d\n", len(uds), cap(uds))
		fmt.Printf("\033[K") // Clear the line
		fmt.Printf("  Audit log : %d / %d\n", len(auditlog), cap(auditlog))
		// fmt.Printf("  XI file   : %d / %d\n", len(logwriter), cap(logwriter))
		// fmt.Printf("  Audit file: %d / %d\n", len(ximessage), cap(ximessage))
		fmt.Printf("\033[K") // Clear the line
		fmt.Println("")
		fmt.Printf("\033[K") // Clear the line
		fmt.Println("Exported so far:")
		fmt.Printf("\033[K") // Clear the line
		fmt.Printf("  XI Messages     : %d+\n", statistics.XIMessageLines)
		fmt.Printf("\033[K") // Clear the line
		fmt.Printf("  Auditlog lines  : %d+\n", statistics.AuditlogLines)
		fmt.Printf("\033[K") // Clear the line
		fmt.Println("--------------------")
	}
}
