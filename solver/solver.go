package solver

import (
	"fmt"
	"math"
	"math/rand"
	"IMO/reader"
	"IMO/utils"
)

const (
	NumCycles int     = 2
	Split     float64 = 0.5
)

func Solve(nodes []reader.Node, algorithm string, distance_matrix [][]int) ([][]int, error) {
	var (
		order           [][]int = make([][]int, NumCycles)  // kolejność odwiedzania wierzchołków dla obydwu cykli
		nodes_cycle_one int                                 // liczba wierzchołków w cyklu 1
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
	case "snr": // steepest node random - losowe rozwiązenie początowe, wersja stroma lokalnego przeszukiwania z wymianą wierzchołków w cyklu
		f = SteepestNodeRandom
	case "ser": // steepest edge random - losowe rozwiązenie początowe, wersja stroma lokalnego przeszukiwania z wymianą krawędzi w cyklu
		f = SteepestEdgeRandom
	case "gnr": // greedy node random - losowe rozwiązenie początowe, wersja stroma lokalnego przeszukiwania z wymianą wierzchołków w cyklu
		f = GreedyNodeRandom
	case "ger": // greedy edge random - losowe rozwiązenie początowe, wersja stroma lokalnego przeszukiwania z wymianą krawędzi w cyklu
		f = GreedyEdgeRandom
	case "sngc": // steepest node greedy cycle - wynik greedy cycle to rozwiązanie początkowe, wersja stroma lokalnego przeszukiwania z wymianą wierzchołków w cyklu
		f = SteepestNodeGreedyCycle
	case "segc": // steepest edge greedy cycle - wynik greedy cycle to rozwiązanie początkowe, wersja stroma lokalnego przeszukiwania z wymianą krawędzi w cyklu
		f = SteepestEdgeGreedyCycle
	case "gngc": // greedy node greedy cycle - wynik greedy cycle to rozwiązanie początkowe, wersja stroma lokalnego przeszukiwania z wymianą wierzchołków w cyklu
		f = GreedyNodeGreedyCycle
	case "gegc": // greedy edge greedy cycle - wynik greedy cycle to rozwiązanie początkowe, wersja stroma lokalnego przeszukiwania z wymianą krawędzi w cyklu
		f = GreedyEdgeGreedyCycle
	case "rwr": // random walk random - losowe rozwiązenie początkowe, wersja losowego błądzenia lokalnego przeszukiwania
		f = RandomWalkRandom
	case "rwgc": // random walk greedy cycle - wynik greedy cycle to rozwiązanie początkowe, wersja losowego błądzenia lokalnego przeszukiwania
		f = RandomWalkGreedyCycle
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

func Local_search(order [][]int, algorithm string, distance_matrix [][]int) ([][]int, error){
	panic("Dupa")
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





