package solver

import (
	"IMO/utils"
	"math"
	"sort"
)

type Applicability int

const (
	Applicable Applicability = iota
	NotApplicable
	MayBeApplicable // 1 krawędź w drugą stronę
)

var ApplicabilityType = map[Applicability]string{
	Applicable:      "Applicable",
	NotApplicable:   "NotApplicable",
	MayBeApplicable: "MayBeApplicable",
}

func (a Applicability) String() string {
	return ApplicabilityType[a]
}

// ruch - zamiana krawędzi N1, SN1 i N2, SN2 w cyklu Cycle
type MoveEdgeDetail struct {
	Cycle int
	N1    int // wierzchołek nr 1
	N2    int // wierzchołek nr 2
	SN1   int // następny wierzchołek po N1
	SN2   int // następny wierzchołek po N2
	Delta int // zmiana długości cyklu po zamianie krawędzi
}

// ruch - zamiana wierzchołków między cyklami
type SwapMoveDetail struct {
	N1    int // wierzchołek z cyklu 1
	N2    int // wierzchołek z cyklu 2
	SN1   int // następny wierzchołek po N1
	SN2   int // następny wierzchołek po N2
	PN1   int // poprzedni wierzchołek po N1
	PN2   int // poprzedni wierzchołek po N2
	Delta int // zmiana długości cyklów po dodaniu krawędzi
}

func (m *MoveEdgeDetail) ExecuteMove(order [][]int) {
	nodes := []int{m.N1, m.N2}                        // wierzchołki do zamiany
	indexes := utils.IndexesOf(order[m.Cycle], nodes) // znajdź indeksy w cyklu
	if len(indexes) != 2 {
		panic("MoveEdgeDetail: indexes not found")
	}

	n1_index, n2_index := indexes[0], indexes[1]
	if n1_index > n2_index { // zamień indeksy
		n1_index, n2_index = n2_index, n1_index
	}

	for i, j := n1_index+1, n2_index; i < j; i, j = i+1, j-1 { // jak w MoveEdge
		order[m.Cycle][i], order[m.Cycle][j] = order[m.Cycle][j], order[m.Cycle][i] // zamiana krawędzi wewnątrz cyklu
	}
}

func (m *MoveEdgeDetail) GetDelta() int {
	return m.Delta // zmiana długości cyklu po dodaniu krawędzi
}

func (m *MoveEdgeDetail) SetDelta(delta int) {
	m.Delta = delta
}

func (m *SwapMoveDetail) ExecuteMove(order [][]int) {
	indexes := utils.IndexesOf(order[0], []int{m.N1})                    // znajdź indeksy w cyklu
	indexes = append(indexes, utils.IndexesOf(order[1], []int{m.N2})...) // znajdź indeksy w cyklu

	if len(indexes) != 2 {
		panic("SwapMoveDetail: indexes not found")
	}
	n1_index, n2_index := indexes[0], indexes[1]

	order[0][n1_index], order[1][n2_index] = order[1][n2_index], order[0][n1_index] // zamiana wierzchołków między cyklami
}

func (m *SwapMoveDetail) GetDelta() int {
	return m.Delta // zmiana długości cyklu po dodaniu krawędzi
}

func (m *SwapMoveDetail) SetDelta(delta int) {
	m.Delta = delta
}

