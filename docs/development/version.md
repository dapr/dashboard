# Versioning and releases

To release a versioned binary for Dapr Dashboard, simply add the following to the go build:
```bash
go build ... -ldflags="-X 'github.com/dapr/dashboard/pkg/version.Version=<your-version>"
```