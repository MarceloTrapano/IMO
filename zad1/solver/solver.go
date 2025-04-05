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
	start_node_1, start_node_2, _ := PickRandomNodes(nodes) // wybór startowych punktów

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
	start_node_1, start_node_2, _ := PickRandomClosestNodes(distance_matrix, nodes) // wybór startowych punktów

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
					cost = 0                                // występowały leaki pamięci i program miał pr
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
	start_node_1, start_node_2, _ := PickRandomNodes(nodes) // wybór startowych punktów

	var (
		visited []bool = make([]bool, len(nodes)) // tablica dodanych wierzchołków
		cycle1  []int  = make([]int, 0, len(order[0]))
		cycle2  []int  = make([]int, 0, len(order[1]))
	)
	visited[start_node_1] = true
	visited[start_node_2] = true

	cycle1 = append(cycle1, start_node_1)
	cycle2 = append(cycle2, start_node_2)

	node_val := 10000
	node_idx := -1
	for i, val := range distance_matrix[start_node_1] {
		if val == 0 || visited[i] {
			continue
		}
		if val < node_val {
			node_val = val
			node_idx = i
		}
	}
	cycle1 = append(cycle1, node_idx)
	visited[node_idx] = true

	node_val = 10000
	node_idx = -1
	for i, val := range distance_matrix[start_node_2] {
		if val == 0 || visited[i] {
			continue
		}
		if val < node_val {
			node_val = val
			node_idx = i
		}
	}
	cycle2 = append(cycle2, node_idx)
	visited[node_idx] = true

	for len(cycle1) < len(order[0]) || len(cycle2) < len(order[1]) {
		node1, node2, _ := BestNodes(cycle1, distance_matrix, visited)
		best_score1, second_best_score1, idx1, _ := Calculate4Regret(node1, cycle1, distance_matrix)
		regret1 := second_best_score1 - best_score1
		best_score2, second_best_score2, idx2, _ := Calculate4Regret(node2, cycle1, distance_matrix)
		regret2 := second_best_score2 - best_score2
		if regret1 > regret2 {
			cycle1 = utils.Insert(cycle1, idx1, node1)
			visited[node1] = true
		} else {
			cycle1 = utils.Insert(cycle1, idx2, node2)
			visited[node2] = true
		}

		node3, node4, _ := BestNodes(cycle2, distance_matrix, visited)

		best_score3, second_best_score3, idx3, _ := Calculate4Regret(node3, cycle2, distance_matrix)
		regret3 := second_best_score3 - best_score3
		best_score4, second_best_score4, idx4, _ := Calculate4Regret(node4, cycle2, distance_matrix)
		regret4 := second_best_score4 - best_score4
		if regret3 > regret4 {
			cycle2 = utils.Insert(cycle2, idx3, node3)
			visited[node3] = true
		} else {
			cycle2 = utils.Insert(cycle2, idx4, node4)
			visited[node4] = true
		}

	}
	order[0] = cycle1
	order[1] = cycle2
	return nil
}
func WeightedRegret(distance_matrix [][]int, order [][]int, nodes []reader.Node) error {
	start_node_1, start_node_2, _ := PickRandomClosestNodes(distance_matrix, nodes) // wybór startowych punktów

	var (
		visited       []bool = make([]bool, len(nodes)) // tablica dodanych wierzchołków
		cycle1        []int  = make([]int, 0, len(order[0]))
		cycle2        []int  = make([]int, 0, len(order[1]))
		weight_regret int    = 1
		weight_change int    = -4
	)
	visited[start_node_1] = true
	visited[start_node_2] = true

	cycle1 = append(cycle1, start_node_1)
	cycle2 = append(cycle2, start_node_2)

	node_val := 10000
	node_idx := -1
	for i, val := range distance_matrix[start_node_1] {
		if val == 0 || visited[i] {
			continue
		}
		if val < node_val {
			node_val = val
			node_idx = i
		}
	}
	cycle1 = append(cycle1, node_idx)
	visited[node_idx] = true

	node_val = 10000
	node_idx = -1
	for i, val := range distance_matrix[start_node_2] {
		if val == 0 || visited[i] {
			continue
		}
		if val < node_val {
			node_val = val
			node_idx = i
		}
	}
	cycle2 = append(cycle2, node_idx)
	visited[node_idx] = true

	for len(cycle1) < len(order[0]) || len(cycle2) < len(order[1]) {
		node1, node2, _ := BestNodes(cycle1, distance_matrix, visited)

		best_score1, second_best_score1, idx1, _ := Calculate4Regret(node1, cycle1, distance_matrix)
		regret1 := second_best_score1 - best_score1
		total_cost1 := regret1*weight_regret + best_score1*weight_change
		best_score2, second_best_score2, idx2, _ := Calculate4Regret(node2, cycle1, distance_matrix)
		regret2 := second_best_score2 - best_score2
		total_cost2 := regret2*weight_regret + best_score2*weight_change
		if total_cost1 > total_cost2 {
			cycle1 = utils.Insert(cycle1, idx1, node1)
			visited[node1] = true
		} else {
			cycle1 = utils.Insert(cycle1, idx2, node2)
			visited[node2] = true
		}

		node3, node4, _ := BestNodes(cycle2, distance_matrix, visited)

		best_score3, second_best_score3, idx3, _ := Calculate4Regret(node3, cycle2, distance_matrix)
		regret3 := second_best_score3 - best_score3
		total_cost3 := regret3*weight_regret + best_score3*weight_change
		best_score4, second_best_score4, idx4, _ := Calculate4Regret(node4, cycle2, distance_matrix)
		regret4 := second_best_score4 - best_score4
		total_cost4 := regret4*weight_regret + best_score4*weight_change
		if total_cost3 > total_cost4 {
			cycle2 = utils.Insert(cycle2, idx3, node3)
			visited[node3] = true
		} else {
			cycle2 = utils.Insert(cycle2, idx4, node4)
			visited[node4] = true
		}

	}
	order[0] = cycle1
	order[1] = cycle2
	return nil
}

