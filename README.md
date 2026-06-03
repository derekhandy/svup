# svup

A simple command-line tool for uploading photos and files to IPFS (InterPlanetary File System) using the Pinata API. This tool provides an easy way to pin files to IPFS and get permanent, decentralized URLs for your content.

## Overview

IPFS Photo Uploader is a lightweight Go application that interfaces with Pinata's IPFS pinning service. It allows you to upload files to IPFS with a single command, receiving an IPFS hash and gateway URL that can be used to access your content from anywhere on the IPFS network.

## Features

- **Simple CLI Interface**: Upload files with a single command
- **Pinata Integration**: Uses Pinata's reliable IPFS pinning service
- **Connection Testing**: Verify API credentials before uploading
- **Detailed Results**: Returns IPFS hash, gateway URL, file size, and timestamp
- **Error Handling**: Clear error messages for troubleshooting
- **Fast Uploads**: Efficient multipart form data handling

## Prerequisites

- **Go 1.21+** (for building from source)
- **Pinata API Account** - Get your API keys from [Pinata](https://app.pinata.cloud/keys)
- **Internet Connection** - Required for API communication

## Installation

### Option 1: Build from Source

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd ipfs
   ```

2. Build the executable:
   ```bash
   cd src
   go build -o svup.exe
   ```

### Option 2: Use Pre-built Binary

Download the pre-built `svup.exe` from the releases page (if available).

## Setup

### Get Pinata API Credentials

1. Sign up for a free account at [Pinata](https://app.pinata.cloud/)
2. Navigate to [API Keys](https://app.pinata.cloud/keys)
3. Create a new API key pair
4. Copy your `API Key` and `Secret API Key`

### Configure Environment Variables

**Windows (Command Prompt):**
```cmd
set PINATA_API_KEY=your_api_key_here
set PINATA_API_SECRET=your_secret_api_key_here
```

**Windows (PowerShell):**
```powershell
$env:PINATA_API_KEY="your_api_key_here"
$env:PINATA_API_SECRET="your_secret_api_key_here"
```

**Linux/Mac:**
```bash
export PINATA_API_KEY=your_api_key_here
export PINATA_API_SECRET=your_secret_api_key_here
```

**Permanent Setup (Windows):**
1. Open System Properties → Environment Variables
2. Add `PINATA_API_KEY` and `PINATA_API_SECRET` to User variables
3. Restart your terminal

## Usage

### Basic Upload

Upload a single file to IPFS:

```bash
svup.exe path/to/photo.jpg
```
### Output

The tool will display:
- Connection test result
- Upload progress
- Upload result with:
  - Success status
  - IPFS Hash (CID)
  - Gateway URL
  - Filename
  - File size in bytes
  - Upload timestamp

**Example Output:**
```
Testing connection to Pinata...
Successfully connected to Pinata IPFS API

Uploading photo: path/to/photo.png

Upload Result:
  Success: true
  Hash: QmXxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
  URL: https://gateway.pinata.cloud/ipfs/QmXxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
  Filename: vacation.jpg
  Size: 2456789 bytes
  Timestamp: 2024-01-15 14:30:25
```

## How It Works

1. **Connection Test**: Verifies API credentials by calling Pinata's authentication endpoint
2. **File Reading**: Reads the specified file from disk
3. **Multipart Upload**: Creates a multipart form with:
   - File data
   - Metadata (filename)
   - Pinata options (CID version)
4. **API Request**: Sends POST request to Pinata's pinning endpoint
5. **Response Processing**: Extracts IPFS hash and constructs gateway URL
6. **Result Display**: Shows upload result with all relevant information

## Project Structure

```
ipfs/
├── src/
│   ├── unified_main.go        # Main entry point and CLI interface
│   ├── pinata_uploader.go     # Pinata API integration and upload logic
│   ├── svup.exe               # Precompiled executable
│   └── go.mod                 # Go module definition
└── README.md                  # This file
```

## API Reference

### PinataUploader

The main uploader struct that handles all IPFS operations.

#### Methods

**`NewPinataUploader(apiKey, apiSecret string) *PinataUploader`**
- Creates a new Pinata uploader instance
- Parameters:
  - `apiKey`: Your Pinata API key
  - `apiSecret`: Your Pinata secret API key
- Returns: Configured uploader instance

**`TestConnection() error`**
- Tests the connection to Pinata API
- Verifies API credentials are valid
- Returns error if authentication fails

**`UploadPhoto(photoPath string, filename string) (*UploadResult, error)`**
- Uploads a file from disk to IPFS
- Parameters:
  - `photoPath`: Full path to the file to upload
  - `filename`: Optional custom filename (empty string uses file's basename)
- Returns: Upload result with hash and URL, or error

**`UploadPhotoFromBytes(data []byte, filename string) (*UploadResult, error)`**
- Uploads file data from memory (byte slice) to IPFS
- Parameters:
  - `data`: File content as byte slice
  - `filename`: Filename for the upload
- Returns: Upload result with hash and URL, or error

### UploadResult

Structure containing upload operation results.

```go
type UploadResult struct {
    Success   bool      // Whether upload succeeded
    Hash      string    // IPFS hash (CID)
    URL       string    // Gateway URL for accessing the file
    Filename  string   // Name of the uploaded file
    Size      int64    // File size in bytes
    Timestamp time.Time // Upload timestamp
    Error     string   // Error message (if failed)
}
```

## Pinata Configuration

### Default Gateway

The tool uses Pinata's public gateway by default:
```
https://gateway.pinata.cloud/ipfs/
```

### CID Version

Currently uses CIDv0 for compatibility. This can be modified in `pinata_uploader.go`:
```go
options := map[string]interface{}{
    "cidVersion": 0,  // Change to 1 for CIDv1
}
```

## Integration Examples

### Go Program

```go
package main

import (
    "os"
    "fmt"
)

func main() {
    uploader := NewPinataUploader(
        os.Getenv("PINATA_API_KEY"),
        os.Getenv("PINATA_API_SECRET"),
    )
    
    result, err := uploader.UploadPhoto("image.jpg", "")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("Uploaded to: %s\n", result.URL)
}
```

### Batch Upload Script

```bash
#!/bin/bash
for file in images/*.jpg; do
    echo "Uploading $file..."
    svup.exe "$file"
    echo ""
done
```

## Troubleshooting

### API Credentials Not Found

**Error:** `Pinata API credentials not found!`

**Solution:**
- Ensure `PINATA_API_KEY` and `PINATA_API_SECRET` are set
- Verify environment variables are set in the current terminal session
- Check for typos in variable names

### Connection Test Failed

**Error:** `Failed to connect to Pinata: authentication failed`

**Solution:**
- Verify API keys are correct
- Check that keys haven't been revoked in Pinata dashboard
- Ensure you're using the correct key pair (API Key + Secret API Key)

### File Not Found

**Error:** `Failed to read file: open path/to/file: no such file or directory`

**Solution:**
- Verify the file path is correct
- Use absolute paths if relative paths don't work
- Check file permissions

### Upload Failed

**Error:** `pinata API error (400/500): ...`

**Solution:**
- Check file size limits (Pinata free tier has limits)
- Verify file format is supported
- Check Pinata service status
- Review error message for specific details

### Network Timeout

**Error:** `Failed to send request: context deadline exceeded`

**Solution:**
- Check internet connection
- Try again (may be temporary network issue)
- Increase timeout in code if needed (default: 60 seconds)

## Pinata Limits

### Free Tier
- 1 GB storage
- 1 GB bandwidth per month
- Unlimited files

### Paid Tiers
- Higher storage limits
- More bandwidth
- Priority support

Check [Pinata Pricing](https://www.pinata.cloud/pricing) for current limits.

## Security Notes

- **Never commit API keys** to version control
- Use environment variables for credentials
- Rotate API keys periodically
- Use separate API keys for development and production
- Review Pinata's [security best practices](https://docs.pinata.cloud/)

## IPFS Gateway Access

Once uploaded, your file is accessible via:

1. **Pinata Gateway** (default):
   ```
   https://gateway.pinata.cloud/ipfs/<hash>
   ```

2. **Public IPFS Gateways**:
   ```
   https://ipfs.io/ipfs/<hash>
   https://cloudflare-ipfs.com/ipfs/<hash>
   ```

3. **Local IPFS Node** (if running):
   ```
   http://localhost:8080/ipfs/<hash>
   ```

## Future Enhancements

Potential improvements:
- Batch upload support
- Directory upload
- Progress bars for large files
- Custom metadata support
- Multiple gateway selection
- CID version selection
- Upload history/logging
- Configuration file support

## Contributing

When contributing:
1. Follow Go best practices
2. Add error handling for edge cases
3. Update this README for new features
4. Test with various file types and sizes
5. Ensure backward compatibility

## License

[Add your license information here]

## Acknowledgments

- [Pinata](https://www.pinata.cloud/) - IPFS pinning service
- [IPFS](https://ipfs.io/) - InterPlanetary File System protocol
- Go standard library for HTTP and file operations

## Resources

- [Pinata Documentation](https://docs.pinata.cloud/)
- [IPFS Documentation](https://docs.ipfs.io/)
- [Pinata API Reference](https://docs.pinata.cloud/api-reference)
- [Get Pinata API Keys](https://app.pinata.cloud/keys)

