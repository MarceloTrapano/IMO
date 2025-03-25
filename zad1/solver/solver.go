package solver

import (
	"fmt"
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

func PickRandomNode(nodes []reader.Node) (int, error) {
	node1 := rand.Intn(len(nodes))
	return node1, nil
}

func PickRandomFarthest(distance_matrix [][]int, nodes []reader.Node, visited []bool) (int, int, error) {
	node1, err := PickRandomNode(nodes)
	if err != nil {
		return -1, -1, err
	}
	node2, err := utils.FarthestNode(nodes, distance_matrix, node1, visited)
	if err != nil {
		return -1, -1, err
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

func ValidateOrder(order [][]int, nodes []reader.Node) error {
	var visited []bool = make([]bool, len(nodes))
	if len(order[0])+len(order[1]) < len(nodes) {
		return fmt.Errorf("not all nodes visited")
	}
	for i := range order {
		for j := range order[i] {
			if visited[order[i][j]] {
				return fmt.Errorf("node %v visited more than once", order[i][j])
			}
			visited[order[i][j]] = true
		}
	}
	return nil
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
	err = ValidateOrder(order, nodes)
	if err != nil {
		return order, err
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
			switch {
			case min_1 == -1:
				min_1 = order1_nn
				order[0][j] = i
			case min_2 == -1:
				min_2 = order2_nn
				order[1][j] = i
			case order1_nn < min_1:
				min_1 = order1_nn
				order[0][j] = i
			case order2_nn < min_2:
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
					cost = 0                                // występowały leaki pamięci i program odpierdalał
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
	var (
		visited   []bool                  = make([]bool, len(nodes))                  // tablica dodanych wierzchołków
		cycles    []*utils.Edge           = make([]*utils.Edge, NumCycles)            // cykle krawędzi
		distances []*utils.EdgeLinkedList = make([]*utils.EdgeLinkedList, len(nodes)) // tablica linked list z dystansami do krawędzi
		edges     []utils.Edge                                                        // tablica wszystkich krawędzi
	)

	start_node_1, start_node_2, _ := PickRandomFarthest(distance_matrix, nodes, visited) // wybór startowych punktów

	visited[start_node_1] = true
	visited[start_node_2] = true

	// stworzenie pierwszych krawędzi na podstawie najbliższych sąsiadów
	nearest1, err := utils.NearestNode(nodes, distance_matrix, start_node_1, visited)
	if err != nil {
		return err
	}
	visited[nearest1] = true
	nearest2, err := utils.NearestNode(nodes, distance_matrix, start_node_2, visited)
	if err != nil {
		return err
	}
	visited[nearest2] = true

	edge1 := utils.NewEdge(start_node_1, nearest1, distance_matrix, nil, nil)
	edge2 := utils.NewEdge(start_node_2, nearest2, distance_matrix, nil, nil)
	edges = append(edges, edge1, edge2)
	cycles[0] = &edge1
	cycles[1] = &edge2

	// aktualizacja dystansów
	for i := range distances {
		if visited[i] {
			continue
		}
		distances[i] = &utils.EdgeLinkedList{Edge: 0, Next: nil, Value: utils.EdgeInsertValue(distance_matrix, i, &edges[0])}
		newEdges := []utils.EdgeLinkedList{{Edge: 1, Next: nil, Value: utils.EdgeInsertValue(distance_matrix, i, &edges[1])}}
		distances[i] = utils.UpdateDistances(distances[i], distance_matrix, newEdges, false)
	}

	return nil
}
