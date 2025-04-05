package solver

import (
	"fmt"
	"math"
	"zad1/reader"
	"zad1/utils"
)

// ruch - zamiana wierzchołka N1 z N2 w cyklu Cycle
type MoveNode struct {
	Cycle int // numer cyklu
	N1    int // wierzchołek nr 1 - nr w cyklu
	N2    int // wierzchołek nr 2 - nr w cyklu
	Delta int // zmiana długości cyklu po dodaniu krawędzi
}

// ruch - zamiana wierzchołków między cyklami
type SwapMove struct {
	N1    int // wierzchołek z cyklu 1
	N2    int // wierzchołek z cyklu 2
	Delta int // zmiana długości cyklów po dodaniu krawędzi
}

type MoveEdge struct {
	Cycle int
}

type Move interface {
	ExecuteMove(distance_matrix [][]int, order [][]int) // wykonanie ruchu
	GetDelta() int                                      // zmiana długości cyklu po dodaniu krawędzi
}

func (m *SwapMove) ExecuteMove(distance_matrix [][]int, order [][]int) {
	order[0][m.N1], order[1][m.N2] = order[1][m.N2], order[0][m.N1] // zamiana wierzchołków między cyklami
}

func (m *SwapMove) GetDelta() int {
	return m.Delta // zmiana długości cyklu po dodaniu krawędzi
}

func (m *MoveNode) ExecuteMove(distance_matrix [][]int, order [][]int) {
	order[m.Cycle][m.N1], order[m.Cycle][m.N2] = order[m.Cycle][m.N2], order[m.Cycle][m.N1] // zamiana wierzchołków między cyklami
}

func (m *MoveNode) GetDelta() int {
	return m.Delta // zmiana długości cyklu po dodaniu krawędzi
}

func FindBestMove(moves []Move) (Move, int) {
	min_delta := math.MaxInt // minimalna zmiana długości cyklu
	var best_move Move = nil // najlepszy ruch
	for m := range moves {   // dla każdego ruchu
		move := moves[m]
		if move.GetDelta() < min_delta { // jeśli zmiana długości cyklu jest mniejsza od aktualnej
			best_move = move            // zapisz ruch jako najlepszy
			min_delta = move.GetDelta() // zapisz zmianę długości cyklu jako minimalną
		}
	}

	return best_move, min_delta // zwróć najlepszy ruch i minimalną zmianę długości cyklu
}

func DistancesBefore(distance_matrix [][]int, order [][]int) [][]int {
	var distances_before [][]int = make([][]int, NumCycles) // suma dystansów do wierzchołków przed i po aktualnym w cyklu

	for i := range distances_before { // dla każdego cyklu
		distances_before[i] = make([]int, len(order[i]))

		for j := range distances_before[i] { // dla każdego wierzchołka w cyklu
			curr_node := order[i][j]
			// dystans do wierzchołka przed i po aktualnym
			bj := utils.ElemBefore(order[i], j) // wierzchołek przed j
			aj := utils.ElemAfter(order[i], j)  // wierzchołek przed j
			distances_before[i][j] = distance_matrix[bj][curr_node] + distance_matrix[curr_node][aj]
		}
	}

	return distances_before
}

