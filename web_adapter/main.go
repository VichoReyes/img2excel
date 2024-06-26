package main

import (
	"bytes"
	"fmt"
	"syscall/js"

	"img2excel/converter"
)

func main() {
	keepAlive := make(chan bool)
	convertToExcel := js.FuncOf(convertToExcel)
	js.Global().Set("convertToExcel", convertToExcel)
	fmt.Println("WebAssembly initialized")
	<-keepAlive
}

func convertToExcel(this js.Value, p []js.Value) any {
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
	// create UInt8Array from dest.Bytes()
	jsExcelBytes := js.Global().Get("Uint8Array").New(len(dest.Bytes()))
	js.CopyBytesToJS(jsExcelBytes, dest.Bytes())
	fmt.Println("convertToExcel: conversion successful")
	return jsExcelBytes
}
