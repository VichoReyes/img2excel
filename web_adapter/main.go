package main

import "syscall/js"

func main() {
	keepAlive := make(chan bool)
	helloWorld := js.FuncOf(func(this js.Value, p []js.Value) interface{} {
		js.Global().Get("document").Call("getElementById", "output").Set("textContent", "Hello, WebAssembly!")
		return nil
	})
	js.Global().Set("helloWorld", helloWorld)
	<-keepAlive
}