func AllMovesBetweenCycles(distance_matrix [][]int, order [][]int, distances_before [][]int) ([]SwapMove, error) {
	var (
		moves []SwapMove // aktualnie dostępne ruchy
		// distances_before [][]int    = make([][]int, NumCycles) // suma dystansów do wierzchołków przed i po aktualnym w cyklu
	)
	// for i := range distances_before { // dla każdego cyklu
	// 	distances_before[i] = make([]int, len(order[i]))

	// 	for j := range distances_before[i] { // dla każdego wierzchołka w cyklu
	// 		curr_node := order[i][j]
	// 		// dystans do wierzchołka przed i po aktualnym
	// 		bj := utils.ElemBefore(order[i], j) // wierzchołek przed j
	// 		aj := utils.ElemAfter(order[i], j)  // wierzchołek przed j
	// 		distances_before[i][j] = distance_matrix[bj][curr_node] + distance_matrix[curr_node][aj]
	// 	}
	// }

	for i := 0; i < len(order[0]); i++ {
		for j := 0; j < len(order[1]); j++ {
			// zamiana wierzchołka i z cyklu 1 z j z cyklu 2
			curr_node1 := order[0][i]
			curr_node2 := order[1][j]
			bi := utils.ElemBefore(order[0], i) // wierzchołek przed i w cyklu 1
			bj := utils.ElemBefore(order[1], j) // wierzchołek przed j w cyklu 2
			ai := utils.ElemAfter(order[0], i)  // wierzchołek po i w cyklu 1
			aj := utils.ElemAfter(order[1], j)  // wierzchołek po j w cyklu 2

			delta := distance_matrix[bi][curr_node2] + distance_matrix[curr_node2][ai] + // dystansy od wierzchołków przed i po aktualnych po zamianie
				distance_matrix[bj][curr_node1] + distance_matrix[curr_node1][aj] -
				distances_before[0][i] - distances_before[1][j]

			// // debug czy działa distance_matrix
			// if delta != distance_matrix[bi][curr_node2]+distance_matrix[curr_node2][ai]+ // dystansy od wierzchołków przed i po aktualnych po zamianie
			// 	distance_matrix[bj][curr_node1]+distance_matrix[curr_node1][aj]-
			// 	distance_matrix[bi][curr_node1]-distance_matrix[curr_node1][ai]- // dystansy od wierzchołków przed i po aktualnych przed zamianą
			// 	distance_matrix[bj][curr_node2]-distance_matrix[curr_node2][aj] {
			// 	panic("Invalid delta")
			// }

			// dodaj ruch do listy
			moves = append(moves, SwapMove{
				N1:    i,
				N2:    j,
				Delta: delta,
			})
		}
	}

	return moves, nil
}

func AllMovesNodesCycle(distance_matrix [][]int, order []int, cycle int, distances_before []int) []MoveNode {
	var (
		n1         int        // wierzchołek 1
		n2         int        // wierzchołek 2
		delta      int        // zmiana długości cyklu po dodaniu krawędzi
		moves_node []MoveNode // aktualnie dostępne ruchy
	)

	// dla każdej pary wierzchołków w cyklu; kolejność nie ma znaczenia
	for i := 0; i < len(order); i++ {
		for j := i + 1; j < len(order); j++ {
			n1 = order[i]                    // wierzchołek 1
			n2 = order[j]                    // wierzchołek 2
			bi := utils.ElemBefore(order, i) // wierzchołek przed i w cyklu
			bj := utils.ElemBefore(order, j) // wierzchołek przed j w cyklu
			ai := utils.ElemAfter(order, i)  // wierzchołek po i w cyklu
			aj := utils.ElemAfter(order, j)  // wierzchołek po j w cyklu

			if bi == n2 { // jeśli wierzchołki są sąsiadami w cyklu (j przed i)
				delta = distance_matrix[n1][bj] + distance_matrix[n2][ai] - // dystansy od wierzchołków przed i po aktualnych po zamianie
					distance_matrix[n1][ai] - distance_matrix[n2][bj] // dystansy od wierzchołków przed i po aktualnych przed zamianą
			} else if ai == n2 { // jeśli wierzchołki są sąsiadami w cyklu (i przed j)
				delta = distance_matrix[n1][aj] + distance_matrix[n2][bi] - // dystansy od wierzchołków przed i po aktualnych po zamianie
					distance_matrix[n1][bi] - distance_matrix[n2][aj] // dystansy od wierzchołków przed i po aktualnych przed zamianą
			} else { // jeśli wierzchołki nie są sąsiadami w cyklu
				delta = distance_matrix[bi][n2] + distance_matrix[n2][ai] + // dystansy od wierzchołków przed i po aktualnych po zamianie
					distance_matrix[bj][n1] + distance_matrix[n1][aj] -
					distances_before[i] - distances_before[j] // dystansy od wierzchołków przed i po aktualnych przed zamianą
			}

			// dodaj ruch do listy
			moves_node = append(moves_node, MoveNode{
				Cycle: cycle,
				N1:    i,
				N2:    j,
				Delta: delta,
			})
		}
	}

	return moves_node
}

