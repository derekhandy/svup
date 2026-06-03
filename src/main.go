package main

import (
	"fmt"
	"log"
	"os"
)

type Flags struct {
	verbose bool
}

func printResult(result *UploadResult) {
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
}

func main() {
	flags := Flags{
		verbose: false,
	}

	if len(os.Args) > 2 {
		if os.Args[2] != "-v" {
			fmt.Println("Invalid arguments.")
			fmt.Println("Usage: svup path/to/photo.jpg")
			fmt.Println("	-v : Prints request metadata.")
			return
		}
		flags.verbose = true
	}

	if len(os.Args) <= 1 {
		fmt.Println("No photo path provided.")
		fmt.Println("Usage: svup path/to/photo.jpg [-v]")
		return
	}

	apiKey := os.Getenv("PINATA_API_KEY")
	apiSecret := os.Getenv("PINATA_API_SECRET")

	if apiKey == "" || apiSecret == "" {
		fmt.Println("Pinata API credentials not found!")
		fmt.Println("Please set your Pinata API credentials")
		fmt.Println("Or get them from: https://app.pinata.cloud/keys")
		fmt.Println()
		fmt.Println("Usage: svup path/to/photo.jpg")
		return
	}

	uploader := NewPinataUploader(apiKey, apiSecret)

	if err := uploader.TestConnection(); err != nil {
		log.Fatalf("Failed to connect to Pinata: %v", err)
	}

	photoPath := os.Args[1]

	result, err := uploader.UploadPhoto(photoPath, "")
	if err != nil {
		log.Printf("Error uploading photo: %v", err)
		return
	}

	if flags.verbose {
		printResult(result)
	} else {
		fmt.Println(result.URL)
	}
}
