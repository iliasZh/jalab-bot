echo "[formatting code]"
gofumpt -w .

echo "[sorting imports]"
goimports-reviser ./...

echo "[running go mod tidy]"
go mod tidy

echo "[compiling project]"
go build -o jalab-bot cmd/main.go

echo "[running golangci-lint]"
golangci-lint run
