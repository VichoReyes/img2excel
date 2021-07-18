package main

import (
	"testing"
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
	}
	for _, c := range cases {
		got := excelIndex(c.inX, c.inY)
		if got != c.want {
			t.Errorf("excelIndex(%q, %q) == %q, want %q", c.inX, c.inY, got, c.want)
		}
	}

}
