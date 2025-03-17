# OCR Server

This is a simple OCR (Optical Character Recognition) server written in Go. It accepts image uploads via a REST API and extracts text using Tesseract OCR.

## Features
- Supports OCR processing of images via HTTP
- Accepts multiple languages: English, German, French, and Italian
- Uses `gosseract` (a Go wrapper for Tesseract OCR)
- Simple and lightweight

## Requirements
- Go 1.18+
- Tesseract OCR installed on the system
- The `gosseract` package

### Installing Tesseract OCR
Tesseract OCR must be installed on your system for this server to work. Install it using:

**Ubuntu/Debian:**
```sh
sudo apt update && sudo apt install tesseract-ocr -y
```

**MacOS (Homebrew):**
```sh
brew install tesseract
```

**Windows:**
Download and install Tesseract from [UB Mannheim](https://github.com/UB-Mannheim/tesseract/wiki).

## Installation
1. Clone this repository:
```sh
git clone https://github.com/your-repo/ocr-server.git
cd ocr-server
```

2. Install dependencies:
```sh
go mod tidy
```

3. Run the server:
```sh
go run main.go
```
By default, the server runs on port `8080`.

## Usage
### Upload an Image for OCR
Send a `POST` request to `/ocr` with an image file and an optional language parameter.

#### Request
```sh
curl -X POST "http://localhost:8080/ocr" \
  -F "file=@path/to/image.png" \
  -F "lang=english"
```

#### Response
The response will be the extracted text from the image.

### Supported Languages
The server supports the following languages:
- English (`eng`)
- German (`deu`)
- French (`fra`)
- Italian (`ita`)

You can specify the language by using either the full name (`english`, `german`, etc.) or the Tesseract language code (`eng`, `deu`, etc.). If no language is provided, the default is English (`eng`).

## Docker Setup
A Docker image is available for running the OCR server in a containerized environment.

### Build the Docker Image
```sh
docker build -t ocr-server .
```

### Run the Container
```sh
docker run -p 8080:8080 ocr-server
```

### Dockerfile Details
The Docker image is built using a multi-stage approach:

#### 1. Build Stage
- Uses `golang:1.24-bullseye` as the base image
- Installs Tesseract development dependencies (`libtesseract-dev`, `libleptonica-dev`, `pkg-config`)
- Copies the source code and builds the Go binary

#### 2. Runtime Stage
- Uses `debian:bullseye-slim` as the base image
- Installs runtime dependencies (`libtesseract4`, `tesseract-ocr`, and language data for `eng`, `deu`, `fra`, `ita`)
- Copies the compiled binary from the build stage
- Exposes port `8080` (configurable via `PORT` environment variable)
- Runs the server using `CMD ["./ocr-server"]`

## Configuration
You can configure the server using environment variables:
- `PORT`: Change the default port (default: `8080`)

Example:
```sh
PORT=9090 go run main.go
```

## License
This project is licensed under the MIT License.

## Acknowledgments
- [gosseract](https://github.com/otiai10/gosseract) for Go bindings to Tesseract OCR
- Tesseract OCR for text extraction
