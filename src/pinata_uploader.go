package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// UploadResult represents the result of an IPFS upload operation
type UploadResult struct {
	Success   bool      `json:"success"`
	Hash      string    `json:"hash"`
	URL       string    `json:"url"`
	Filename  string    `json:"filename"`
	Size      int64     `json:"size"`
	Timestamp time.Time `json:"timestamp"`
	Error     string    `json:"error,omitempty"`
}

// PinataConfig holds Pinata API configuration
type PinataConfig struct {
	APIKey    string
	APISecret string
	Gateway   string
}

// PinataResponse represents Pinata API response
type PinataResponse struct {
	IpfsHash  string `json:"IpfsHash"`
	PinSize   int    `json:"PinSize"`
	Timestamp string `json:"Timestamp"`
}

// PinataUploader handles photo uploads using Pinata API
type PinataUploader struct {
	config     PinataConfig
	httpClient *http.Client
}

// NewPinataUploader creates a new Pinata uploader
func NewPinataUploader(apiKey, apiSecret string) *PinataUploader {
	return &PinataUploader{
		config: PinataConfig{
			APIKey:    apiKey,
			APISecret: apiSecret,
			Gateway:   "https://gateway.pinata.cloud/ipfs/",
		},
		httpClient: &http.Client{Timeout: 60 * time.Second},
	}
}

// UploadPhoto uploads a photo to Pinata IPFS
func (p *PinataUploader) UploadPhoto(photoPath string, filename string) (*UploadResult, error) {
	// Read the photo file
	fileData, err := os.ReadFile(photoPath)
	if err != nil {
		return &UploadResult{
			Success: false,
			Error:   fmt.Sprintf("Failed to read file: %v", err),
		}, err
	}

	// Use provided filename or extract from path
	if filename == "" {
		filename = filepath.Base(photoPath)
	}

	return p.uploadFileData(fileData, filename)
}

// UploadPhotoFromBytes uploads photo data from a byte slice
func (p *PinataUploader) UploadPhotoFromBytes(data []byte, filename string) (*UploadResult, error) {
	return p.uploadFileData(data, filename)
}

// uploadFileData handles the actual upload to Pinata
func (p *PinataUploader) uploadFileData(data []byte, filename string) (*UploadResult, error) {
	// Create multipart form data
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add file field
	fileWriter, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return &UploadResult{
			Success: false,
			Error:   fmt.Sprintf("Failed to create form file: %v", err),
		}, err
	}

	_, err = fileWriter.Write(data)
	if err != nil {
		return &UploadResult{
			Success: false,
			Error:   fmt.Sprintf("Failed to write file data: %v", err),
		}, err
	}

	// Add metadata
	metadata := map[string]string{
		"name": filename,
	}
	metadataJSON, _ := json.Marshal(metadata)
	writer.WriteField("pinataMetadata", string(metadataJSON))

	// Add options
	options := map[string]interface{}{
		"cidVersion": 0,
	}
	optionsJSON, _ := json.Marshal(options)
	writer.WriteField("pinataOptions", string(optionsJSON))

	writer.Close()

	// Create HTTP request
	req, err := http.NewRequest("POST", "https://api.pinata.cloud/pinning/pinFileToIPFS", &buf)
	if err != nil {
		return &UploadResult{
			Success: false,
			Error:   fmt.Sprintf("Failed to create request: %v", err),
		}, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("pinata_api_key", p.config.APIKey)
	req.Header.Set("pinata_secret_api_key", p.config.APISecret)

	// Send request
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return &UploadResult{
			Success: false,
			Error:   fmt.Sprintf("Failed to send request: %v", err),
		}, err
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &UploadResult{
			Success: false,
			Error:   fmt.Sprintf("Failed to read response: %v", err),
		}, err
	}

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("pinata API error (%d): %s", resp.StatusCode, string(body))
		return &UploadResult{
			Success: false,
			Error:   errMsg,
		}, fmt.Errorf(errMsg)
	}

	// Parse response
	var result PinataResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return &UploadResult{
			Success: false,
			Error:   fmt.Sprintf("Failed to parse response: %v", err),
		}, err
	}

	// Construct the URL
	url := p.config.Gateway + result.IpfsHash

	return &UploadResult{
		Success:   true,
		Hash:      result.IpfsHash,
		URL:       url,
		Filename:  filename,
		Size:      int64(len(data)),
		Timestamp: time.Now(),
	}, nil
}

// TestConnection tests the connection to Pinata API
func (p *PinataUploader) TestConnection() error {
	req, err := http.NewRequest("GET", "https://api.pinata.cloud/data/testAuthentication", nil)
	if err != nil {
		return fmt.Errorf("failed to create test request: %v", err)
	}

	req.Header.Set("pinata_api_key", p.config.APIKey)
	req.Header.Set("pinata_secret_api_key", p.config.APISecret)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to Pinata API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Pinata API authentication failed (%d): %s", resp.StatusCode, string(body))
	}

	fmt.Println("✅ Successfully connected to Pinata IPFS API")
	return nil
}

// GetPhotoInfo previously provided metadata lookup against Pinata.
// If needed in the future consider re-adding a lean version focused on the required fields.
