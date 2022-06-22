default:
    @just --list

# Format and lint project
fmt:
    treefmt

# Build the project
build:
    go build .

# Run linters not covered by treefmt
lint:
    golangci-lint run