func FastLocalSearch(distance_matrix [][]int, order [][]int) error {
	// inicjacja tablicy z najlepszymi ruchami
	var best_moves []Move // aktualnie najlepsze ruchy posortowane od najlepszego do najgorszego

	distance_before := DistancesBefore(distance_matrix, order)
	swap_moves, err := BestMovesBetweenCycles(distance_matrix, order, distance_before) // wybieranie ruchów między cyklami poprawiających wynik
	if err != nil {
		return err
	}
	for m := range swap_moves {
		best_moves = append(best_moves, &swap_moves[m])
	}

	for c := 0; c < NumCycles; c++ {
		moves_cycle := BestMovesEdgesCycle(distance_matrix, order[c], c) // wybieranie ruchów krawędzi poprawiających wynik
		for m := range moves_cycle {
			best_moves = append(best_moves, &moves_cycle[m])
		}
	}
	sort.Slice(best_moves, func(i, j int) bool {
		return best_moves[i].GetDelta() < best_moves[j].GetDelta() // sortowanie po najwyższych deltach
	})

	for len(best_moves) > 0 {
		to_delete := []int{}  // indeksy do usunięcia
		new_moves := []Move{} // nowe ruchy do dodania

	Loop: // label loopa do breakowania
		for i, move := range best_moves {
			applicability := CheckApplicability(move, order)
			switch applicability {
			case Applicable:
				move.ExecuteMove(order)

				new_moves, err = FindNewMoves(distance_matrix, order, move) // znajdź nowe ruchy
				if err != nil {
					return err
				}

				to_delete = append(to_delete, i) // usuwamy po wykonaniu
				break Loop
			case NotApplicable:
				to_delete = append(to_delete, i) // usuwamy - ruch nie jest aplikowalny
			case MayBeApplicable: // nie usuwamy - ruch może być aplikowalny po czasie - krawędzie istnieją ale w drugą stronę
			}
			if i == len(best_moves)-1 { // jak nie znajdzie ruchu aplikowalnego to przerywa
				return nil
			}
		}
		best_moves = utils.RemoveIndexes(best_moves, to_delete) // usuń ruchy, które nie są aplikowalne
		best_moves = AddNewMoves(best_moves, new_moves)         // dodaj nowe ruchy
	}

	return nil
}
func BestMovesBetweenCycles(distance_matrix [][]int, order [][]int, distances_before [][]int) ([]SwapMoveDetail, error) {
	var (
		moves []SwapMoveDetail // aktualnie dostępne ruchy
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
			if delta < 0 {
				// dodaj ruch do listy
				moves = append(moves, SwapMoveDetail{
					N1:    curr_node1,
					N2:    curr_node2,
					SN1:   ai,
					SN2:   aj,
					PN1:   bi,
					PN2:   bj,
					Delta: delta,
				})
			}
		}
	}

	return moves, nil
}

func BestMovesEdgesCycle(distance_matrix [][]int, order []int, cycle int) []MoveEdgeDetail {
	var (
		n1         int              // wierzchołek 1
		n2         int              // wierzchołek 2
		delta      int              // zmiana długości cyklu po dodaniu krawędzi
		moves_node []MoveEdgeDetail // aktualnie dostępne ruchy
	)

	// dla każdej pary wierzchołków w cyklu; kolejność nie ma znaczenia
	for i := 0; i < len(order); i++ {
		for j := i + 1; j < len(order); j++ {
			n1, n2 = order[i], order[j] // wierzchołki 1 i 2 - nr w cyklu
			ai := utils.ElemAfter(order, i)
			aj := utils.ElemAfter(order, j)
			delta = distance_matrix[n1][n2] + distance_matrix[ai][aj] - // dystansy po zamianie krawędzi
				distance_matrix[ai][n1] - distance_matrix[aj][n2] // dystansy przed zamianą krawędzi
			if delta < 0 {
				moves_node = append(moves_node, MoveEdgeDetail{
					N1:    n1,
					N2:    n2,
					SN1:   ai,
					SN2:   aj,
					Delta: delta,
					Cycle: cycle,
				})
				moves_node = append(moves_node, MoveEdgeDetail{ // krawędź w drugą stronę
					N1:    ai,
					N2:    aj,
					SN1:   n1,
					SN2:   n2,
					Delta: delta,
					Cycle: cycle,
				})
			}
		}
	}

	return moves_node
}

func CheckApplicability(move Move, order [][]int) Applicability {
	switch m := move.(type) {
	case *MoveEdgeDetail:
		i := utils.IndexOf(order[m.Cycle], m.N1)
		j := utils.IndexOf(order[m.Cycle], m.N2)
		if i == -1 || j == -1 { // jeśli nie w tym samym cyklu
			return NotApplicable
		}
		ai, aj := utils.ElemAfter(order[m.Cycle], i), utils.ElemAfter(order[m.Cycle], j)   // wierzchołki po i i j w cyklu
		bi, bj := utils.ElemBefore(order[m.Cycle], i), utils.ElemBefore(order[m.Cycle], j) // wierzchołki przed i i j w cyklu

		if m.SN1 != ai && m.SN2 != aj && m.SN1 != bi && m.SN2 != bj { // jeśli różne następne wierzchołki niż wcześniej
			return NotApplicable
		}
		if m.SN1 == ai && m.SN2 == aj { // jeśli wszystko takie jak wcześniej
			return Applicable
		}
		if (m.SN1 == ai || m.SN2 == aj) && (m.SN1 == ai || m.SN1 == bi) && (m.SN2 == aj || m.SN2 == bj) { // jeśli różne sąsiedzi niż wcześniej
			return MayBeApplicable // jakaś krawędź w drugą stronę
		} else {
			return NotApplicable
		}

	case *SwapMoveDetail:
		i := utils.IndexOf(order[0], m.N1)
		j := utils.IndexOf(order[1], m.N2)
		if i == -1 || j == -1 { // jeśli różne cykle niż wcześniej
			return NotApplicable
		}
		ai := utils.ElemAfter(order[0], i)                            // wierzchołek po i w cyklu 1
		aj := utils.ElemAfter(order[1], j)                            // wierzchołek po j w cyklu 2
		bi := utils.ElemBefore(order[0], i)                           // wierzchołek przed i w cyklu 1
		bj := utils.ElemBefore(order[1], j)                           // wierzchołek przed j w cyklu 2
		if m.PN1 != bi || m.PN2 != bj || m.SN1 != ai || m.SN2 != aj { // jeśli różni sąsiedzi niż wcześniej
			return NotApplicable
		}
		return Applicable
	}

	return Applicable
}

