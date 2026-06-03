# svup

`svup` uploads a file to IPFS through Pinata and prints the CID and gateway URL.

It is a small command-line program written in Go. Give it a file path, and it will:

- check that your Pinata API keys work
- upload the file to Pinata
- print the IPFS url for the uploaded file (-v : print all metadata)

## Requirements

- Go 1.26 or newer, if you are building it yourself
- a Pinata account
- a Pinata API key and secret

You can create Pinata API keys here:

```text
https://app.pinata.cloud/keys
```

## Build

From the repository root:

```bash
cd src
go build -o svup
```

This creates a `svup` binary in the `src` directory.

## Set up your Pinata keys

`svup` reads your Pinata credentials from environment variables:

```bash
export PINATA_API_KEY="your_api_key"
export PINATA_API_SECRET="your_secret_api_key"
```

On Windows PowerShell:

```powershell
$env:PINATA_API_KEY="your_api_key"
$env:PINATA_API_SECRET="your_secret_api_key"
```

Do not commit these keys to the repository.

## Use

Run `svup` with the file you want to upload:

```bash
./svup path/to/file.jpg
```

If `svup` is already on your `PATH`, you can run:

```bash
svup path/to/file.jpg
```

The output looks like this:

```text
Testing connection to Pinata...
Successfully connected to Pinata IPFS API

Uploading photo: path/to/file.jpg
Upload Result:
  Success: true
  Hash: Qm...
  URL: https://gateway.pinata.cloud/ipfs/Qm...
  Filename: file.jpg
  Size: 12345 bytes
  Timestamp: 2026-06-03 12:00:00
```

## Notes

`svup` uses Pinata's public gateway:

```text
https://gateway.pinata.cloud/ipfs/
```

The program currently uploads one file at a time. It uses the file name from the path you pass in.