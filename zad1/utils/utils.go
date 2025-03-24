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

func Insert(array []int, i int, j int) []int {
	var new_arr []int
	for idx, val := range array {
		if idx == i {
			new_arr = append(new_arr, j)
		}
		new_arr = append(new_arr, val)
	}
	return new_arr
}
