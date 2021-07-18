package main

import (
	"encoding/csv"
	"image"
	"image/png"
	"log"
	"os"
	"strconv"
)

const LettersInAlphabet = 26

// Converts zero-indexed coordinates to letter-and-number excel ones
func excelIndex(x, y int) string {
	y++
	colIndex := ""
	for {
		colIndex = string(rune(x%LettersInAlphabet+'A')) + colIndex
		x /= LettersInAlphabet
		if x != 0 {
			x--
			continue
		}
		break
	}
	if colIndex == "" {
		colIndex = "A"
	}
	rowIndex := strconv.Itoa(y)
	return colIndex + rowIndex
}

func main() {
	if len(os.Args) != 3 {
		log.Fatalln("USAGE: img2excel <inputFile.png> <outputFile.csv>")
	}
	imgFilePath := os.Args[1]
	imageReader, err := os.Open(imgFilePath)
	if err != nil {
		log.Fatalf("%v: %v\n", imgFilePath, err)
	}
	img, err := png.Decode(imageReader)
	if err != nil {
		log.Fatal(err)
	}

	csvFilePath := os.Args[2]
	writer, err := os.Create(csvFilePath)
	if err != nil {
		log.Fatalf("%v: %v\n", csvFilePath, err)
	}

	matrix := matrixFromImage(img)
	saveAsCSV(matrix, writer)
}

func saveAsCSV(matrix [][]uint32, writer *os.File) {
	w := csv.NewWriter(writer)

	// strMatrix is also transposed
	strMatrix := make([][]string, len(matrix[0]))
	for i := 0; i < len(matrix[0]); i++ {
		strMatrix[i] = make([]string, len(matrix))
		for j := 0; j < len(matrix); j++ {
			strMatrix[i][j] = strconv.Itoa(int(matrix[j][i]))
		}
	}

	err := w.WriteAll(strMatrix)
	if err != nil {
		log.Fatalf("saveAsCSV: %v\n", err)
	}
}

func matrixFromImage(image image.Image) [][]uint32 {
	size := image.Bounds().Size()
	matrix := make([][]uint32, size.X)
	for i := range matrix {
		matrix[i] = make([]uint32, size.Y*3)
	}
	for x, column := range matrix {
		for y := 0; y < len(column); y += 3 {
			pixelColor := image.At(x, y/3)
			r, g, b, _ := pixelColor.RGBA()
			column[y] = r / 0xff
			column[y+1] = g / 0xff
			column[y+2] = b / 0xff
		}
	}
	return matrix
}