func Calculate4Regret(node1 int, cycle []int, distance_matrix [][]int) (int, int, int, error) {
	minimal_cost := -1
	second_minimal_cost := -1
	idx := -1
	for i := range cycle {
		temp_cycle := utils.Insert(cycle, i, node1)
		cost := utils.CalculateCycleLen(temp_cycle, distance_matrix)
		if minimal_cost == -1 || cost < minimal_cost {
			minimal_cost = cost
			idx = i
		} else if second_minimal_cost == -1 || cost < second_minimal_cost {
			second_minimal_cost = cost
		}
	}
	return minimal_cost, second_minimal_cost, idx, nil
}
func BestNodes(cycle []int, distance_matrix [][]int, visited []bool) (int, int, error) {
	var (
		node1 int
		node2 int
	)
	worst_cost := -1
	second_worst_cost := -1
	for i := range visited {
		if visited[i] {
			continue
		}
		for j := range cycle {
			temp_cycle := utils.Insert(cycle, j, i) // musiałem sam napisać funkcję do dodwawania elementu do macierzy XD
			cost := utils.CalculateCycleLen(temp_cycle, distance_matrix)
			if cost > second_worst_cost {
				second_worst_cost = cost
				node2 = i
			} else if cost > worst_cost {
				worst_cost = cost
				node1 = i
			}
		}
	}
	return node1, node2, nil
}

