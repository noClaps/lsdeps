_default:
	@echo "lsdeps"

# Builds the lsdeps binary
build:
	@go build
	@echo "Built lsdeps"

# Installs the lsdeps binary to ~/.local/bin
install: build
	@install ./lsdeps ~/.local/bin
	@rm ./lsdeps
	@echo "Installed lsdeps to ~/.local/bin"
