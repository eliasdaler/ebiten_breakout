package main

import "syscall/js"

// this makes sure that focus is enabled by default when running on itch.io
func focus() {
	doc := js.Global().Get("document")
	canvas := doc.Call("getElementsByTagName", "canvas").Index(0)
	canvas.Call("focus")
}
