# Go compiler
GO := go

# Output file names
WASM_OUTPUT := main.wasm
SERVER_OUTPUT := server-bin
WASM_EXEC_OUTPUT := wasm_exec.js

# Build target
build: build-server get-wasm-script
	GOOS=js GOARCH=wasm $(GO) build -o $(WASM_OUTPUT) ./web_adapter/main.go

# Build server target
build-server:
	$(GO) build -o $(SERVER_OUTPUT) ./server/main.go

# Build wasm_exec.js target
get-wasm-script:
	cp "$(shell go env GOROOT)/misc/wasm/$(WASM_EXEC_OUTPUT)" .

# Clean target
clean:
	rm -f $(WASM_OUTPUT) $(SERVER_OUTPUT) $(WASM_EXEC_OUTPUT)
