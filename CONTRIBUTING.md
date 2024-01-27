# Contributing

## Getting Started

1. Fork the repository on GitHub
2. Clone the forked repository to your local machine
3. Build the project: `go build cmd/sepp/main.go`
4. Run the project: `./main --config=config.yaml`

## Testing

### Unit Tests

```bash
go test ./...
```

### Lint

```bash
golangci-lint run ./...
```