func FindNewMoves(distance_matrix [][]int, order [][]int, move Move) ([]Move, error) {
	var (
		delta     int               // zmiana długości cyklu po dodaniu krawędzi)
		new_moves []Move = []Move{} // nowe ruchy do dodania
	)
	nodes_inner := [2][]int{{}, {}} // wierzchołki do rozważenia po zmianach krawędzi, bierzemy pod uwagę nowe krawędzie N1-N2, SN1-SN2
	nodes_outer := [2][]int{{}, {}} // wierzchołki cyklu 0 do rozważenia po zmianach krawędzi
	indexes_inner := [2][]int{{}, {}}
	indexes_outer := [2][]int{{}, {}}

	switch m := move.(type) {
	case *SwapMoveDetail:
		// na nowo obliczyć dla wszystkich wierzchołków
		nodes_outer[0] = []int{m.N2, m.SN1, m.PN1}                                              // wierzchołki do rozważenia przy zamianie wierzchołków
		nodes_outer[1] = []int{m.N1, m.SN2, m.PN2}                                              // drugi cykl
		indexes_outer[0] = utils.IndexesOf(order[0], nodes_outer[0])                            // znajdź indeksy w cyklu
		indexes_outer[1] = utils.IndexesOf(order[1], nodes_outer[1])                            // drugi cykl
		nodes_inner[0] = []int{m.N2, utils.ElemBefore(order[0], utils.IndexOf(order[0], m.N2))} // wierzchołki do rozważenia przy zamianie wierzchołków
		nodes_inner[1] = []int{m.N1, utils.ElemBefore(order[1], utils.IndexOf(order[1], m.N1))} // drugi cykl
		indexes_inner[0] = utils.IndexesOf(order[0], nodes_inner[0])                            // znajdź indeksy w cyklu
		indexes_inner[1] = utils.IndexesOf(order[1], nodes_inner[1])                            // drugi cykl

	case *MoveEdgeDetail:
		// nowe krawędzie: N1 - N2, SN1 - SN2, usunięcie krawędzi N1 - SN1, N2 - SN2
		if m.Cycle == 0 {
			nodes_inner[0] = []int{m.N1, m.SN1}
			nodes_outer[0] = []int{m.N1, m.N2, m.SN1, m.SN2}
		} else {
			nodes_inner[1] = []int{m.N1, m.SN1}
			nodes_outer[1] = []int{m.N1, m.N2, m.SN1, m.SN2}
		}
		indexes_inner[0] = utils.IndexesOf(order[m.Cycle], nodes_inner[0]) // znajdź indeksy w cyklu
		indexes_inner[1] = utils.IndexesOf(order[m.Cycle], nodes_inner[1]) // znajdź indeksy w cyklu
		indexes_outer[0] = utils.IndexesOf(order[m.Cycle], nodes_outer[0]) // znajdź indeksy w cyklu
		indexes_outer[1] = utils.IndexesOf(order[m.Cycle], nodes_outer[1]) // znajdź indeksy w cyklu
	}

	for c, no := range nodes_outer {
		for z, n1 := range no {
			i := indexes_outer[c][z] // indeks w cyklu

			cycle := c
			other_cycle := 1 - c                    // drugi cykl
			bi := utils.ElemBefore(order[cycle], i) // wierzchołek przed i w cyklu 1
			ai := utils.ElemAfter(order[cycle], i)  // wierzchołek po i w cyklu 1
			moves_node := []SwapMoveDetail{}        // aktualnie dostępne ruchy

			for j := 0; j < len(order[other_cycle]); j++ {
				n2 := order[other_cycle][j]

				bj := utils.ElemBefore(order[other_cycle], j) // wierzchołek przed j w cyklu 2
				aj := utils.ElemAfter(order[other_cycle], j)  // wierzchołek po j w cyklu 2

				delta := distance_matrix[bi][n2] + distance_matrix[n2][ai] + // dystansy od wierzchołków przed i po aktualnych po zamianie
					distance_matrix[bj][n1] + distance_matrix[n1][aj] -
					distance_matrix[bi][n1] - distance_matrix[bj][n2] - // dystansy przed zamianą krawędzi
					distance_matrix[ai][n1] - distance_matrix[aj][n2] // dystansy po zamianie krawędzi
				if delta < 0 {
					// dodaj ruch do listy
					if cycle == 0 {
						moves_node = append(moves_node, SwapMoveDetail{
							N1:    n1,
							N2:    n2,
							SN1:   ai,
							SN2:   aj,
							PN1:   bi,
							PN2:   bj,
							Delta: delta,
						})
					} else {
						moves_node = append(moves_node, SwapMoveDetail{
							N1:    n2,
							N2:    n1,
							SN1:   aj,
							SN2:   ai,
							PN1:   bj,
							PN2:   bi,
							Delta: delta,
						})
					}
				}
			}
			for move := range moves_node {
				new_moves = append(new_moves, &moves_node[move])
			}
		}
	}

	for c, ni := range nodes_inner {
		for n := range ni {
			n1 := ni[n]                        // wierzchołek 1
			i := indexes_inner[c][n]           // indeks w cyklu
			ai := utils.ElemAfter(order[c], i) // wierzchołek po i w cyklu

			moves_node := []MoveEdgeDetail{} // aktualnie dostępne ruchy

			for j := 0; j < len(order[c]); j++ {
				n2 := order[c][j]
				if n2 == n1 {
					continue
				}
				aj := utils.ElemAfter(order[c], j)

				delta = distance_matrix[n1][n2] + distance_matrix[ai][aj] - // dystansy po zamianie krawędzi
					distance_matrix[ai][n1] - distance_matrix[aj][n2] // dystansy przed zamianą krawędzi
				if delta < 0 {
					moves_node = append(moves_node, MoveEdgeDetail{
						N1:    n1,
						N2:    n2,
						SN1:   ai,
						SN2:   aj,
						Delta: delta,
						Cycle: c,
					})
					moves_node = append(moves_node, MoveEdgeDetail{ // krawędź w drugą stronę
						N1:    ai,
						N2:    aj,
						SN1:   n1,
						SN2:   n2,
						Delta: delta,
						Cycle: c,
					})
				}

			}
			for move := range moves_node {
				new_moves = append(new_moves, &moves_node[move])
			}
		}
	}

	return new_moves, nil
}

