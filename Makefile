run:
	go run .

serve:
	@echo "Hosting game on http://localhost:4242"
	wasmserve --http=":4242" --allow-origin='*' .

WASM_EXEC_PATH="$(shell go env GOROOT)/misc/wasm/wasm_exec.js"

itch:
	GOOS=js GOARCH=wasm go build -o game.wasm .
	zip -r -j site.zip game.wasm $(WASM_EXEC_PATH) html/index.html
	rm game.wasm

clean:
	rm -f game.wasm site.zip
