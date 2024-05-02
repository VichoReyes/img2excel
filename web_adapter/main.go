package main

import (
	"bytes"
	"fmt"
	"syscall/js"

	"img2excel/converter"
)

func main() {
	keepAlive := make(chan bool)
	convertToExcel := js.FuncOf(func(this js.Value, p []js.Value) any {
		if len(p) != 1 {
			fmt.Println("convertToExcel: expected 1 argument")
			return nil
		}
		jsImageBytes := p[0]
		byteLength := jsImageBytes.Get("byteLength").Int()
		buffer := make([]byte, byteLength)
		byteLengthCopied := js.CopyBytesToGo(buffer, jsImageBytes)
		if byteLength != byteLengthCopied {
			fmt.Printf("convertToExcel: expected to copy %d bytes, but got %d\n", byteLength, byteLengthCopied)
			return nil
		}
		fmt.Println("convertToExcel: converting image to excel")
		dest := new(bytes.Buffer)
		err := converter.Convert(bytes.NewReader(buffer), dest)
		if err != nil {
			fmt.Printf("convertToExcel converter: %v\n", err)
			return nil
		}
		return nil
	})
	js.Global().Set("convertToExcel", convertToExcel)
	fmt.Println("WebAssembly initialized") // not printing for some reason
	<-keepAlive
}
