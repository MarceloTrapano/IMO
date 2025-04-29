package solver

import (
	"IMO/reader"
	"IMO/utils"
	"fmt"
	"math"
	"math/rand"
)

const (
	NumCycles int     = 2
	Split     float64 = 0.5
)

func Solve(nodes []reader.Node, algorithm string, distance_matrix [][]int) ([][]int, error) {
	var (
		order           [][]int = make([][]int, NumCycles) // kolejność odwiedzania wierzchołków dla obydwu cykli
		nodes_cycle_one int                                // liczba wierzchołków w cyklu 1
	)
	// stworzenie macierzy odległości

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
	case "wreg": // weighted regret - żal ważony
		f = WeightedRegret
	case "rand": // rozwiązanie losowe
		f = Random
	default:
		f = InOrder
	}

	err := f(distance_matrix, order, nodes)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func Local_search(start_order [][]int, algorithm string, distance_matrix [][]int) ([][]int, error) {
	var order [][]int = make([][]int, NumCycles)
	copy(order, start_order)
	order = append(start_order[:0:0], start_order...)
	var f func([][]int, [][]int) error
	switch algorithm {
	case "sn":
		f = SteepestNode
	case "se":
		f = SteepestEdge
	case "gn":
		f = GreedyNode
	case "ge":
		f = GreedyEdge
	case "rw":
		f = RandomWalk
	case "fls":
		f = FastLocalSearch
	case "c":
		f = CandidateSearch
	default:
		f = SteepestEdge
	}
	err := f(distance_matrix, order)
	if err != nil {
        return nil, err
    }
	return order, nil
}

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

func PickRandomFarthest(distance_matrix [][]int, nodes []reader.Node) (int, int, error) {
	visited := make([]bool, len(nodes))
	node1, err := PickRandomNode(nodes)
	visited[node1] = true
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
