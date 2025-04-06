package solver

import (
	"IMO/reader"
	"IMO/utils"
	"math"
	"math/rand"
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

// ruch - zamiana krawędzi po wierzchołku N1 z krawędzią po wierzchołku N2 w cyklu Cycle
type MoveEdge struct {
	Cycle int
	N1    int // wierzchołek nr 1 - nr w cyklu
	N2    int // wierzchołek nr 2 - nr w cyklu
	Delta int // zmiana długości cyklu po zamianie krawędzi
}

type Move interface {
	ExecuteMove(order [][]int) // wykonanie ruchu
	GetDelta() int             // zmiana długości cyklu po dodaniu krawędzi
	SetDelta(delta int)        // ustawienie zmiany długości cyklu po dodaniu krawędzi
}

func (m *SwapMove) ExecuteMove(order [][]int) {
	order[0][m.N1], order[1][m.N2] = order[1][m.N2], order[0][m.N1] // zamiana wierzchołków między cyklami
}

func (m *SwapMove) GetDelta() int {
	return m.Delta // zmiana długości cyklu po dodaniu krawędzi
}

func (m *SwapMove) SetDelta(delta int) {
	m.Delta = delta
}

func (m *MoveNode) ExecuteMove(order [][]int) {
	order[m.Cycle][m.N1], order[m.Cycle][m.N2] = order[m.Cycle][m.N2], order[m.Cycle][m.N1] // zamiana wierzchołków wewnątrz cyklu
}

func (m *MoveNode) GetDelta() int {
	return m.Delta // zmiana długości cyklu po dodaniu krawędzi
}

func (m *MoveNode) SetDelta(delta int) {
	m.Delta = delta
}

func (m *MoveEdge) ExecuteMove(order [][]int) {
	for i, j := m.N1+1, m.N2; i < j; i, j = i+1, j-1 {
		order[m.Cycle][i], order[m.Cycle][j] = order[m.Cycle][j], order[m.Cycle][i] // zamiana krawędzi wewnątrz cyklu
	}
}

func (m *MoveEdge) GetDelta() int {
	return m.Delta // zmiana długości cyklu po dodaniu krawędzi
}

func (m *MoveEdge) SetDelta(delta int) {
	m.Delta = delta
}

func CalculateDelta(move Move, distance_matrix [][]int, order [][]int) int {
	var (
		delta      int = 0 // zmiana długości cyklu po dodaniu krawędzi
		n1         int     // wierzchołek 1 - nr w cyklu
		n2         int     // wierzchołek 2 - nr w cyklu
		curr_node1 int     // numer wierzchołka 1
		curr_node2 int     // numer wierzchołka 2
		bi         int     // wierzchołek przed i w cyklu 1
		bj         int     // wierzchołek przed j w cyklu 2
		ai         int     // wierzchołek po i w cyklu 1
		aj         int     // wierzchołek po j w cyklu 2
	)
	switch m := move.(type) {
	case *MoveNode:
		n1, n2 = m.N1, m.N2                       // wierzchołki 1 i 2 - nr w cyklu
		curr_node1 = order[m.Cycle][m.N1]         // wierzchołek aktualny w cyklu
		curr_node2 = order[m.Cycle][m.N2]         // wierzchołek aktualny w cyklu
		bi = utils.ElemBefore(order[m.Cycle], n1) // wierzchołek przed i w cyklu 1
		bj = utils.ElemBefore(order[m.Cycle], n2) // wierzchołek przed j w cyklu 2
		ai = utils.ElemAfter(order[m.Cycle], n1)  // wierzchołek po i w cyklu 1
		aj = utils.ElemAfter(order[m.Cycle], n2)  // wierzchołek po j w cyklu 2
	case *SwapMove:
		n1, n2 = m.N1, m.N2                 // wierzchołki 1 i 2 - nr w cyklu
		curr_node1 = order[0][m.N1]         // wierzchołek aktualny w cyklu 1
		curr_node2 = order[1][m.N2]         // wierzchołek aktualny w cyklu 2
		bi = utils.ElemBefore(order[0], n1) // wierzchołek przed i w cyklu 1
		bj = utils.ElemBefore(order[1], n2) // wierzchołek przed j w cyklu 2
		ai = utils.ElemAfter(order[0], n1)  // wierzchołek po i w cyklu 1
		aj = utils.ElemAfter(order[1], n2)  // wierzchołek po j w cyklu 2
	case *MoveEdge:
		n1, n2 = m.N1, m.N2 // wierzchołki 1 i 2 - nr w cyklu
		curr_node1 = order[m.Cycle][m.N1]
		curr_node2 = order[m.Cycle][m.N2]
		ai = utils.ElemAfter(order[m.Cycle], n1)
		aj = utils.ElemAfter(order[m.Cycle], n2)
	}

	switch m := move.(type) {
	case *SwapMove:
		delta = distance_matrix[bi][curr_node2] + distance_matrix[curr_node2][ai] + // dystansy od wierzchołków przed i po aktualnych po zamianie
			distance_matrix[bj][curr_node1] + distance_matrix[curr_node1][aj] -
			distance_matrix[bi][curr_node1] - distance_matrix[curr_node1][ai] - // dystansy od wierzchołków przed i po aktualnych przed zamianą
			distance_matrix[bj][curr_node2] - distance_matrix[curr_node2][aj] // dystansy od wierzchołków przed i po aktualnych przed zamianą
		m.Delta = delta // ustaw zmianę długości cyklu na mniejszą
	case *MoveNode:
		if bi == curr_node2 { // jeśli wierzchołki są sąsiadami w cyklu (j przed i)
			delta = distance_matrix[curr_node1][bj] + distance_matrix[curr_node2][ai] - // dystansy od wierzchołków przed i po aktualnych po zamianie
				distance_matrix[curr_node1][ai] - distance_matrix[curr_node2][bj] // dystansy od wierzchołków przed i po aktualnych przed zamianą
		} else if ai == curr_node2 { // jeśli wierzchołki są sąsiadami w cyklu (i przed j)
			delta = distance_matrix[curr_node1][aj] + distance_matrix[curr_node2][bi] - // dystansy od wierzchołków przed i po aktualnych po zamianie
				distance_matrix[curr_node1][bi] - distance_matrix[curr_node2][aj] // dystansy od wierzchołków przed i po aktualnych przed zamianą
		} else { // jeśli wierzchołki nie są sąsiadami w cyklu - tak jak w SwapMove
			delta = distance_matrix[bi][curr_node2] + distance_matrix[curr_node2][ai] + // dystansy od wierzchołków przed i po aktualnych po zamianie
				distance_matrix[bj][curr_node1] + distance_matrix[curr_node1][aj] -
				distance_matrix[bi][curr_node1] - distance_matrix[curr_node1][ai] - // dystansy od wierzchołków przed i po aktualnych przed zamianą
				distance_matrix[bj][curr_node2] - distance_matrix[curr_node2][aj] // dystansy od wierzchołków przed i po aktualnych przed zamianą
		}
		m.Delta = delta // ustaw zmianę długości cyklu na mniejszą
	case *MoveEdge:
		delta = distance_matrix[curr_node1][curr_node2] + distance_matrix[ai][aj] - // dystansy po zamianie krawędzi
			distance_matrix[ai][curr_node1] + distance_matrix[aj][curr_node2] // dystansy przed zamianą krawędzi
		m.Delta = delta // ustaw zmianę długości cyklu na mniejszą
	}
	return delta
}

func FisherYatesShuffle[T comparable](arr []T) []T {
	for i := len(arr) - 1; i > 0; i-- { // iteracja po arr od końca
		j := rand.Intn(i + 1)           // losowy indeks od 0 do i
		arr[i], arr[j] = arr[j], arr[i] // zamień elementy miejscami
	}
	return arr // zwróć przetasowaną tablicę
}

func FindBestMoveGreedy(moves []Move, distance_matrix [][]int, order [][]int) (Move, int) {
	moves = FisherYatesShuffle(moves) // przetasuj ruchy
	for m := range moves {            // dla każdego ruchu
		move := moves[m]
		delta := CalculateDelta(move, distance_matrix, order)
		if delta < 0 { // jeśli zmiana długości cyklu jest mniejsza od aktualnej i mniejsza od 0
			move.SetDelta(delta) // ustaw zmianę długości cyklu na mniejszą
			return move, delta
		}
	}
	return nil, math.MaxInt // zwróć najlepszy ruch i minimalną zmianę długości cyklu
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
	)

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

func AllMovesBetweenCyclesNoDistance(order [][]int) ([]SwapMove, error) {
	var (
		moves []SwapMove // aktualnie dostępne ruchy
	)

	for i := 0; i < len(order[0]); i++ {
		for j := 0; j < len(order[1]); j++ {
			// zamiana wierzchołka i z cyklu 1 z j z cyklu 2
			moves = append(moves, SwapMove{
				N1:    i,
				N2:    j,
				Delta: math.MaxInt,
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

func AllMovesNodesCycleNoDistance(order []int, cycle int) []MoveNode {
	var (
		moves_node []MoveNode // aktualnie dostępne ruchy
	)

	// dla każdej pary wierzchołków w cyklu; kolejność nie ma znaczenia
	for i := 0; i < len(order); i++ {
		for j := i + 1; j < len(order); j++ {

			// dodaj ruch do listy
			moves_node = append(moves_node, MoveNode{
				Cycle: cycle,
				N1:    i,
				N2:    j,
				Delta: math.MaxInt,
			})
		}
	}

	return moves_node
}

func AllMovesEdgesCycle(distance_matrix [][]int, order []int, cycle int, distances_before []int) []MoveEdge {
	panic("not implemented") // TODO: implement
	return nil
}

func AllMovesEdgesCycleNoDistance(order []int, cycle int) []MoveEdge {
	panic("not implemented") // TODO: implement
	return nil
}

func SteepestNode(distance_matrix [][]int, order [][]int) error {
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
		best_move.ExecuteMove(order)                // wykonaj najlepszy ruch
		current_length = current_length + min_delta // aktualizuj długość cyklu
		best_move, min_delta = nil, math.MaxInt     // ustaw najlepszy ruch na nil i delta MaxInt
	}

	return nil
}

func GreedyNode(distance_matrix [][]int, order [][]int) error {
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
		swap_moves, err := AllMovesBetweenCyclesNoDistance(order) // wszystkie ruchy między cyklami
		if err != nil {
			return err
		}
		all_moves = make([]Move, len(swap_moves)) // lista ruchów
		for i := range swap_moves {               // dla każdego ruchu
			all_moves[i] = &swap_moves[i] // dodaj ruch do listy
		}

		// ruchy w obrębie cyklu - zamiana wierzchołków w cyklu
		for c := 0; c < NumCycles; c++ { // dla każdego cyklu
			moves_cycle := AllMovesNodesCycleNoDistance(order[c], c) // wszystkie ruchy w cyklu zamiany wierzchołków

			for m := range moves_cycle { // dla każdego ruchu
				all_moves = append(all_moves, &moves_cycle[m]) // dodaj ruch do listy
			}
		}
		best_move, min_delta = FindBestMoveGreedy(all_moves, distance_matrix, order) // najlepszy ruch i minimalna zmiana długości cyklu

		// koniec iteracji
		if min_delta >= 0 { // jeśli nie znaleziono ruchu, który zmniejsza długość cyklu skończ przeszukiwanie
			best_move = nil // ustaw najlepszy ruch na nil
			break
		}
		// jeśli znaleziono ruch, to wykonaj go
		best_move.ExecuteMove(order)                // wykonaj najlepszy ruch
		current_length = current_length + min_delta // aktualizuj długość cyklu

		best_move, min_delta = nil, math.MaxInt // ustaw najlepszy ruch na nil i delta MaxInt
	}

	return nil
}

// TODO: SteepestEdge: implementacja AllMovesEdgesCycle i AllMovesEdgesCycleNoDistance -> to co SteepestNode ale AllMovesNodesCycle zamienić na AllMovesEdgesCycle
// TODO: GreedyEdge: implementacja AllMovesEdgesCycle i AllMovesEdgesCycleNoDistance -> to co GreedyNode ale AllMovesNodesCycle zamienić na AllMovesEdgesCycle

func SteepestEdge(distance_matrix [][]int, order [][]int) error {
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

		// ruchy w obrębie cyklu - zamiana wierzchołków w cyklu
		for c := 0; c < NumCycles; c++ { // dla każdego cyklu
			moves_cycle := AllMovesEdgesCycle(distance_matrix, order[c], c, distances_before[c]) // wszystkie ruchy w cyklu zamiany wierzchołków

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
		best_move.ExecuteMove(order)                // wykonaj najlepszy ruch
		current_length = current_length + min_delta // aktualizuj długość cyklu
		best_move, min_delta = nil, math.MaxInt     // ustaw najlepszy ruch na nil i delta MaxInt
	}

	return nil
}

func GreedyEdge(distance_matrix [][]int, order [][]int) error {
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
		swap_moves, err := AllMovesBetweenCyclesNoDistance(order) // wszystkie ruchy między cyklami
		if err != nil {
			return err
		}
		all_moves = make([]Move, len(swap_moves)) // lista ruchów
		for i := range swap_moves {               // dla każdego ruchu
			all_moves[i] = &swap_moves[i] // dodaj ruch do listy
		}

		// ruchy w obrębie cyklu - zamiana wierzchołków w cyklu
		for c := 0; c < NumCycles; c++ { // dla każdego cyklu
			moves_cycle := AllMovesEdgesCycleNoDistance(order[c], c) // wszystkie ruchy w cyklu zamiany wierzchołków

			for m := range moves_cycle { // dla każdego ruchu
				all_moves = append(all_moves, &moves_cycle[m]) // dodaj ruch do listy
			}
		}
		best_move, min_delta = FindBestMoveGreedy(all_moves, distance_matrix, order) // najlepszy ruch i minimalna zmiana długości cyklu

		// koniec iteracji
		if min_delta >= 0 { // jeśli nie znaleziono ruchu, który zmniejsza długość cyklu skończ przeszukiwanie
			best_move = nil // ustaw najlepszy ruch na nil
			break
		}
		// jeśli znaleziono ruch, to wykonaj go
		best_move.ExecuteMove(order)                // wykonaj najlepszy ruch
		current_length = current_length + min_delta // aktualizuj długość cyklu

		best_move, min_delta = nil, math.MaxInt // ustaw najlepszy ruch na nil i delta MaxInt
	}

	return nil
}

func SteepestNodeRandom(distance_matrix [][]int, order [][]int, nodes []reader.Node) error {
	Random(distance_matrix, order, nodes)
	SteepestNode(distance_matrix, order)

	return nil
}

func SteepestEdgeRandom(distance_matrix [][]int, order [][]int, nodes []reader.Node) error {
	Random(distance_matrix, order, nodes)
	SteepestEdge(distance_matrix, order)

	return nil
}

func GreedyNodeRandom(distance_matrix [][]int, order [][]int, nodes []reader.Node) error {
	Random(distance_matrix, order, nodes)
	GreedyNode(distance_matrix, order)

	return nil
}

func GreedyEdgeRandom(distance_matrix [][]int, order [][]int, nodes []reader.Node) error {
	Random(distance_matrix, order, nodes)
	GreedyEdge(distance_matrix, order)

	return nil
}

func SteepestNodeGreedyCycle(distance_matrix [][]int, order [][]int, nodes []reader.Node) error {
	GreedyCycle(distance_matrix, order, nodes)
	SteepestNode(distance_matrix, order)

	return nil
}

func SteepestEdgeGreedyCycle(distance_matrix [][]int, order [][]int, nodes []reader.Node) error {
	GreedyCycle(distance_matrix, order, nodes)
	SteepestEdge(distance_matrix, order)

	return nil
}

func GreedyNodeGreedyCycle(distance_matrix [][]int, order [][]int, nodes []reader.Node) error {
	GreedyCycle(distance_matrix, order, nodes)
	GreedyNode(distance_matrix, order)

	return nil
}

func GreedyEdgeGreedyCycle(distance_matrix [][]int, order [][]int, nodes []reader.Node) error {
	GreedyCycle(distance_matrix, order, nodes)
	GreedyEdge(distance_matrix, order)

	return nil
}

func RandomWalkRandom(distance_matrix [][]int, order [][]int, nodes []reader.Node) error {
	Random(distance_matrix, order, nodes)
	panic("not implemented") // TODO: implement
	// TODO: RandomWalk

	return nil
}

func RandomWalkGreedyCycle(distance_matrix [][]int, order [][]int, nodes []reader.Node) error {
	GreedyCycle(distance_matrix, order, nodes)
	panic("not implemented") // TODO: implement
	// TODO: RandomWalk

	return nil
}
