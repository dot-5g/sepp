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

## Container image

```bash
rockcraft pack -v
version=$(yq '.version' rockcraft.yaml)
sudo skopeo --insecure-policy copy oci-archive:sepp_${version}_amd64.rock docker-daemon:sepp:${version}
docker run sepp:${version}
```
