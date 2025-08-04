#!/bin/bash

# format.sh - Comprehensive Go code formatting script
# Usage: ./scripts/format.sh [options]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to install missing tools
install_tools() {
    print_status "Installing missing Go formatting tools..."
    
    if ! command_exists goimports; then
        go install golang.org/x/tools/cmd/goimports@latest
    fi
    
    if ! command_exists ineffassign; then
        go install github.com/gordonklaus/ineffassign@latest
    fi
    
    if ! command_exists misspell; then
        go install github.com/client9/misspell/cmd/misspell@latest
    fi
    
    if ! command_exists staticcheck; then
        go install honnef.co/go/tools/cmd/staticcheck@latest
    fi
    
    print_success "All tools installed!"
}

# Function to run formatting
run_formatting() {
    print_status "Starting Go code formatting..."
    
    # 1. gofmt - Standard Go formatter
    print_status "Running gofmt (standard Go formatter)..."
    if command_exists gofmt; then
        gofmt -w -s .
        print_success "gofmt completed"
    else
        print_error "gofmt not found"
        exit 1
    fi
    
    # 2. goimports - Import management
    print_status "Running goimports (import management)..."
    if command_exists goimports; then
        goimports -w .
        print_success "goimports completed"
    else
        print_warning "goimports not found, installing..."
        go install golang.org/x/tools/cmd/goimports@latest
        goimports -w .
        print_success "goimports completed"
    fi
    
    # 3. go mod tidy - Clean up dependencies
    print_status "Running go mod tidy..."
    go mod tidy
    print_success "go mod tidy completed"
    
    # 4. go vet - Basic static analysis
    print_status "Running go vet (static analysis)..."
    go vet ./...
    print_success "go vet completed"
    
    # 5. ineffassign - Find ineffectual assignments
    print_status "Checking for ineffectual assignments..."
    if command_exists ineffassign; then
        ineffassign ./...
        print_success "ineffassign check completed"
    else
        print_warning "ineffassign not found, skipping..."
    fi
    
    # 6. misspell - Fix common misspellings
    print_status "Checking for misspellings..."
    if command_exists misspell; then
        misspell -w .
        print_success "misspell check completed"
    else
        print_warning "misspell not found, skipping..."
    fi
    
    # 7. staticcheck - Advanced static analysis (optional)
    if [[ "$ADVANCED" == "true" ]]; then
        print_status "Running staticcheck (advanced static analysis)..."
        if command_exists staticcheck; then
            staticcheck ./...
            print_success "staticcheck completed"
        else
            print_warning "staticcheck not found, skipping..."
        fi
    fi
    
    print_success "All formatting completed successfully!"
}

# Function to show help
show_help() {
    echo "Go Code Formatting Script"
    echo ""
    echo "Usage: $0 [options]"
    echo ""
    echo "Options:"
    echo "  --install     Install missing formatting tools"
    echo "  --advanced    Run advanced static analysis (staticcheck)"
    echo "  --help        Show this help message"
    echo ""
    echo "Default behavior: Run standard formatting (gofmt, goimports, go vet, etc.)"
    echo ""
    echo "Tools used:"
    echo "  - gofmt: Standard Go formatter"
    echo "  - goimports: Import management"
    echo "  - go mod tidy: Dependency cleanup"
    echo "  - go vet: Basic static analysis"
    echo "  - ineffassign: Find ineffectual assignments"
    echo "  - misspell: Fix common misspellings"
    echo "  - staticcheck: Advanced static analysis (with --advanced)"
}

# Parse command line arguments
INSTALL_ONLY=false
ADVANCED=false

for arg in "$@"; do
    case $arg in
        --install)
            INSTALL_ONLY=true
            shift
            ;;
        --advanced)
            ADVANCED=true
            shift
            ;;
        --help)
            show_help
            exit 0
            ;;
        *)
            print_error "Unknown option: $arg"
            show_help
            exit 1
            ;;
    esac
done

# Main execution
print_status "Go Code Formatter Script"
print_status "========================"

if [[ "$INSTALL_ONLY" == "true" ]]; then
    install_tools
    exit 0
fi

# Check if we're in a Go project
if [[ ! -f "go.mod" ]]; then
    print_error "No go.mod file found. Please run this script from the root of a Go project."
    exit 1
fi

run_formatting

echo ""
print_success "🎉 Code formatting completed! Your Go code is now properly formatted."
echo ""
print_status "Next steps:"
echo "  - Review the changes with: git diff"
echo "  - Run tests to ensure everything works: go test ./..."
echo "  - Consider setting up pre-commit hooks for automatic formatting" 