func SteepestNodeRandom(distance_matrix [][]int, order [][]int, nodes []reader.Node) error {
	Random(distance_matrix, order, nodes)
	var (
		best_move       Move   = nil                                                // najlepszy ruch w iteracji
		min_delta       int    = math.MaxInt                                        // minimalna zmiana długości cyklu
		current_length1 int    = utils.CalculateCycleLen(order[0], distance_matrix) // akutalna długość cyklu 1
		current_length2 int    = utils.CalculateCycleLen(order[1], distance_matrix) // akutalna długość cyklu 2
		current_length  int    = current_length1 + current_length2
		all_moves       []Move // aktualnie dostępne ruchy
	)

	for {
		// ruchy pomiędzy cyklami
		distances_before := DistancesBefore(distance_matrix, order)                        // dystans do wierzchołków przed i po aktualnym w cyklu
		swap_moves, err := AllMovesBetweenCycles(distance_matrix, order, distances_before) // wszystkie ruchy między cyklami
		if err != nil {
			return err
		}
		all_moves = make([]Move, len(swap_moves)) // lista ruchów
		for i := range swap_moves {               // dla każdego ruchu
			all_moves[i] = &swap_moves[i] // dodaj ruch do listy
		}

		// best_swap_move, min_delta_swap := FindBestMove(all_moves) // najlepszy ruch i minimalna zmiana długości cyklu
		// best_move, min_delta = best_swap_move, min_delta_swap     // minimalna zmiana długości cyklu

		// ruchy w obrębie cyklu - zamiana wierzchołków w cyklu
		for c := 0; c < NumCycles; c++ { // dla każdego cyklu
			moves_cycle := AllMovesNodesCycle(distance_matrix, order[c], c, distances_before[c]) // wszystkie ruchy w cyklu zamiany wierzchołków

			for m := range moves_cycle { // dla każdego ruchu
				all_moves = append(all_moves, &moves_cycle[m]) // dodaj ruch do listy
			}
		}
		best_move, min_delta = FindBestMove(all_moves) // najlepszy ruch i minimalna zmiana długości cyklu

		// koniec iteracji
		if min_delta >= 0 { // jeśli nie znaleziono ruchu, który zmniejsza długość cyklu skończ przeszukiwanie
			best_move = nil // ustaw najlepszy ruch na nil
			break
		}
		// jeśli znaleziono ruch, to wykonaj go
		best_move.ExecuteMove(distance_matrix, order) // wykonaj najlepszy ruch
		current_length = current_length + min_delta   // aktualizuj długość cyklu
		best_move, min_delta = nil, math.MaxInt       // ustaw najlepszy ruch na nil i delta MaxInt
	}

	// debug
	current_length1 = utils.CalculateCycleLen(order[0], distance_matrix) // aktualizuj długość cyklu 1
	current_length2 = utils.CalculateCycleLen(order[1], distance_matrix) // aktualizuj długość cyklu 2
	if current_length1+current_length2 != current_length {
		fmt.Println("Długość cyklu:", current_length)
		fmt.Println("Długość cyklu 1 + 2:", current_length1+current_length2)
		panic("Invalid cycle length")
	}
	fmt.Println("SteepestNodeRandom: current length:", current_length)

	return nil
}

func SteepestEdgeRandom(distance_matrix [][]int, order [][]int, nodes []reader.Node) error {
	Random(distance_matrix, order, nodes)

	return nil
}

func GreedyNodeRandom(distance_matrix [][]int, order [][]int, nodes []reader.Node) error {
	Random(distance_matrix, order, nodes)

	return nil
}

func GreedyEdgeRandom(distance_matrix [][]int, order [][]int, nodes []reader.Node) error {
	Random(distance_matrix, order, nodes)

	return nil
}

func SteepestNodeGreedyCycle(distance_matrix [][]int, order [][]int, nodes []reader.Node) error {
	GreedyCycle(distance_matrix, order, nodes)

	return nil
}

func SteepestEdgeGreedyCycle(distance_matrix [][]int, order [][]int, nodes []reader.Node) error {
	GreedyCycle(distance_matrix, order, nodes)

	return nil
}

func GreedyNodeGreedyCycle(distance_matrix [][]int, order [][]int, nodes []reader.Node) error {
	GreedyCycle(distance_matrix, order, nodes)

	return nil
}

func GreedyEdgeGreedyCycle(distance_matrix [][]int, order [][]int, nodes []reader.Node) error {
	GreedyCycle(distance_matrix, order, nodes)

	return nil
}

func RandomWalkRandom(distance_matrix [][]int, order [][]int, nodes []reader.Node) error {
	Random(distance_matrix, order, nodes)

	return nil
}

func RandomWalkGreedyCycle(distance_matrix [][]int, order [][]int, nodes []reader.Node) error {
	GreedyCycle(distance_matrix, order, nodes)

	return nil
}
