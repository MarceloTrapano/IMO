package utils

import (
	"math"
)

func MatrixMax(matrix [][]int) (int, int, int) {
	max := math.MinInt64
	x := 0
	y := 0
	for i, row := range matrix {
		for j, value := range row {
			if value > max {
				max = value
				x = j
				y = i
			}
		}
	}
	return x, y, max
}
