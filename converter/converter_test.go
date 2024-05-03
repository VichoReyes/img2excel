package converter

import (
	"testing"
	"testing/quick"
)

func TestExcelIndex(t *testing.T) {
	cases := []struct {
		inX  int
		inY  int
		want string
	}{
		{0, 0, "A1"},
		{25, 2, "Z3"},
		{26, 13, "AA14"},
		{26 + 25, 0, "AZ1"},
		{26 + 3, 0, "AD1"},
		{751, 0, "ABX1"},
	}
	for _, c := range cases {
		got := excelIndex(c.inX, c.inY)
		if got != c.want {
			t.Errorf("excelIndex(%q, %q) == %q, want %q", c.inX, c.inY, got, c.want)
		}
	}

}

// // makes matrix maxImageWidth pixels wide at most
// func limitWidth(matrix [][][3]uint8, maxImageWidth int) [][][3]uint8 {

func TestLimitWidth(t *testing.T) {
	// we'll test a single case
	// resizing a 5x5 matrix to 2x2
	original := [][][3]uint8{
		{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10, 11, 12}, {13, 14, 15}},
		{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10, 11, 12}, {13, 14, 15}},
		{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10, 11, 12}, {13, 14, 15}},
		{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10, 11, 12}, {13, 14, 15}},
		{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10, 11, 12}, {13, 14, 15}},
	}
	resized := limitWidth(original, 2)
	expected := [][][3]uint8{
		{{4, 5, 6}, {10, 11, 12}},
		{{4, 5, 6}, {10, 11, 12}},
	}
	if len(resized) != len(expected) {
		t.Errorf("limitWidth: unexpected result length: got %v want %v", len(resized), len(expected))
	}
	for i, v := range resized {
		for j, u := range v {
			if u != expected[i][j] {
				t.Errorf("limitWidth: unexpected result at %d,%d: got %v want %v", i, j, u, expected[i][j])
			}
		}
	}
}

// check that limitWidth doesn't panic
// and that the resulting size is reasonable
func TestLimitWidthProperties(t *testing.T) {
	prop := func(x, y, maxImageWidth uint8) bool {
		if x == 0 || y == 0 || maxImageWidth == 0 {
			return true
		}
		matrix := make([][][3]uint8, x)
		for i := 0; i < int(x); i++ {
			matrix[i] = make([][3]uint8, y)
		}
		resized := limitWidth(matrix, int(maxImageWidth))
		if len(resized) == 0 || len(resized[0]) == 0 {
			t.Error("limitWidth: resized shouldn't be empty")
		}
		if len(resized) > int(x) || len(resized[0]) > int(y) {
			t.Error("limitWidth: resized shouldn't be bigger than original")
		}
		stride := (x + maxImageWidth - 1) / maxImageWidth
		if int(x)-int(stride) >= len(resized)*int(stride) {
			t.Error("limitWidth: resized too small")
		}
		if len(resized)*int(stride) > int(x) {
			t.Error("limitWidth: resized too big")
		}
		return true
	}
	if err := quick.Check(prop, nil); err != nil {
		t.Error("limitWidth: quick.Check failed")
	}
}
