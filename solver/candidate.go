package solver

import (
	"IMO/utils"
	"sort"
)

func FastLocalSearch(distance_matrix [][]int, order [][]int) error {
	// inicjacja tablicy z najlepszymi ruchami
	var best_moves []Move       // aktualnie najlepsze ruchy posortowane od najlepszego do najgorszego
	
	for {
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
			return best_moves[i].GetDelta() > best_moves[j].GetDelta() // sortowanie po najwyższych deltach
		})
		for i, move := range best_moves{
			if CheckAplicability(move, order){ // sprawdzenie czy ruch jest aplikowalny
				move.ExecuteMove(order)
				break
			}
			if i == len(best_moves){ // jak nie znajdzie ruchu aplikowalnego to przerywa
				return nil
			}
		}
	}
}
func BestMovesBetweenCycles(distance_matrix [][]int, order [][]int, distances_before [][]int) ([]SwapMove, error) {
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
			if delta > 0 {
				// dodaj ruch do listy
				moves = append(moves, SwapMove{
					N1:    i,
					N2:    j,
					Delta: delta,
				})
			}
		}
	}

	return moves, nil
}

func BestMovesEdgesCycle(distance_matrix [][]int, order []int, cycle int) []MoveEdge {
	var (
		n1         int        // wierzchołek 1
		n2         int        // wierzchołek 2
		delta      int        // zmiana długości cyklu po dodaniu krawędzi
		moves_node []MoveEdge // aktualnie dostępne ruchy
	)

	// dla każdej pary wierzchołków w cyklu; kolejność nie ma znaczenia
	for i := 0; i < len(order); i++ {
		for j := i + 1; j < len(order); j++ {
			n1, n2 = order[i], order[j] // wierzchołki 1 i 2 - nr w cyklu
			ai := utils.ElemAfter(order, i)
			aj := utils.ElemAfter(order, j)
			delta = distance_matrix[n1][n2] + distance_matrix[ai][aj] - // dystansy po zamianie krawędzi
				distance_matrix[ai][n1] - distance_matrix[aj][n2] // dystansy przed zamianą krawędzi
			if delta > 0 {
				moves_node = append(moves_node, MoveEdge{
					N1:    i,
					N2:    j,
					Delta: delta,
					Cycle: cycle,
				})
			}
		}
	}
	return moves_node
}
func CheckAplicability(move Move, order [][]int) bool {
	switch m := move.(type){
	case *MoveEdge:
		for node := range order[m.Cycle]{
			// Tu jest problem bo w sumie nasze wykonywanie ruchu nie zapisuje które krawędzie usuwa, tylko wylicza samo na bazie aktualnego cyklu
			// i teraz nie wiem czy my musimy jakoś zapisywać, które krawędzie dodajemy i zabieramy
			// też kwestia interpretacji tutaj wchodzi, bo zamiana krawędzi u nas jest pod względem wierzchołka
			// gdzie dla cyklu 12345678 zamianę dwóch krawędzi interpretujemy jako stworzenie krawędzi z wierzchołka do wierzchołka
			// więc automatycznie jak chcemy walnąć krawędź 36 to program sam se wylicza pasujące 12365478 i wykonuje zamianę krawędzi 3-4 z 6-7
			// czyli w sumie powinniśmy gdzieś zapisywać krawędź 3-4 i 6-7. Hmm... 
		}
	}
	return true
}
