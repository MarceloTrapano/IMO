package solver

import (
	"math"
	"math/rand"
	"zad1/reader"
	"zad1/utils"
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
func PickFarthestNodes(distance_matrix [][]int, nodes []reader.Node) (int, int, error) {
	x, y, _ := utils.MatrixMax(distance_matrix)
	return x, y, nil
}
func PickRandomNodes(nodes []reader.Node) (int, int, error) {
	node1 := rand.Intn(len(nodes))
	node2 := node1
	for node1 == node2 {
		node2 = rand.Intn(len(nodes))
	}
	return node1, node2, nil
}
func PickRandomClosestNodes(distance_matrix [][]int, nodes []reader.Node) (int, int, error) {
	idx := rand.Intn(len(nodes))
	node_val := 10000
	node2_idx := -1
	for i, val := range distance_matrix[idx] {
		if val == 0 {
			continue
		}
		if val < node_val {
			node_val = val
			node2_idx = i
		}
	}
	return idx, node2_idx, nil
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
	var f func([][]int, [][]int, []reader.Node) error
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

	err := f(distance_matrix, order, nodes)
	if err != nil {
		return nil, err
	}
	return order, nil
}

// testowo jak może struktura wyglądać funkcji - paramtetry
func InOrder(distance_matrix [][]int, order [][]int, nodes []reader.Node) error {
	for i := range distance_matrix {
		if i < len(order[0]) {
			order[0][i] = i
		} else {
			order[1][i-len(order[0])] = i
		}
	}
	return nil
}

func NearestNeighbour(distance_matrix [][]int, order [][]int, nodes []reader.Node) error {
	start_node_1, start_node_2, _ := PickRandomClosestNodes(distance_matrix, nodes) // wybór startowych punktów

	order[0][len(order[0])-1] = -1
	order[1][len(order[1])-1] = -1 // znakowanie końca tablic order
	order[0][0] = start_node_1
	order[1][0] = start_node_2 // przypisywanie pierwszych wierzchołków

	var visited []bool = make([]bool, len(nodes)) // tablica dodanych wierzchołków

	visited[start_node_1] = true
	visited[start_node_2] = true
	for j := 1; order[0][len(order[0])-1] == -1 || order[1][len(order[1])-1] == -1; j++ {
		min_1 := -1
		min_2 := -1 // wartości najbliższych krawędzi
		if j < len(order[0]) {
			order[0][j] = -1
		}
		if j < len(order[1]) {
			order[1][j] = -1
		}
		for i := range nodes {
			if visited[i] {
				continue // pomijanie wierzchołków dodanych
			}

			if j >= len(order[0]) { // po osiągnięciu maksymalnej długości na jednycm cyklu resztę sąsiadów szuka dla jednego cyklu
				order2_nn := distance_matrix[i][order[1][j-1]]
				if min_2 == -1 || order2_nn < min_2 {
					min_2 = order2_nn
					order[1][j] = i
				}
				continue
			}
			if j >= len(order[1]) {
				order1_nn := distance_matrix[i][order[0][j-1]]
				if min_1 == -1 || order1_nn < min_1 {
					min_1 = order1_nn
					order[0][j] = i
				}
				continue
			}
			order1_nn := distance_matrix[i][order[0][j-1]]
			order2_nn := distance_matrix[i][order[1][j-1]]
			if min_1 == -1 || (order1_nn < min_1) {
				min_1 = order1_nn
				order[0][j] = i
				if order[1][j] != -1 { // warunek sprawdzający czy drugi cykl ma przypisaną wartość
					continue
				}
			}
			if min_2 == -1 || (order2_nn < min_2) {
				min_2 = order2_nn
				order[1][j] = i
			}
		}
		if j < len(order[0]) {
			visited[order[0][j]] = true
		}
		if j < len(order[1]) {
			visited[order[1][j]] = true
		}
	}
	return nil
}

func GreedyCycle(distance_matrix [][]int, order [][]int, nodes []reader.Node) error {
	start_node_1, start_node_2, _ := PickRandomNodes(nodes) // wybór startowych punktów

	var (
		visited      []bool = make([]bool, len(nodes)) // tablica dodanych wierzchołków
		cycle1       []int
		cycle2       []int
		new_cycle    []int
		visit        int
		cost         int
		temp_cycle   []int
		minimal_cost int
	)

	visited[start_node_1] = true
	visited[start_node_2] = true

	cycle1 = append(cycle1, start_node_1)
	cycle2 = append(cycle2, start_node_2)

	lenCycle1 := len(order[0])
	lenCycle2 := len(order[1])

	for len(cycle1) < lenCycle1 || len(cycle2) < lenCycle2 {
		if len(cycle1) < lenCycle1 {
			visit = -1
			minimal_cost = -1
			for i := range nodes {
				if visited[i] {
					continue
				}
				for j := range cycle1 {
					temp_cycle = utils.Insert(cycle1, j, i) // musiałem sam napisać funkcję do dodwawania elementu do macierzy XD
					cost = 0								// występowały leaki pamięci i program odpierdalał
					for node_idx := range temp_cycle {
						cost += distance_matrix[temp_cycle[node_idx]][temp_cycle[(node_idx+1)%len(temp_cycle)]]
					}
					if minimal_cost == -1 || cost < minimal_cost {
						new_cycle = append(temp_cycle[:0:0], temp_cycle...)
						minimal_cost = cost
						visit = i
					}
				}
			}
			
			cycle1 = append(new_cycle[:0:0], new_cycle...)
			visited[visit] = true

		}
		if len(cycle2) < lenCycle2 {
			visit = -1
			minimal_cost = -1
			for i := range nodes {
				if visited[i] {
					continue
				}
				for j := range cycle2 {
					temp_cycle = utils.Insert(cycle2, j, i)
					cost = 0
					for node_idx := range temp_cycle {
						cost += distance_matrix[temp_cycle[node_idx]][temp_cycle[(node_idx+1)%len(temp_cycle)]]
					}
					if minimal_cost == -1 || cost < minimal_cost {
						new_cycle = append(temp_cycle[:0:0], temp_cycle...)
						minimal_cost = cost
						visit = i
					}
				}
			}
			cycle2 = append(new_cycle[:0:0], new_cycle...)
			visited[visit] = true
		}
	}
	order[0] = cycle1
	order[1] = cycle2

	return nil
}

func Regret(distance_matrix [][]int, order [][]int, nodes []reader.Node) error {
	// TODO
	return nil
}
