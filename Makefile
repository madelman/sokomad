default: run

run:
	@echo "Running SokoMAD..."
	go run *.go

build:
	@echo "Building for current platform..."
	go build -o sokomad *.go

build-windows:
	@echo "Cross-building for Windows..."
	GOOS=windows GOARCH=amd64 go build -o sokomad.exe *.go
