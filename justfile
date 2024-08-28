default:
    @just --list

# Builds the lsdeps binary
build:
    @go build -o lsdeps lsdeps.go
    @echo "Built lsdeps"

# Installs the lsdeps binary with Go
go-install:
    @go install lsdeps.go
    @echo "Installed lsdeps with Go"

# Installs the lsdeps binary to ~/.local/bin
install: build
    @install lsdeps ~/.local/bin
    @rm lsdeps
    @echo "Installed lsdeps to ~/.local/bin"