func AddNewMoves(s []Move, new_moves []Move) []Move {
	if len(new_moves) == 0 {
		return s
	}

	// sortowanie new_moves po delcie
	sort.Slice(new_moves, func(i, j int) bool {
		return new_moves[i].GetDelta() < new_moves[j].GetDelta()
	})

	result := make([]Move, 0, len(s)+len(new_moves))
	i, j := 0, 0

	for i < len(s) && j < len(new_moves) {
		if s[i].GetDelta() < new_moves[j].GetDelta() { // stare lepsze
			result = append(result, s[i])
			i++
		} else { // nowe lepsze
			result = append(result, new_moves[j])
			j++
		}
	}

	// pozostałr ruchu po wypełnieniu 1 z list
	result = append(result, s[i:]...)
	result = append(result, new_moves[j:]...)

	return result
}

func AddSorted(s []Move, move Move) []Move {
	for i, v := range s {
		if move.GetDelta() <= v.GetDelta() {
			s = utils.Insert(s, i, move) // dodaj ruch w odpowiednie miejsce
			return s
		}
	}

	s = append(s, move)
	return s
}

type Pair[T comparable] struct {
	A T
	B T
}

func AllCandidateMoves(distance_matrix [][]int, order [][]int, candidates [][]int, which_cycle map[int]int) ([]Move, error) {
	var (
		delta           int                                               // zmiana długości cyklu po dodaniu krawędzi
		moves_edge      []MoveEdgeDetail                                  // ruchy zamiany krawędzi
		moves_swap      []SwapMoveDetail                                  // ruchy zamiany wierzchołków między cyklami
		candidate_moves []Move                                            // wyszystkie ruchy
		num_nodes       int              = len(distance_matrix)           // liczba wierzchołków
		pairs           []Pair[int]                                       // pary wierzchołków/początek krawędzi do zamiany
		nodeToIndex     []map[int]int    = make([]map[int]int, num_nodes) // mapa wierzchołków do indeksów
	)
	for i := range order {
		nodeToIndex[i] = make(map[int]int, len(order[i]))
		for j, n := range order[i] {
			nodeToIndex[i][n] = j // mapa wierzchołków do indeksów
		}
	}

	for i := 0; i < num_nodes; i++ {
		cycle := which_cycle[i]                       // w którym cyklu jest dany wierzchołek
		index_i := nodeToIndex[cycle][i]              // indeks i w cyklu
		ai := utils.ElemAfter(order[cycle], index_i)  // wierzchołek po i w cyklu
		bi := utils.ElemBefore(order[cycle], index_i) // wierzchołek przed i w cyklu

		for _, candidate := range candidates[i] {
			if candidate == ai || candidate == bi { // jeśli już sąsiedzi nie rozważamy
				continue
			}
			cycle_candidate := which_cycle[candidate]                       // w którym cyklu jest dany wierzchołek
			index_candidate := nodeToIndex[cycle_candidate][candidate]      // indeks kandydata w cyklu
			aj := utils.ElemAfter(order[cycle_candidate], index_candidate)  // wierzchołek po candidate w cyklu
			bj := utils.ElemBefore(order[cycle_candidate], index_candidate) // wierzchołek przed candidate w cyklu

			if cycle == cycle_candidate { // jeśli w tym samym cyklu -> zamiana krawędzi
				pairs = []Pair[int]{
					{A: i, B: candidate},
					{A: bi, B: bj},
				} // wierchołki rozpoczynające krawędź do zamiany by otrzymać krawędź i-candidate

				for _, pair := range pairs {
					a, b := pair.A, pair.B
					index_a, index_b := nodeToIndex[cycle_candidate][a], nodeToIndex[cycle_candidate][b] // indeksy w cyklu
					aa := utils.ElemAfter(order[cycle], index_a)                                         // wierzchołek po a w cyklu
					ab := utils.ElemAfter(order[cycle], index_b)                                         // wierzchołek po b w cyklu

					delta = distance_matrix[a][b] + distance_matrix[aa][ab] - // dystansy po zamianie krawędzi
						distance_matrix[a][aa] - distance_matrix[b][ab] // dystansy przed zamianą krawędzi

					moves_edge = append(moves_edge, MoveEdgeDetail{
						N1:    order[cycle][index_a],
						N2:    order[cycle][index_b],
						SN1:   aa,
						SN2:   ab,
						Delta: delta,
						Cycle: cycle,
					})
				}
			} else { // jeśli w różnych cyklach -> zamiana wierzchołków
				if cycle == 0 { // pierwszy z pary cykl = 0
					pairs = []Pair[int]{
						{A: i, B: bj},
						{A: i, B: aj},
						{A: bi, B: candidate},
						{A: ai, B: candidate},
					} // pary wierzchołków do zamiany - wierzchołki obok tego z którym chcemy mieć krawędź
				} else {
					pairs = []Pair[int]{
						{A: bj, B: i},
						{A: aj, B: i},
						{A: candidate, B: bi},
						{A: candidate, B: ai},
					} // pary wierzchołków do zamiany - wierzchołki obok tego z którym chcemy mieć krawędź
				}

				for _, pair := range pairs {
					a, b := pair.A, pair.B
					index_a, index_b := nodeToIndex[0][a], nodeToIndex[1][b] // indeksy w cyklu
					aa := utils.ElemAfter(order[0], index_a)                 // wierzchołek po a w cyklu
					ab := utils.ElemAfter(order[1], index_b)                 // wierzchołek po b w cyklu
					ba := utils.ElemBefore(order[0], index_a)                // wierzchołek przed a w cyklu
					bb := utils.ElemBefore(order[1], index_b)                // wierzchołek przed b w cyklu

					delta = distance_matrix[ba][b] + distance_matrix[b][aa] + // dystansy od wierzchołków przed i po aktualnych po zamianie
						distance_matrix[bb][a] + distance_matrix[a][ab] -
						distance_matrix[ba][a] - distance_matrix[bb][b] - // dystansy przed zamianą krawędzi
						distance_matrix[aa][a] - distance_matrix[ab][b] // dystansy po zamianie krawędzi

					moves_swap = append(moves_swap, SwapMoveDetail{
						N1:    a,
						N2:    b,
						SN1:   aa,
						SN2:   ab,
						PN1:   ba,
						PN2:   bb,
						Delta: delta,
					})
				}
			}
		}
	}

	for i := range moves_edge {
		candidate_moves = append(candidate_moves, &moves_edge[i])
	}
	for i := range moves_swap {
		candidate_moves = append(candidate_moves, &moves_swap[i])
	}

	return candidate_moves, nil
}

