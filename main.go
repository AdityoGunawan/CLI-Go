package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

// LogFile struct
type LogFile struct {
	DateTime string `json:"datetime"`
	Level    string `json:"level"`
	Message  string `json:"message"`
}

// LogFileList struct
type LogFileList struct {
	LogFiles []LogFile `json:"logfiles"`
}

func main() {
	// Get flag
	var (
		outputFile = flag.String("o", "", "Output file")
		textType   = flag.String("t", "text", "Text type")
		help       = flag.Bool("h", false, "Help")
	)

	flag.Parse()

	if *help {
		fmt.Println("Usage: go run main.go /var/log/file.log -t json, -t text")
		os.Exit(0)
	}

	// Get file
	file, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Get output file
	var output io.Writer
	if *outputFile != "" {
		f, err := os.Create(*outputFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		output = f
	} else {
		output = os.Stdout
	}

	// Get text type
	switch *textType {
	case "text":
		// Read file
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			// Get line
			line := scanner.Text()

			// Split line
			splitLine := strings.SplitN(line, " ", 4)

			// Get datetime
			datetime, err := time.Parse("2006-01-02T15:04:05.000Z", splitLine[0])
			if err != nil {
				log.Fatal(err)
			}

			// Get level
			level := splitLine[1]

			// Get message
			message := splitLine[3]

			// Write to output
			fmt.Fprintf(output, "%s %s %s %s", datetime.Format("2006-01-02"), datetime.Format("15:04:05"), level, message)
		}
	case "json":
		// Create LogFileList
		logFileList := LogFileList{}

		// Read file
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			// Get line
			line := scanner.Text()

			// Split line
			splitLine := strings.SplitN(line, " ", 4)

			// Get datetime
			datetime, err := time.Parse("2006-01-02T15:04:05.000Z", splitLine[0])
			if err != nil {
				log.Fatal(err)
			}

			// Get level
			level := splitLine[1]

			// Get message
			message := splitLine[3]

			// Create LogFile
			logFile := LogFile{
				DateTime: datetime.Format("2006-01-02 15:04:05"),
				Level:    level,
				Message:  message,
			}

			// Append LogFile to LogFileList
			logFileList.LogFiles = append(logFileList.LogFiles, logFile)
		}

		// Marshal LogFileList
		logFileListJSON, err := json.Marshal(logFileList)
		if err != nil {
			log.Fatal(err)
		}

		// Write to output
		fmt.Fprintf(output, "%s", logFileListJSON)
	default:
		log.Fatal("Invalid text type")
	}

	// Check error
	if err != nil {
		log.Fatal(err)
	}
}

/*
Command
$ go run main.go /var/log/file.log -t json
$ go run main.go -h
*/
