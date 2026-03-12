.PHONY: build run install clean help test

BINARY_NAME=noxar
BINARY_PATH=./$(BINARY_NAME)
INSTALL_PATH=/usr/local/bin/$(BINARY_NAME)
GO=go
GOFLAGS=-v

help:
	@echo "NOXAR - Network Reconnaissance Tool"
	@echo ""
	@echo "Available commands:"
	@echo "  make build      - Build binary (./noxar)"
	@echo "  make run        - Build and run"
	@echo "  make install    - Install globally (/usr/local/bin/noxar)"
	@echo "  make uninstall  - Remove from /usr/local/bin/"
	@echo "  make clean      - Remove binary"
	@echo "  make test       - Run tests"
	@echo "  make dev        - Build for development (with race detector)"

build:
	@echo "[*] Building NOXAR..."
	$(GO) build $(GOFLAGS) -o $(BINARY_PATH) ./cmd/gorecon
	@echo "[+] Binary created: $(BINARY_PATH)"
	@echo "[*] Usage: ./$(BINARY_NAME)"

dev:
	@echo "[*] Building with race detector..."
	$(GO) build -race -o $(BINARY_PATH) ./cmd/gorecon
	@echo "[+] Development build ready"

run: build
	@echo "[*] Running NOXAR..."
	$(BINARY_PATH)

install: build
	@echo "[*] Installing to $(INSTALL_PATH)..."
	@sudo cp $(BINARY_PATH) $(INSTALL_PATH)
	@sudo chmod +x $(INSTALL_PATH)
	@echo "[+] Installed successfully!"
	@echo "[*] Usage: Just type 'noxar' from anywhere"

uninstall:
	@echo "[*] Removing from $(INSTALL_PATH)..."
	@sudo rm -f $(INSTALL_PATH)
	@echo "[+] Uninstalled successfully!"

clean:
	@echo "[*] Cleaning up..."
	@rm -f $(BINARY_PATH)
	$(GO) clean
	@echo "[+] Clean complete"

test:
	@echo "[*] Running tests..."
	$(GO) test -v ./...
	@echo "[+] Tests completed"

coverage:
	@echo "[*] Running tests with coverage..."
	$(GO) test -v -cover ./...

fmt:
	@echo "[*] Formatting code..."
	$(GO) fmt ./...
	@echo "[+] Format complete"

lint:
	@echo "[*] Running linter..."
	@which golangci-lint > /dev/null || echo "[!] golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
	@golangci-lint run ./... 2>/dev/null || echo "[!] Lint check skipped"

version:
	@$(GO) version
	@./$(BINARY_NAME) --version 2>/dev/null || echo "[*] Run 'make build' first"

.DEFAULT_GOAL := help