func CandidateSearch(distance_matrix [][]int, order [][]int) error {
	var (
		candidate_moves []Move      // aktualnie dostępne ruchy
		candidates      [][]int     // numery wierzchołków kandydackich dla każdego wierzchołka
		top_candidates  int         = 10
		which_cycle     map[int]int // w którym cyklu jest dany wierzchołek
	)

	candidates = CalculateCandidates(distance_matrix, top_candidates) // obliczanie kandydatów
	// inicjalizacja which_cycle
	which_cycle = make(map[int]int)
	for c := range order {
		for i := range order[c] {
			which_cycle[order[c][i]] = c // przypisanie cyklu do wierzchołka
		}
	}

	var (
		best_move       Move  = nil                                                // najlepszy ruch w iteracji
		min_delta       int   = math.MaxInt                                        // minimalna zmiana długości cyklu
		current_length1 int   = utils.CalculateCycleLen(order[0], distance_matrix) // akutalna długość cyklu 1
		current_length2 int   = utils.CalculateCycleLen(order[1], distance_matrix) // akutalna długość cyklu 2
		current_length  int   = current_length1 + current_length2
		err             error = nil
	)

	for {
		// ruchy pomiędzy cyklami
		candidate_moves, err = AllCandidateMoves(distance_matrix, order, candidates, which_cycle) // wszystkie ruchy między cyklami
		if err != nil {
			return err
		}

		best_move, min_delta = FindBestMove(candidate_moves) // najlepszy ruch i minimalna zmiana długości cyklu

		// koniec iteracji
		if min_delta >= 0 { // jeśli nie znaleziono ruchu, który zmniejsza długość cyklu skończ przeszukiwanie
			best_move = nil // ustaw najlepszy ruch na nil
			break
		}
		// jeśli znaleziono ruch, to wykonaj go
		best_move.ExecuteMove(order) // wykonaj najlepszy ruch

		// aktualizacja which_cycle
		bm, ok := best_move.(*SwapMoveDetail)
		if ok { // jeśli ruch to zamiana wierzchołków
			// zamień cykle
			which_cycle[bm.N1] = 1 - which_cycle[bm.N1] // zamień cykle
			which_cycle[bm.N2] = 1 - which_cycle[bm.N2] // zamień cykle
		}

		current_length = current_length + min_delta // aktualizuj długość cyklu
		best_move, min_delta = nil, math.MaxInt     // ustaw najlepszy ruch na nil i delta MaxInt
	}

	return nil
}

