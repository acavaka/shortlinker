# URL Shortener Service

A robust and efficient URL shortening service written in Go. This service provides both in-memory and file-based storage options for URL mappings.

## Features

- Generate short URLs from long URLs
- Support for both in-memory and file-based storage
- Thread-safe operations
- Gzip compression support
- Logging middleware
- RESTful API endpoints

## Installation

### Prerequisites

- Go 1.23.4 or higher
- Git

### Getting Started

1. Clone the repository:
```bash
git clone https://github.com/acavaka/shortlinker.git
cd shortlinker
```

2. Install dependencies:
```bash
go mod download
```

3. Build the project:
```bash
go build -o shortlinker ./cmd/shortlinker
```

## Usage

### Running the Service

```bash
./shortlinker
```

By default, the service runs in memory-only mode. To enable file storage, set the appropriate configuration.

### API Endpoints

1. Create Short URL
```http
POST /api/shorten
Content-Type: application/json

{
    "url": "https://example.com/very/long/url"
}
```

2. Get Original URL
```http
GET /{shortURL}
```

3. Save URL (Alternative Endpoint)
```http
POST /
Content-Type: application/json

{
    "url": "https://example.com/very/long/url"
}
```

## Configuration

The service can be configured using environment variables or a configuration file. Key configuration options include:

- `FILE_STORAGE_PATH`: Path to the file storage (optional)
- `BASE_URL`: Base URL for the shortened links
- Other configuration options can be found in the `internal/config` package

## Development

### Project Structure

```
.
├── cmd/
│   ├── shortlinker/    # Main application
│   └── shortener/      # CLI tool
├── internal/
│   ├── config/         # Configuration
│   ├── handlers/       # HTTP handlers
│   ├── logger/         # Logging
│   ├── middleware/     # HTTP middleware
│   ├── models/         # Data models
│   ├── service/        # Business logic
│   └── storage/        # Storage implementations
└── README.md
```

### Running Tests

To run all tests:
```bash
go test ./...
```

To run tests with coverage:
```bash
go test -cover ./...
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
