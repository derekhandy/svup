# svup


<b> svup </b> is a cross-platform CLI tool for uploading a file to IPFS through Pinata and returning the CID gateway URL. Available for Linux, Windows and MacOS architectures.

Requires <b> Go 1.26 </b> or newer and a <b> Pinata </b> account.


```bash
### Usage Examples

# Returns only the gateway URL

$ svup path/to/file.jpg
https://gateway.pinata.cloud/ipfs/Qm...


# Returns upload metadata

$ svup -v path/to/file.jpg
Upload Result:
  Success: true
  Hash: Qm...
  URL: https://gateway.pinata.cloud/ipfs/Qm...
  Filename: file.jpg
  Size: 12345 bytes
  Timestamp: 2026-06-03 12:00:00
```

## Install

Download from repository then build and move the executable to environment path:

```bash
go get github.com/derekhandy/svup
go mod tidy
```

```bash
# Linux
go build -o svup
sudo mv svup /usr/local/bin/svup

# Mac OS
go build -o svup
sudo mv svup /usr/local/bin/svup

# Windows PowerShell
go build -o svup.exe
```

## Set Environment Variables

Create Pinata API credentials at:

```text
https://app.pinata.cloud/
```

```bash
# Linux / Mac OS
export PINATA_API_KEY="your_api_key"
export PINATA_API_SECRET="your_secret_api_key"
```

```powershell
# Windows PowerShell
$env:PINATA_API_KEY="your_api_key"
$env:PINATA_API_SECRET="your_secret_api_key"
```

## Use

```bash
# Uploads a single file and prints only the gateway URL
svup path/to/file.jpg

# Uploads a single file and prints metadata
svup -v path/to/file.jpg
svup path/to/file.jpg -v
```

## Commands

```bash
# Upload
$ svup path/to/file.jpg
https://gateway.pinata.cloud/ipfs/Qm...

# Verbose Upload
$ svup -v path/to/file.jpg
Upload Result:
  Success: true
  Hash: Qm...
  URL: https://gateway.pinata.cloud/ipfs/Qm...
  Filename: file.jpg
  Size: 12345 bytes
  Timestamp: 2026-06-03 12:00:00
```

## Info

`svup` uses the `svuplib` uploader library:

```text
https://github.com/derekhandy/svuplib
```

Only one file path is accepted per upload. Directories and multi-file uploads are not supported.

The CLI validates Pinata credentials before uploading and exits early if `PINATA_API_KEY` or `PINATA_API_SECRET` is missing.

## NOTICE

<b> Files are uploaded to Pinata using the credentials available in the current shell environment. Treat uploaded files as externally hosted content and verify your Pinata gateway, account permissions, and file contents before sharing returned URLs.</b>
