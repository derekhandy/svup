package main

import (
	"fmt"
	"log"
	"os"
)

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
	apiKey := os.Getenv("PINATA_API_KEY")
	apiSecret := os.Getenv("PINATA_API_SECRET")

	if apiKey == "" || apiSecret == "" {
		fmt.Println("🔑 Pinata API credentials not found!")
		fmt.Println("Please set your Pinata API credentials:")
		fmt.Println("  set PINATA_API_KEY=your_api_key")
		fmt.Println("  set PINATA_API_SECRET=your_api_secret")
		fmt.Println()
		fmt.Println("Or get them from: https://app.pinata.cloud/keys")
		fmt.Println()
		fmt.Println("Usage: pinata-uploader.exe path/to/photo.jpg")
		return
	}

	uploader := NewPinataUploader(apiKey, apiSecret)

	fmt.Println("Testing connection to Pinata...")
	if err := uploader.TestConnection(); err != nil {
		log.Fatalf("Failed to connect to Pinata: %v", err)
	}

	if len(os.Args) <= 1 {
		fmt.Println("No photo path provided.")
		fmt.Println("Usage: pinata-uploader.exe path/to/photo.jpg")
		return
	}

	photoPath := os.Args[1]
	fmt.Printf("\n📸 Uploading photo: %s\n", photoPath)

	result, err := uploader.UploadPhoto(photoPath, "")
	if err != nil {
		log.Printf("Error uploading photo: %v", err)
		return
	}

	printResult(result)
}
