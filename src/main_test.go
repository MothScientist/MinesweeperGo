package main

import (
	"math/rand"
	"testing"
)

func TestFillSquare(t *testing.T) {
	minesCount := rand.Intn(2) + 1
	square := fillSquare(minesCount)

	if len(square) != 3 {
		t.Errorf("Invalid field size: %d", len(square))
	}
	for i := range 3 {
		if len(square[i]) != 3 {
			t.Errorf("Invalid field size: %d, element: %d", len(square[i]), i)
		}
	}

	minesFound := 0

	for i := range square {
		for j := range square[i] {
			if square[i][j] >= 0 && square[i][j] <= 2 {
				continue
			} else if square[i][j] == -3 {
				minesFound++
			} else {
				t.Errorf("Invalid element: %d, position [i][j] = [%d][%d]", square[i][j], i, j)
			}
		}
	}

	if minesFound != minesCount {
		t.Errorf("minesFound != minesCount. minesFound: %d, , minesCount: %d", minesFound, minesCount)
	}

}
