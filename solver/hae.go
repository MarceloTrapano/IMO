package solver

import (
	"IMO/reader"
	"IMO/utils"
	"time"
)

func SameSolution[T comparable](s1 [][]T, s2 [][]T) bool {
	for i := 0; i < len(s1); i++ {
		for j := 0; j < len(s1[i]); j++ {
			if s1[i][j] != s2[i][j] {
				return false // jeśli 1 wartość się nie zgadza to różne
			}
		}
	}
	return true
}

func CreateStartPopulation(distance_matrix [][]int, nodes []reader.Node, population_size int, heuristic_algorithm string, local_search_algorithm string) ([][][]int, []int) {
	var (
		population            [][][]int = make([][][]int, population_size) // eltarna
		population_cycles_len []int     = make([]int, population_size)     // długości cykli
	)

	// 1. Stworzenie populacji elitarnej
	for i := 0; i < population_size; i++ {
		start_order, err := Solve(nodes, heuristic_algorithm, distance_matrix) // domyślnie Random
		if err != nil {
			panic("Error")
		}
		ls_order, err := Local_search(start_order, local_search_algorithm, distance_matrix) // lokalne wyszukiwanie; domyślnie SteepestEdge
		if err != nil {
			panic("Error")
		}
		cycle_len := utils.CalculateCycleLen(ls_order[0], distance_matrix) + utils.CalculateCycleLen(ls_order[1], distance_matrix)
		index_better := utils.IndexBetterInSortedArray(population_cycles_len[:i], cycle_len)
		if index_better == -1 {
			index_better = i
		}

		// // sprawdzenie czy jest rozwiązanie - ta sama suma dł. cyklów, oraz te same wartości
		// for prev_index:=index_better-1; prev_index >= 0 && cycle_len == population_cycles_len[prev_index];{
		// 	if SameSolution(ls_order, population[prev_index]) == true {
		// 		// to samo rozwiązanie już jest

		// 	}
		// }
		if prev_index := index_better - 1; prev_index >= 0 && cycle_len == population_cycles_len[prev_index] {
			// ta sama suma cykli - najpewniej to samo rozwiązanie, ciężko o taką samą dł. cykli przy różnych rozwiązaniach
			// wygeneruj nowe
			i--
			continue
		}

		if index_better == i {
			population[i] = ls_order
			population_cycles_len[i] = cycle_len
		} else {
			utils.Insert(population_cycles_len, index_better, cycle_len)
			// insert ale na slice nie comparable do populacji
			// TODO
		}

	}
	return population, population_cycles_len
}

func CrossOver(p1 [][]int, p2 [][]int, distance_matrix [][]int) ([][]int, error) {
	var (
		crossed_order [][]int = make([][]int, len(p1))
	)

	// 1. Dla obu rodziców: Znajdź osobno dla każdego cyklu wszyskie krawędzie - utwórz macierz sąsiedztwa; jeśli jest krawędź w jedną stronę to utwórz krawędź w drugą stronę
	// 2. Połącz macierze sąsiedztwa obu rodziców, osobno dla każdego cyklu - operator AND
	// Teraz w wierszu mamy 3 możliwości: 0, 1, 2 sąsiedzi dla danego wierzchołka
	// 0 - nie ma sąsiadów nie będzie w pozostałym cyklu, dołączy do puli wierzchołków nieprzypisanych do żadnego cyklu
	// 1, 2 - dołączymy do cyklu
	// jeśli nie ma 1 i 0 to znaczy, że cykl taki sam w obu rodzicach, nie trzeba tworzyć pozostałego cyklu - drugi musi być zatem inny bo rozwiązania w populacji się nie powtarzają
	// zaczynamy od 1 - wierzchołki, które mają krawędź tylko z 1 strony - tworzą z kolejnymi sąsiadami łańcuchy, które chcemy połączyć, by mieć cykl do działania z GreedyCycle
	// łączymy wierzchołki z 1 sąsiadem z najbliższym innym wierzchołkiem z 1, które nie są w tym samym łańcuchu
	// łączenie po kolei od najmniejszego indeksu wierzchołka - nie zawsze stworzy się najkrótszy cykl, ale i tak robimy GreedyCycle; też większa różnorodność w populacji i szybkie działanie
	// ew. minimum ze wszystkich możliwości połączeń
	// ostatecznie z pojedynczego łańcuch tworzymy cykl - teraz to do GreedyCycle razem z wierzchołkami z 0
}

func HAEWithoutLS(distance_matrix [][]int, order [][]int, nodes []reader.Node, time_limit int, heuristic_algorithm string, local_search_algorithm string, population_size int) (int, error) {
	var (
		iter                  int                        // wykonane iteracje
		population            [][][]int                  // eltarna
		population_cycles_len []int                      // długości cykli
		start_time            time.Time     = time.Now() // czas rozpoczęcia algorytmu
		elapsed               time.Duration              // czas od rozpoczęnia algorytmu
		time_limit_reached    bool
	)

	// 1. Stworzenie populacji elitarnej
	population, population_cycles_len = CreateStartPopulation(distance_matrix, nodes, population_size, heuristic_algorithm, local_search_algorithm)
	var (
		p1, p2 [][]int // rodzice
	)

	// główna pętla algorytmu
	for time_limit_reached = false; !time_limit_reached; {
		// 2 losowi rodzice z populacji
		i1, i2, _ := utils.Pick2RandomValues(population_size)
		p1, p2 = population[i1], population[i2]

		// krzyżowanie rodziców

		// sprawdzenie czy koniec czasu
		elapsed = time.Since(start_time)
		if elapsed.Milliseconds() > int64(time_limit) {
			time_limit_reached = true
		}
		iter++
	}
	utils.CopyCycles(order, population[0]) // kopiowanie najlepszego rozwiązania do order
	return iter, nil
}

func HAEWithLS(distance_matrix [][]int, order [][]int, nodes []reader.Node, time_limit int, heuristic_algorithm string, local_search_algorithm string, population_size int) (int, error) {
	return 0, nil
}
