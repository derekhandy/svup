//
//												svup @ v1.1.2
//
//									MIT License, Copyright (c) 2026 Derek Handy
//							Project can be found at: https://github.com/derekhandy/svup
//

package main

import (
	"fmt"
	"log"
	"os"

	svup "github.com/derekhandy/svuplib"
)

type Flags struct {
	verbose bool
}

func parseArgs(args []string) (string, Flags, error) {
	flags := Flags{}
	var filePath string

	for _, arg := range args {
		switch arg {
		case "-v", "--verbose":
			flags.verbose = true
		default:
			if filePath != "" {
				return "", flags, fmt.Errorf("invalid arguments")
			}
			filePath = arg
		}
	}

	if filePath == "" {
		return "", flags, fmt.Errorf("no file path provided")
	}

	return filePath, flags, nil
}

func printUsage() {
	fmt.Println("Usage: svup [-v] path/to/file")
	fmt.Println("  -v, --verbose : Prints request metadata.")
}

func main() {
	filePath, flags, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		printUsage()
		return
	}

	apiKey := os.Getenv("PINATA_API_KEY")
	apiSecret := os.Getenv("PINATA_API_SECRET")

	if apiKey == "" || apiSecret == "" {
		fmt.Println("Pinata API credentials not found")
		return
	}

	uploader := svup.NewPinataUploader(apiKey, apiSecret)

	if err := uploader.TestConnection(); err != nil {
		log.Fatalf("Failed to connect to Pinata: %v", err)
	}

	result, err := uploader.UploadFile(filePath, "")
	if err != nil {
		log.Printf("Error uploading file: %v", err)
		return
	}

	if flags.verbose {
		fmt.Printf("Upload Result:\n")
		fmt.Printf("  Success: %t\n", result.Success)
		if result.Success {
			fmt.Printf("  Hash: %s\n", result.Hash)
			fmt.Printf("  URL: %s\n", result.URL)
			fmt.Printf("  Filename: %s\n", result.Filename)
			fmt.Printf("  Size: %d bytes\n", result.Size)
			fmt.Printf("  Timestamp: %s\n", result.Timestamp.Format("2006-01-02 15:04:05"))
		} else {
			fmt.Printf("  Error: %s\n", result.Error)
		}
	} else {
		fmt.Println(result.URL)
	}
}
