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

type UploadResult struct {
	Success   bool      `json:"success"`
	Hash      string    `json:"hash"`
	URL       string    `json:"url"`
	Filename  string    `json:"filename"`
	Size      int64     `json:"size"`
	Timestamp time.Time `json:"timestamp"`
	Error     string    `json:"error,omitempty"`
}

type PinataConfig struct {
	APIKey    string
	APISecret string
	Gateway   string
}

type PinataResponse struct {
	IpfsHash  string `json:"IpfsHash"`
	PinSize   int    `json:"PinSize"`
	Timestamp string `json:"Timestamp"`
}

type PinataUploader struct {
	config     PinataConfig
	httpClient *http.Client
}

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

func (p *PinataUploader) UploadPhoto(photoPath string, filename string) (*UploadResult, error) {
	fileData, err := os.ReadFile(photoPath)
	if err != nil {
		return &UploadResult{
			Success: false,
			Error:   fmt.Sprintf("Failed to read file: %v", err),
		}, err
	}

	if filename == "" {
		filename = filepath.Base(photoPath)
	}

	return p.uploadFileData(fileData, filename)
}

func (p *PinataUploader) UploadPhotoFromBytes(data []byte, filename string) (*UploadResult, error) {
	return p.uploadFileData(data, filename)
}

func (p *PinataUploader) uploadFileData(data []byte, filename string) (*UploadResult, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

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

	metadata := map[string]string{
		"name": filename,
	}
	metadataJSON, _ := json.Marshal(metadata)
	writer.WriteField("pinataMetadata", string(metadataJSON))

	options := map[string]interface{}{
		"cidVersion": 0,
	}
	optionsJSON, _ := json.Marshal(options)
	writer.WriteField("pinataOptions", string(optionsJSON))

	writer.Close()

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

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return &UploadResult{
			Success: false,
			Error:   fmt.Sprintf("Failed to send request: %v", err),
		}, err
	}
	defer resp.Body.Close()

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

	return nil
}
