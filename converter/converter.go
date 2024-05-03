package converter

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
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

func Convert(imageReader io.Reader, excelWriter io.Writer) error {
	img, formatName, err := image.Decode(imageReader)
	if err != nil {
		return err
	}
	log.Printf("Decoded image format: %s", formatName)

	bigPxMatrix := matrixFromImage(img)
	smallerPxMatrix := limitWidth(bigPxMatrix, 120)
	numberMatrix := spreadY(smallerPxMatrix)
	saveAsExcel(numberMatrix, excelWriter)
	return nil
}

func saveAsExcel(matrix [][]uint8, writer io.Writer) {
	const sheet = "Sheet1"
	f := excelize.NewFile()
	for x, col := range matrix {
		for y, val := range col {
			coords := excelIndex(x, y)
			err := f.SetCellInt(sheet, coords, int(val))
			if err != nil {
				log.Fatalf("setting %s to be %d: %v", coords, val, err)
			}
		}
	}

	maxColors := []string{"#FF0000", "#00FF00", "#0000FF"}
	for y := 0; y < len(matrix[0]); y++ {
		excelRow := excelIndex(0, y) + ":" + excelIndex(len(matrix), y)
		err := f.SetConditionalFormat(sheet, excelRow,
			`[{"type":"2_color_scale","criteria":"=","min_type":"num","min_value":"0","max_type":"num","max_value":"255","min_color":"#000000","max_color":"`+maxColors[y%3]+`"}]`)
		if err != nil {
			log.Fatalf("SetConditionalFormat: %v", err)
		}
	}

	err := f.Write(writer)
	if err != nil {
		log.Fatalf("writing to excel writer file: %v", err)
	}
}

// turns each pixel into a [r, g, b] array
func matrixFromImage(image image.Image) [][][3]uint8 {
	size := image.Bounds().Size()
	matrix := make([][][3]uint8, size.X)
	for i := range matrix {
		matrix[i] = make([][3]uint8, size.Y)
	}
	for x, column := range matrix {
		for y := 0; y < len(column); y++ {
			pixelColor := image.At(x, y/3)
			r, g, b, _ := pixelColor.RGBA()
			column[y][0] = uint8(r / 0xff)
			column[y][1] = uint8(g / 0xff)
			column[y][2] = uint8(b / 0xff)
		}
	}
	return matrix
}

// turns each [r, g, b] array into three rows
func spreadY(imgMatrix [][][3]uint8) [][]uint8 {
	result := make([][]uint8, len(imgMatrix))
	for i := range result {
		result[i] = make([]uint8, len(imgMatrix[i])*3)
	}
	for x, column := range imgMatrix {
		for y, pixel := range column {
			result[x][y*3] = pixel[0]
			result[x][y*3+1] = pixel[1]
			result[x][y*3+2] = pixel[2]
		}
	}
	return result
}

// makes matrix maxImageWidth pixels wide at most
func limitWidth(matrix [][][3]uint8, maxImageWidth int) [][][3]uint8 {
	// round up of len(matrix) / maxImageWidth
	stride := (len(matrix) + maxImageWidth - 1) / maxImageWidth

	// averages each (stride x stride) square into a single pixel
	averageOf := func(x, y int) [3]uint8 {
		var r, g, b uint64
		for i := x; i < x+stride; i++ {
			for j := y; j < y+stride; j++ {
				r += uint64(matrix[i][j][0])
				g += uint64(matrix[i][j][1])
				b += uint64(matrix[i][j][2])
			}
		}
		area := uint64(stride * stride)
		return [3]uint8{uint8(r / area), uint8(g / area), uint8(b / area)}
	}

	result := make([][][3]uint8, len(matrix)/stride)
	for i := range result {
		result[i] = make([][3]uint8, len(matrix[0])/stride)
		for j := range result[i] {
			result[i][j] = averageOf(i*stride, j*stride)
		}
	}
	return result
}