func CalculateCandidates(distance_matrix [][]int, top_candidates int) (candidates [][]int) {
	candidates = make([][]int, len(distance_matrix)) // numery wierzchołków kandydackich dla każdego wierzchołka

	for i := 0; i < len(distance_matrix); i++ {
		for j := 0; j < len(distance_matrix); j++ {
			var (
				c     int  // aktualny kandydat
				k     int  // nr kandydata
				added bool = false
			)
			if i == j {
				continue
			}

			dist := distance_matrix[i][j] // dystans między i - aktualny wierzchołek, a j - potencjalny kandydat
			for k, c = range candidates[i] {
				if dist < distance_matrix[i][c] { // jeśli dystans mniejszy niż aktualny kandydat
					added = true
					candidates[i] = utils.Insert(candidates[i], k, j) // dodaj kandydata w odpowiednie miejsce
					if len(candidates[i]) > top_candidates {          // jeśli za dużo kandydatów
						candidates[i] = candidates[i][:top_candidates] // ogranicz do top_candidates
					}
					break
				}
			}
			if !added && len(candidates[i]) < top_candidates { // jeśli nie dodano kandydata i nie ma wszystkich
				candidates[i] = append(candidates[i], j) // dodaj kandydata na koniec
			}
		}
	}

	return
}
