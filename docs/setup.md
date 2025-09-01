# Development Environment Setup

## Prerequisites

- Go 1.21 or higher
- Git
- Chrome browser (for extension testing)
- curl or Postman (for API testing)

## Installation

### 1. Clone the Repository
```bash
git clone <repository-url>
cd monarchmoney-sync-backend
```

### 2. Install Dependencies
```bash
go mod download
go mod tidy
```

### 3. Set Up Environment Variables
```bash
cp .env.example .env
```

Edit `.env` and set your values:
```bash
PORT=8080
MONARCH_API_KEY=your-monarch-api-key
EXTENSION_SECRET_KEY=your-shared-secret
```

### 4. Run Tests (TDD Workflow)
```bash
# Run all tests
go test ./... -v

# Check coverage
go test ./... -cover

# Run specific test
go test -run TestHealthCheck -v
```

### 5. Run the Server
```bash
# Only run after all tests pass!
go run main.go
```

The server will start on `http://localhost:8080`

## Development Workflow

### TDD Development Flow
1. Write a test first
2. Run the test - watch it fail
3. Write minimal code to pass
4. Run the test - watch it pass
5. Refactor if needed
6. Update documentation

### Adding New Features
1. Create a test file: `feature_test.go`
2. Write failing tests
3. Implement the feature
4. Ensure all tests pass
5. Update `/docs/progress.md`
6. Update API documentation if needed

## Project Structure

```
monarchmoney-sync-backend/
├── main.go                 # Application entry point
├── main_test.go           # Main integration tests
├── go.mod                 # Go module definition
├── go.sum                 # Dependency checksums
├── .env.example           # Environment variables template
├── .gitignore            # Git ignore rules
├── README.md             # Project overview
│
├── handlers/             # HTTP request handlers
│   ├── health.go        # Health check endpoint
│   ├── health_test.go   # Health check tests
│   ├── walmart.go       # Walmart order handlers
│   └── walmart_test.go  # Walmart handler tests
│
├── models/              # Data models
│   ├── order.go        # Order structures
│   └── order_test.go   # Model tests
│
├── services/           # Business logic
│   ├── monarch.go      # Monarch Money integration
│   ├── monarch_test.go # Monarch service tests
│   ├── processor.go    # Order processing logic
│   └── processor_test.go # Processor tests
│
├── config/            # Configuration
│   ├── config.go     # Config loading
│   └── config_test.go # Config tests
│
├── testdata/         # Test fixtures
│   └── sample_order.json # Sample order data
│
└── docs/            # Documentation
    ├── progress.md  # Development progress
    ├── api.md      # API documentation
    ├── testing.md  # Testing guidelines
    ├── setup.md    # This file
    └── bug-fixes.md # Bug fix log
```

## Testing the API

### Health Check
```bash
curl http://localhost:8080/health
```

### Send Walmart Order
```bash
curl -X POST http://localhost:8080/api/walmart/orders \
  -H "Content-Type: application/json" \
  -H "X-Extension-Key: your-secret-key" \
  -d @testdata/sample_order.json
```

## Debugging

### Enable Debug Logging
Set Gin to debug mode:
```go
gin.SetMode(gin.DebugMode)
```

### View Logs
```bash
go run main.go 2>&1 | tee server.log
```

### Common Issues

#### Port Already in Use
```bash
# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>
```

#### Module Download Issues
```bash
# Clear module cache
go clean -modcache

# Re-download dependencies
go mod download
```

#### Test Failures
```bash
# Run with verbose output
go test ./... -v

# Run with race detection
go test ./... -race
```

## IDE Setup

### VS Code
Install extensions:
- Go (official)
- Go Test Explorer
- REST Client (for API testing)

Settings:
```json
{
  "go.testOnSave": true,
  "go.coverOnSave": true,
  "go.lintOnSave": "package",
  "go.formatTool": "goimports"
}
```

### GoLand
- Enable Go Modules support
- Configure test runner for TDD
- Set up file watchers for auto-formatting

## Linting and Formatting

### Format Code
```bash
go fmt ./...
goimports -w .
```

### Run Linter
```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run
```

## Database Setup (Future)

### PostgreSQL
```bash
# Create database
createdb walmart_monarch_sync

# Run migrations
migrate -path migrations -database "postgresql://localhost/walmart_monarch_sync" up
```

### Redis
```bash
# Start Redis
redis-server

# Test connection
redis-cli ping
```

## Deployment (Future)

### Docker
```bash
# Build image
docker build -t walmart-monarch-sync .

# Run container
docker run -p 8080:8080 --env-file .env walmart-monarch-sync
```

### Production Checklist
- [ ] Set Gin to release mode
- [ ] Configure proper logging
- [ ] Set up monitoring
- [ ] Configure rate limiting
- [ ] Enable HTTPS
- [ ] Set up database backups
- [ ] Configure secrets management