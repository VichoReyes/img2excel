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

	matrix := matrixFromImage(img)
	saveAsExcel(matrix, excelWriter)
	return nil
}

func saveAsExcel(matrix [][]uint32, writer io.Writer) {
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