/*
func Regret(distance_matrix [][]int, order [][]int, nodes []reader.Node) error {
	var (
		visited   []bool                    = make([]bool, len(nodes))                   // tablica dodanych wierzchołków
		cycles    []*utils.Edge             = make([]*utils.Edge, NumCycles)             // cykle krawędzi
		distances [][]*utils.EdgeLinkedList = make([][]*utils.EdgeLinkedList, NumCycles) // tablica linked list z dystansami do krawędzi dla każdego cyklu i każdego wierzchołka
		edges     []*utils.Edge                                                          // tablica wszystkich krawędzi
		lenCycles []int                     = make([]int, NumCycles)                     // długości cykli w wierzchołkach
	)
	// przygotowanie tablicy dystansów
	for i := range distances {
		distances[i] = make([]*utils.EdgeLinkedList, len(nodes))
	}
	start_node_1, start_node_2, _ := PickRandomNodes(nodes) // wybór startowych punktów

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
	cycles[0] = edge1
	cycles[1] = edge2
	lenCycles[0] = 2
	lenCycles[1] = 2

	// aktualizacja początkowych dystansów
	for i := 0; i < len(nodes); i++ {
		if visited[i] {
			continue
		}
		distances[0][i] = &utils.EdgeLinkedList{Edge: 0, Next: nil, Value: utils.EdgeInsertValue(distance_matrix, i, edge1)}
		distances[1][i] = &utils.EdgeLinkedList{Edge: 1, Next: nil, Value: utils.EdgeInsertValue(distance_matrix, i, edge2)}
		// newEdges := []utils.EdgeLinkedList{{Edge: 1, Next: nil, Value: utils.EdgeInsertValue(distance_matrix, i, &edges[1])}}
		// distances[0][i] = utils.UpdateDistances(distances[0][i], distance_matrix, nil, newEdges, false)
	}
	edges = append(edges, edge1, edge2)

	// dodawanie wierzchołków - rozpatrujemy żal dla każdej krawędzi w cyklu i pozostałego wierzchołka dla każdego cyklu osobno
	// dopóki nie dodamy wszystkich wierzchołków
	for lenCycles[0] < len(order[0]) || lenCycles[1] < len(order[1]) { // krawędzi będzie o 2 mniej niż wierzchołków
		// dla każdego cyklu
		for i := 0; i < NumCycles; i++ {
			var (
				max_nodes      int = len(order[i])
				max_regret     int = math.MinInt64
				max_regret_idx int
			)
			// cykl jest ukończony
			if lenCycles[i] == max_nodes {
				continue
			}
			// obliczanie żalu
			regrets := ComputeRegrets(distances[i], visited, 2)
			for i, regret := range regrets {
				if !visited[i] && regret > max_regret {
					max_regret = regret
					max_regret_idx = i
				}
			}
			var (
				best_edge_idx int         = distances[i][max_regret_idx].Edge // posortowane rosnąco (UpdateDistances co iterację) także pierwsza jest najlepsza
				best_edge     *utils.Edge = edges[best_edge_idx]
				delEdges      []int
				newEdges      []utils.EdgeLinkedList
			)

			// tworzenie nowych krawędzi
			e1 := utils.NewEdge(best_edge.From, max_regret_idx, distance_matrix, best_edge.Prev, nil)
			e2 := utils.NewEdge(max_regret_idx, best_edge.To, distance_matrix, e1, best_edge.Next) // niepowtarzające się indeksy na poziomie cyklu
			e1.Next = e2

			// dodanie krawędzi
			edges = append(edges, e1, e2)
			// aktualizacja cyklu - edges i cycles mają wskaźniki na te same krawędzie
			if lenCycles[i] <= 2 { // pierwsze dodanie wierzchołka - dodajemy 2 krawędzie, stara zostaje
				best_edge.Prev = e2
				best_edge.Next = e1
				e1.Prev = best_edge
				e2.Next = best_edge
				e1.From = best_edge.To // to jest trochę nieintuicyjne, bez rozrysowania ciężko mi było i się męczyłem :|
				e2.To = best_edge.From
				e1.Length = distance_matrix[e1.From][e1.To]
				e2.Length = distance_matrix[e2.From][e2.To]
				best_edge_idx = -1 // by nie usuwać w tej iteracji krawędzi
			} else { // kolejne dodanie wierzchołka - dodajemy 2 krawędzie, stara odpada
				best_edge.Prev.Next = e1
				best_edge.Next.Prev = e2
				if best_edge == cycles[i] {
					cycles[i] = e1
				}
				best_edge = nil
			}

			visited[max_regret_idx] = true
			// aktualizacja dystansów
			for j := range distances[i] {
				if visited[j] {
					continue
				}
				newEdges = []utils.EdgeLinkedList{
					{Edge: len(edges), Next: nil, Value: utils.EdgeInsertValue(distance_matrix, j, e1)},
					{Edge: len(edges) + 1, Next: nil, Value: utils.EdgeInsertValue(distance_matrix, j, e2)},
				}

				delEdges = []int{best_edge_idx}
				distances[i][j] = utils.UpdateDistances(distances[i][j], distance_matrix, delEdges, newEdges, false)
			}

			// aktualizacja tablic
			edges = append(edges, e1, e2)
			// visited[max_regret_idx] = true
			lenCycles[i]++
		}
	}
	order[0] = utils.EdgeToNodeCycle(cycles[0])
	order[1] = utils.EdgeToNodeCycle(cycles[1])

	return nil
}

func ComputeRegrets(distances []*utils.EdgeLinkedList, visited []bool, degree int) []int {
	var (
		regrets []int = make([]int, len(distances))
		oneEdge bool  // czy tylko 1 krawędzi; dla każdego wierzchołka to samo bo każdy tak samo długie listy - l. krawędzi
	)
	for i := range visited {
		if !visited[i] {
			oneEdge = distances[i].Next == nil
			break
		}
	}
	for i := range distances {
		var current_reg int = 1
		if visited[i] {
			continue
		}
		var current *utils.EdgeLinkedList = distances[i]
		var first_value int = current.Value
		if oneEdge { // gdyby była jedna krawędź - nie ma sensu liczyć żalu zamiast tego po prostu wzrost długości
			regrets[i] = first_value
			continue
		}

		for current_reg < degree && current.Next != nil { // liczenie kolejnych żali
			regrets[i] += current.Next.Value - first_value // kolejne krawędzie mają większe wartości
			current_reg++
			current = current.Next
		}
	}
	return regrets
}
*/

func Random(distance_matrix [][]int, order [][]int, nodes []reader.Node) error {
	result := make([][]int, NumCycles)
	nodes_copy := make([]reader.Node, len(nodes))
	copy(nodes_copy, nodes)             // kopiowanie tablicy nodes do nowej tablicy
	nodes_nr := make([]int, len(nodes)) // tablica z numerami wierzchołków
	for i := range nodes {
		nodes_nr[i] = i
	}
	for len(nodes_copy) > 0 {
		for i := 0; i < NumCycles; i++ {
			if len(nodes_copy) == 0 {
				break
			}
			node_idx, err := PickRandomNode(nodes_copy)
			if err != nil {
				return err
			}
			result[i] = append(result[i], nodes_nr[node_idx]) // dodanie wierzchołka do cyklu
			// usunięcie wierzchołka z list
			nodes_copy = append(nodes_copy[:node_idx], nodes_copy[node_idx+1:]...)
			nodes_nr = append(nodes_nr[:node_idx], nodes_nr[node_idx+1:]...)
		}
	}
	copy(order, result)
	return nil
}
