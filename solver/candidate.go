package solver

import (
	"IMO/utils"
	"fmt"
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
	// fmt.Println("Best moves: ")
	// for _, move := range best_moves {
	// 	fmt.Println(move)
	// }
	del := false

	for len(best_moves) > 0 {
		// distance_before := DistancesBefore(distance_matrix, order)
		// swap_moves, err := BestMovesBetweenCycles(distance_matrix, order, distance_before) // wybieranie ruchów między cyklami poprawiających wynik
		// if err != nil {
		// 	return err
		// }
		// for m := range swap_moves {
		// 	best_moves = append(best_moves, &swap_moves[m])
		// }
		// for c := 0; c < NumCycles; c++ {
		// 	moves_cycle := BestMovesEdgesCycle(distance_matrix, order[c], c) // wybieranie ruchów krawędzi poprawiających wynik
		// 	for m := range moves_cycle {
		// 		best_moves = append(best_moves, &moves_cycle[m])
		// 	}
		// }

		// fmt.Println("Best moves: ")
		// for _, move := range best_moves {
		// 	fmt.Println(move)
		// }
		// fmt.Println("Best moves count: ", len(best_moves))
		to_delete := []int{}  // indeksy do usunięcia
		new_moves := []Move{} // nowe ruchy do dodania
		// fmt.Println(new_moves)

	Loop: // label loopa do breakowania
		for i, move := range best_moves {
			applicability := CheckApplicability(move, order)
			switch applicability {
			case Applicable:
				del = false
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

			// if CheckApplicability(move, order) { // sprawdzenie czy ruch jest aplikowalny
			// 	move.ExecuteMove(order)

			// 	new_moves, err := FindNewMoves(distance_matrix, order, move) // znajdź nowe ruchy
			// 	break
			// } else {
			// 	// dodaj do usunięcia - ruch nie jest aplikowalny
			// 	to_delete = append(to_delete, i)
			// }

			if del && i == len(best_moves)-1 { // jak nie znajdzie ruchu aplikowalnego to przerywa
				// fmt.Println("No more moves")
				// fmt.Println("Best moves count: ", len(best_moves))
				// fmt.Println("To delete: ", to_delete)
				// fmt.Println("Orders: ")
				// fmt.Println(order[0])
				// fmt.Println(order[1])
				// fmt.Println("Best moves: ")
				// for _, move := range best_moves {
				// 	fmt.Println(move)
				// }
				return nil
			}
			if i == len(best_moves)-1 {
				// fmt.Println("No more moves")
				del = true // jeśli ostatni ruch to przerywamy
			}
		}
		best_moves = utils.RemoveIndexes(best_moves, to_delete) // usuń ruchy, które nie są aplikowalne
		best_moves = AddNewMoves(best_moves, new_moves)         // dodaj nowe ruchy
		best_val := best_moves[0].GetDelta()                    // najlepsza zmiana
		for _, move := range best_moves {
			// fmt.Println(move)
			if move.GetDelta() < best_val { // jeśli nowy ruch jest lepszy od najlepszego to dodaj go do najlepszych ruchów
				panic("FastLocalSearch: new move is better than best move")
			}
			best_val = move.GetDelta()
		}
	}
	fmt.Println("No more moves")
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
	var dodano int = 0 // liczba dodanych ruchów
	var all_moves int = 0

	// dla każdej pary wierzchołków w cyklu; kolejność nie ma znaczenia
	for i := 0; i < len(order); i++ {
		for j := i + 1; j < len(order); j++ {
			n1, n2 = order[i], order[j] // wierzchołki 1 i 2 - nr w cyklu
			ai := utils.ElemAfter(order, i)
			aj := utils.ElemAfter(order, j)
			delta = distance_matrix[n1][n2] + distance_matrix[ai][aj] - // dystansy po zamianie krawędzi
				distance_matrix[ai][n1] - distance_matrix[aj][n2] // dystansy przed zamianą krawędzi
			if delta < 0 {
				dodano++
				if n1 == n2 || ai == aj { // jeśli nie można zamienić krawędzi
					panic("BestMovesEdgesCycle: n1 == n2")
				}
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
			all_moves++
		}
	}
	// fmt.Println("Doano {} ruchów z {} wszystkich, procent", dodano, all_moves, float64(dodano)/float64(all_moves))

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

	switch m := move.(type) {
	case *SwapMoveDetail:
		// na nowo obliczyć dla wszystkich wierzchołków
		nodes := []int{m.N2, m.SN1, m.PN1, m.N1, m.SN2, m.PN2}                             // wierzchołki do rozważenia przy nowych krawędziach, bierzemy pod uwagę nowe krawędzie N1-N2, SN1-SN2
		indexes := utils.IndexesOf(order[0], []int{m.N2, m.SN1, m.PN1})                    // znajdź indeksy w cyklu
		indexes = append(indexes, utils.IndexesOf(order[1], []int{m.N1, m.SN2, m.PN2})...) // znajdź indeksy w cyklu
		if len(indexes) != 6 {
			panic("FindNewMoves: indexes not found")
		}
		// fmt.Println("SwapMove")

		for z, n1 := range nodes {
			i := indexes[z] // indeks w cyklu
			cycle := z / 3
			other_cycle := int(math.Abs(float64(cycle - 1))) // drugi cykl
			bi := utils.ElemBefore(order[cycle], i)          // wierzchołek przed i w cyklu 1
			ai := utils.ElemAfter(order[cycle], i)           // wierzchołek po i w cyklu 1
			moves_node := []SwapMoveDetail{}                 // aktualnie dostępne ruchy

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
					if n1 == n2 || ai == aj { // jeśli nie można zamienić krawędzi
						panic("FindNewMoves: n1 == n2")
					}
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
	case *MoveEdgeDetail:
		// fmt.Println("MoveEdge")
		// nowe krawędzie: N1 - N2, SN1 - SN2, usunięcie krawędzi N1 - SN1, N2 - SN2
		nodes := []int{m.N1, m.SN1}                       // wierzchołki do rozważenia przy nowych krawędziach, bierzemy pod uwagę nowe krawędzie N1-N2, SN1-SN2
		indexes := utils.IndexesOf(order[m.Cycle], nodes) // znajdź indeksy w cyklu
		if len(indexes) != 2 {
			panic("MoveEdgeDetail: indexes not found")
		}

		for n := range nodes {
			n1 := nodes[n]                                              // wierzchołek 1
			i := indexes[n]                                             // indeks w cyklu
			ai := utils.ElemAfter(order[m.Cycle], i)                    // wierzchołek po i w cyklu
			bi := utils.ElemBefore(order[m.Cycle], i)                   // wierzchołek przed i w cyklu
			if ai != m.N2 && ai != m.SN2 && bi != m.N2 && bi != m.SN2 { // jeśli nie jest to wierzchołek po N1 lub N2
				// fmt.Println(move)
				// fmt.Println("n1, i, ai, m.N2, m.SN2", n1, i, ai, m.N2, m.SN2)
				// fmt.Println("m.N1, m.N2, m.SN1, m.SN2", m.N1, m.N2, m.SN1, m.SN2)
				// fmt.Println("bi", utils.ElemBefore(order[m.Cycle], i))
				// fmt.Println("delta", m.Delta)
				// fmt.Println("order", order[m.Cycle])
				panic("MoveEdgeDetail: node not in move")
			}
			moves_node := []MoveEdgeDetail{} // aktualnie dostępne ruchy

			for j := 0; j < len(order[m.Cycle]); j++ {
				n2 := order[m.Cycle][j]
				if n2 == n1 {
					continue
				}
				aj := utils.ElemAfter(order[m.Cycle], j)

				delta = distance_matrix[n1][n2] + distance_matrix[ai][aj] - // dystansy po zamianie krawędzi
					distance_matrix[ai][n1] - distance_matrix[aj][n2] // dystansy przed zamianą krawędzi
				if delta < 0 {
					moves_node = append(moves_node, MoveEdgeDetail{
						N1:    n1,
						N2:    n2,
						SN1:   ai,
						SN2:   aj,
						Delta: delta,
						Cycle: m.Cycle,
					})
					moves_node = append(moves_node, MoveEdgeDetail{ // krawędź w drugą stronę
						N1:    ai,
						N2:    aj,
						SN1:   n1,
						SN2:   n2,
						Delta: delta,
						Cycle: m.Cycle,
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
	for _, move := range new_moves {
		s = AddSorted(s, move)
	}
	return s
}

func AddSorted(s []Move, move Move) []Move {
	for i, v := range s {
		if move.GetDelta() < v.GetDelta() {
			s = utils.Insert(s, i, move) // dodaj ruch w odpowiednie miejsce
			return s
		}
	}

	s = append(s, move)
	return s
}
