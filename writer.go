package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
	"sync/atomic"
)

func ExportWriter(filepath string, input <-chan UnifiedExportFormat, compress bool, counter *int32) {

	writerchan := make(chan []byte, 100)

	go func(in <-chan UnifiedExportFormat, out chan<- []byte) {
		defer close(out)

		for entry := range in {
			jsonText, err := json.Marshal(entry)
			if err != nil {
				fmt.Println("Error during JSON export of xi message: ", entry.XIMessage.MessageKey)
				os.Exit(10)
			}

			out <- jsonText
		}
	}(input, writerchan)

	if compress {
		go outputGZIPFileWriter(filepath, writerchan, counter)
	} else {
		go outputFileWriter(filepath, writerchan, counter)
	}

}

func outputGZIPFileWriter(filepath string, input <-chan []byte, counter *int32) {
	defer wgWriters.Done()

	file, err := os.Create(filepath)
	defer file.Close()

	if err != nil {
		fmt.Printf("Failed creating file [%s]: %s", filepath, err)
		os.Exit(6)
	}

	gzipWriter, err := gzip.NewWriterLevel(file, gzip.BestCompression)
	if err != nil {
		fmt.Printf("Failed to prepare GZIP [%s]: %s", filepath, err)
		os.Exit(6)
	}
	defer gzipWriter.Close()

	for line := range input {

		_, err = gzipWriter.Write(line)

		if err != nil {
			fmt.Printf("Failed writing to file [%s]: %s", filepath, err)
			os.Exit(7)
		}

		_, _ = gzipWriter.Write([]byte("\n"))
		atomic.AddInt32(counter, 1)
	}

}

func outputFileWriter(filepath string, input <-chan []byte, counter *int32) {
	defer wgWriters.Done()

	file, err := os.Create(filepath)
	defer file.Close()

	if err != nil {
		fmt.Printf("Failed creating file [%s]: %s", filepath, err)
		os.Exit(6)
	}

	for line := range input {

		_, err = file.Write(line)
		if err != nil {
			fmt.Printf("Failed writing to file [%s]: %s", filepath, err)
			os.Exit(7)
		}

		file.WriteString("\n")
		atomic.AddInt32(counter, 1)
	}
}
