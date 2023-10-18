echo "[running go mod tidy]"
go mod tidy

echo "[sorting imports]"
goimports-reviser ./...

echo "[compiling project]"
go build -o jalab-bot.exe cmd/main.go

echo "[running golangci-lint]"
golangci-lint run
