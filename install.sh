#!/bin/bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

echo -e "${CYAN}"}
echo "╔════════════════════════════════════════╗"
echo "║  NOXAR Installation Script              ║"
echo "║  Network Reconnaissance Tool            ║"
echo "╚════════════════════════════════════════╝"
echo -e "${NC}"

# Check if Go is installed
echo -e "${YELLOW}[*] Checking Go installation...${NC}"
if ! command -v go &> /dev/null; then
    echo -e "${RED}[!] Go is not installed!${NC}"
    echo -e "${YELLOW}[*] Install Go from: https://golang.org/doc/install${NC}"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}')
echo -e "${GREEN}[+] Go version: ${GO_VERSION}${NC}"

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    echo -e "${RED}[!] Error: go.mod not found!${NC}"
    echo -e "${YELLOW}[*] Make sure you're in the NOXAR root directory${NC}"
    exit 1
fi

# Build
echo -e "${YELLOW}[*] Building NOXAR...${NC}"
make build

# Install
echo ""
echo -e "${YELLOW}[*] Installing NOXAR globally...${NC}"
if [ "$EUID" -ne 0 ]; then 
    echo -e "${YELLOW}[*] Requesting sudo for global installation...${NC}"
    sudo make install
else
    make install
fi

echo ""
echo -e "${GREEN}╔════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║  Installation Complete!                 ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════╝${NC}"
echo ""
echo -e "${CYAN}[+] You can now use NOXAR from anywhere:${NC}"
echo -e "${GREEN}    $ noxar${NC}"
echo ""
echo -e "${CYAN}[+] For help:${NC}"
echo -e "${GREEN}    $ make help${NC}"
echo ""

# Test if noxar command works
if command -v noxar &> /dev/null; then
    echo -e "${GREEN}[+] Verification: noxar command is available!${NC}"
else
    echo -e "${YELLOW}[!] Warning: noxar not found in PATH${NC}"
    echo -e "${YELLOW}[*] Try: source ~/.bashrc${NC}"
fi
