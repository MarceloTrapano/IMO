package solver

import (
    "IMO/reader"
	"IMO/utils"
)
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
					temp_cycle = utils.Insert(cycle1, j, i) 
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
			temp_cycle := utils.Insert(cycle, j, i) 
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