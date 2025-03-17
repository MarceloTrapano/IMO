package solver

import (
	"math"
	"zad1/reader"
)

const (
	NumCycles int     = 2
	Split     float64 = 0.5
)

func EucDist(a, b reader.Node) int {
	return int(
		math.Round(
			math.Sqrt(
				math.Pow(float64(a.X-b.X), 2) + math.Pow(float64(a.Y-b.Y), 2),
			),
		), // fajna opcja żeby lepszą czytelność mieć ale śmieszne, że przecinki się daje przed nową linią żeby dobrze parsował kompilator
	)
}

// zwraca kolejnych indeksów wierzchołków z nodes w kolejności odwiedzania w cyklu
func Solve(nodes []reader.Node, algorithm string) ([][]int, error) {
	var (
		order           [][]int = make([][]int, NumCycles)  // kolejność odwiedzania wierzchołków dla obydwu cykli
		distance_matrix [][]int = make([][]int, len(nodes)) // macierz odległości
		nodes_cycle_one int                                 // liczba wierzchołków w cyklu 1
	)
	// stworzenie macierzy odległości
	for i := range distance_matrix {
		distance_matrix[i] = make([]int, len(nodes))
		for j := range distance_matrix[i] {
			distance_matrix[i][j] = EucDist(nodes[i], nodes[j])
		}
	}
	// zajęcie pamięci dla macierzy order
	nodes_cycle_one = int(float64(len(distance_matrix)) * Split)
	order[0] = make([]int, nodes_cycle_one)
	order[1] = make([]int, len(nodes)-nodes_cycle_one)

	// wybór algorytmu
	var f func([][]int, [][]int, int) error
	switch algorithm {
	case "nn": // nearest neighbour - najbliższy sąsiad
		f = NearestNeighbour
	case "gc": // greedy cycle - rozbudowa cyklu
		f = GreedyCycle
	case "reg": // regret - żal
		f = Regret
	default:
		f = InOrder
	}

	err := f(distance_matrix, order, nodes_cycle_one)
	if err != nil {
		return nil, err
	}
	return order, nil
}

// testowo jak może struktura wyglądać funkcji - paramtetry
func InOrder(distance_matrix [][]int, order [][]int, nodes_cycle_one int) error {
	for i := range distance_matrix {
		if i < nodes_cycle_one {
			order[0][i] = i
		} else {
			order[1][i-nodes_cycle_one] = i
		}
	}
	return nil
}

func NearestNeighbour(distance_matrix [][]int, order [][]int, nodes_cycle_one int) error {
	// TODO
	return nil
}

func GreedyCycle(distance_matrix [][]int, order [][]int, nodes_cycle_one int) error {
	// TODO
	return nil
}

func Regret(distance_matrix [][]int, order [][]int, nodes_cycle_one int) error {
	// TODO
	return nil
}
