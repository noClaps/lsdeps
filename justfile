_default:
    @just --list

# Builds the lsdeps binary
build:
	@swift build -c release
	@cp ./.build/release/lsdeps .
	@echo "Built lsdeps"

# Runs the lsdeps project with the given arguments
run ARGS:
    @swift run lsdeps {{ARGS}}

# Installs the lsdeps binary to ~/.local/bin
install: build
    @install ./lsdeps ~/.local/bin
    @rm ./lsdeps
    @echo "Installed lsdeps to ~/.local/bin"

# Format all Swift files
format:
    @swift-format format -i ./**/*.swift
