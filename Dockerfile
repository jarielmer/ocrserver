FROM golang:1.24-bullseye AS builder

WORKDIR /app

# Install Tesseract development files required for gosseract
RUN apt-get update && apt-get install -y \
    libtesseract-dev \
    libleptonica-dev \
    pkg-config

COPY . .
RUN go mod tidy && \
    go build -o ocr-server .

FROM debian:bullseye-slim

# Install runtime dependencies and language data for English, German, French, and Italian
RUN apt-get update && apt-get install -y \
    libtesseract4 \
    tesseract-ocr \
    tesseract-ocr-eng \
    tesseract-ocr-deu \
    tesseract-ocr-fra \
    tesseract-ocr-ita \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/ocr-server .

# Set the port environment variable with a default value
ENV PORT=8080

# Expose the port
EXPOSE ${PORT}

# Run the OCR server
CMD ["./ocr-server"]
