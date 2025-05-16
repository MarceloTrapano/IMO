package solver

import (
	"IMO/reader"
	"IMO/utils"
	"fmt"
	"math"
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
		population            [][][]int // eltarna
		population_cycles_len []int     // długości cykli
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
			population = append(population, ls_order)                        // dodanie do populacji
			population_cycles_len = append(population_cycles_len, cycle_len) // dodanie do długości cykli
		} else {
			population_cycles_len = utils.Insert(population_cycles_len, index_better, cycle_len)
			population = utils.Insert(population, index_better, ls_order)
		}
	}

	return population, population_cycles_len
}

func CrossOver(p1 [][]int, p2 [][]int, distance_matrix [][]int, nodes []reader.Node) ([][]int, error) {
	var (
		crossed_order     [][]int    = make([][]int, len(p1))
		adjacency_matrix1 [][]bool                               // macierze sąsiedztwa dla p1
		adjacency_matrix2 [][]bool                               // macierze sąsiedztwa dla p2
		adjacency_crossed [][][]bool = make([][][]bool, len(p1)) // macierze sąsiedztwa połączone (2 cykle)
	)
	// return crossed_order, nil

	// ogólny opis algorytmu:
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

	for i := 0; i < len(p1); i++ { // iteracja po cyklach
		adjacency_matrix1 = make([][]bool, len(p1[i])+len(p2[i]))
		adjacency_matrix2 = make([][]bool, len(p1[i])+len(p2[i]))
		for j := 0; j < len(adjacency_matrix1); j++ { // iteracja po wierzchołkach
			adjacency_matrix1[j] = make([]bool, len(adjacency_matrix1))
			adjacency_matrix2[j] = make([]bool, len(adjacency_matrix1))
		}
		for j := 0; j < len(p1[i]); j++ {
			n1 := p1[i][j] // rodzic 1
			n2 := utils.ElemAfter(p1[i], j)
			adjacency_matrix1[n1][n2] = true
			adjacency_matrix1[n2][n1] = true

			n1 = p2[i][j] // rodzic 2
			n2 = utils.ElemAfter(p2[i], j)
			adjacency_matrix2[n1][n2] = true
			adjacency_matrix2[n2][n1] = true
		}
		// połączenie macierzy sąsiedztwa rodziców
		adjacency_AND := utils.MatrixLogicAND(adjacency_matrix1, adjacency_matrix2) // macierz sąsiedztwa połączona
		adjacency_crossed[i] = adjacency_AND
	}

	var node_in_chain []bool

	// iteracja po cyklach
	for i := 0; i < len(adjacency_crossed); i++ {
		var (
			chains             [][]int // łańcuchy
			nr_neighbors_node  []int   // liczba sąsiadów
			neighbors_nodes    [][]int // sąsiedzi
			nr_neighbors_count map[int]int
		)
		nr_neighbors_node = make([]int, len(adjacency_crossed[i]))
		nr_neighbors_count = make(map[int]int)
		neighbors_nodes = make([][]int, len(adjacency_crossed[i]))
		chains = make([][]int, 0)
		node_in_chain = make([]bool, len(adjacency_crossed[i]))

		// uzupełenienie nr_neighbors_node
		for j := 0; j < len(adjacency_crossed[i]); j++ {
			for k := 0; k < len(adjacency_crossed[i][j]); k++ {
				if adjacency_crossed[i][j][k] {
					nr_neighbors_node[j]++
					neighbors_nodes[j] = append(neighbors_nodes[j], k) // dodanie sąsiada
				}
			}
			nr_neighbors_count[nr_neighbors_node[j]]++
		}

		if nr_neighbors_count[1] == 0 { // ten sam cykl w obu rodzicach
			copy(crossed_order[i], p1[i])
			continue
		}

		// łączenie łańcuchów
		for j := 0; j < len(adjacency_crossed[i]); j++ {
			if node_in_chain[j] {
				continue // jeśli wierzchołek już w łańcuchu to pomiń
			}
			if nr_neighbors_node[j] == 1 { // 1 sąsiad
				chain := make([]int, 0)
				chain = append(chain, j) // dodanie wierzchołka do łańcucha
				node_in_chain[j] = true
				previous := j
				// dodanie sąsiada do łańcucha - rekurencja
				for neighbor := neighbors_nodes[j][0]; ; {
					chain = append(chain, neighbor) // dodanie sąsiada do łańcucha
					node_in_chain[neighbor] = true
					// sprawdzenie czy sąsiad ma 1 sąsiada
					if nr_neighbors_node[neighbor] == 1 {
						// koniec łańcucha
						break
					} else {
						// dodanie sąsiada do łańcucha
						if neighbors_nodes[neighbor][0] == previous {
							previous = neighbor
							neighbor = neighbors_nodes[neighbor][1]
						} else {
							previous = neighbor
							neighbor = neighbors_nodes[neighbor][0]
						}
					}
				}
				chains = append(chains, chain) // dodanie łańcucha do listy łańcuchów
			}
		}

		// dopóki nie ma 1 chaina
		for len(chains) > 1 {
			// odległość między pierwszym a pozostałymi łańcuchami (koniec, początek) - branie minimum
			// 2 wektory - koniec, początek do wszystkich wierzchołków w pozostałych łańcuchach, po iteracji kkoniec lub początek znowu minimum ze wszystkich pozostałych
			start, end := chains[0][0], chains[0][len(chains[0])-1]
			min_distance := math.MaxInt
			min_index := -1
			zero_end := false
			join_end := false
			for j := 1; j < len(chains); j++ {
				// sprawdzenie odległości między końcem a początkiem
				distance_start_start := distance_matrix[start][chains[j][0]]
				distance_end_start := distance_matrix[end][chains[j][0]]
				distance_start_end := distance_matrix[start][chains[j][len(chains[j])-1]]
				distance_end_end := distance_matrix[end][chains[j][len(chains[j])-1]]

				// wszystkie możliwości po kolei O(n) a nie O(log(n)) jakby można zrobić ale tylko 4 przypadki więc spoko
				if distance_start_end < min_distance {
					min_distance = distance_start_end
					min_index = j
					join_end = true
					zero_end = false
				}
				if distance_end_start < min_distance {
					min_distance = distance_end_start
					min_index = j
					join_end = false
					zero_end = true
				}
				if distance_start_start < min_distance {
					min_distance = distance_start_start
					min_index = j
					join_end = false
					zero_end = false
				}
				if distance_end_end < min_distance {
					min_distance = distance_end_end
					min_index = j
					join_end = true
					zero_end = true
				}
			}

			if min_index == -1 {
				return nil, fmt.Errorf("błąd: nie znaleziono indeksu")
			}
			if zero_end {
				// odwróć chain 0
				temp_chain := make([]int, 0)
				for j := len(chains[0]) - 1; j >= 0; j-- {
					temp_chain = append(temp_chain, chains[0][j])
				}
				chains[0] = temp_chain
			}
			// dodanie do cyklu
			if join_end {
				// dodaj do chaina 0 elementy chaina min_index zaczynając od końca
				for j := len(chains[min_index]) - 1; j >= 0; j-- {
					chains[0] = append(chains[0], chains[min_index][j])
				}
			}
			if !join_end {
				// dodaj do chaina 0 elementy chaina min_index zaczynając od początku
				chains[0] = append(chains[0], chains[min_index]...)
			}
			// usunięcie chaina min_index
			chains = append(chains[:min_index], chains[min_index+1:]...)
		}

		if len(chains) == 0 { // jeśli nie ma łańcuchów to nie ma cyklu - chyba niepotrzebny warunek
			continue
		}
		crossed_order[i] = append(crossed_order[i], chains[0]...) // dodanie cyklu do crossed_order
	}

	// naprawa cyklu
	err := Repair(crossed_order, distance_matrix, nodes)
	if err != nil {
		return nil, err
	}

	return crossed_order, nil
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
		p1, p2           [][]int                                               // rodzice
		used_parents     map[string]utils.Empty = make(map[string]utils.Empty) // Mapa przechowująca użyte kombinacje rodziców
		num_used_parents int                    = 0
		max_combinations int                    = population_size * (population_size - 1) / 2 // maksymalna liczba kombinacji rodziców
	)
	if len(population) != len(population_cycles_len) && len(population) != population_size {
		return iter, fmt.Errorf("błąd: populacja nie jest tej samej długości co długości cykli")
	}

	// główna pętla algorytmu
	for time_limit_reached = false; !time_limit_reached && max_combinations != num_used_parents; {
		// 2 losowi rodzice z populacji
		i1, i2, _ := utils.Pick2RandomValues(population_size)
		p1, p2 = population[i1], population[i2]
		// jak rodzice byli sprawdzani to ich nie sprawdzaj ponownie
		key := fmt.Sprintf("%d-%d", min(i1, i2), max(i1, i2))
		if _, exists := used_parents[key]; exists {
			continue
		}
		used_parents[key] = utils.Empty{}
		num_used_parents++

		// krzyżowanie rodziców
		new_order, err := CrossOver(p1, p2, distance_matrix, nodes)
		if err != nil {
			return iter, err
		}
		len_new_order := utils.CalculateCycleLen(new_order[0], distance_matrix) + utils.CalculateCycleLen(new_order[1], distance_matrix)

		if len_new_order < population_cycles_len[len(population_cycles_len)-1] {
			exists := false
			// jeśli nowy cykl jest krótszy od najdłuższego cyklu w populacji to dodaj go do populacji
			index_better := utils.IndexBetterInSortedArray(population_cycles_len, len_new_order)
			if prev_index := index_better - 1; prev_index >= 0 && len_new_order == population_cycles_len[prev_index] {
				exists = true
			}
			if !exists {
				utils.InsertRetainSize(population_cycles_len, len_new_order, index_better)
				utils.InsertRetainSize(population, new_order, index_better)

				used_parents = make(map[string]utils.Empty) // reset mapy użytych rodziców
				num_used_parents = 0
			}
		}

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
		p1, p2           [][]int                                               // rodzice
		used_parents     map[string]utils.Empty = make(map[string]utils.Empty) // Mapa przechowująca użyte kombinacje rodziców
		num_used_parents int                    = 0
		max_combinations int                    = population_size * (population_size - 1) / 2 // maksymalna liczba kombinacji rodziców
	)
	if len(population) != len(population_cycles_len) && len(population) != population_size {
		return iter, fmt.Errorf("błąd: populacja nie jest tej samej długości co długości cykli")
	}

	// główna pętla algorytmu
	for time_limit_reached = false; !time_limit_reached && max_combinations != num_used_parents; {
		// 2 losowi rodzice z populacji
		i1, i2, _ := utils.Pick2RandomValues(population_size)
		p1, p2 = population[i1], population[i2]
		// jak rodzice byli sprawdzani to ich nie sprawdzaj ponownie
		key := fmt.Sprintf("%d-%d", min(i1, i2), max(i1, i2))
		if _, exists := used_parents[key]; exists {
			continue
		}
		used_parents[key] = utils.Empty{}
		num_used_parents++

		// krzyżowanie rodziców
		new_order, err := CrossOver(p1, p2, distance_matrix, nodes)
		if err != nil {
			return iter, err
		}
		// local search
		new_order, err = Local_search(new_order, local_search_algorithm, distance_matrix)
		if err != nil {
			return iter, err
		}
		len_new_order := utils.CalculateCycleLen(new_order[0], distance_matrix) + utils.CalculateCycleLen(new_order[1], distance_matrix)

		if len_new_order < population_cycles_len[len(population_cycles_len)-1] {
			exists := false
			// jeśli nowy cykl jest krótszy od najdłuższego cyklu w populacji to dodaj go do populacji
			index_better := utils.IndexBetterInSortedArray(population_cycles_len, len_new_order)
			if prev_index := index_better - 1; prev_index >= 0 && len_new_order == population_cycles_len[prev_index] {
				exists = true
			}
			if !exists {
				utils.InsertRetainSize(population_cycles_len, len_new_order, index_better)
				utils.InsertRetainSize(population, new_order, index_better)

				used_parents = make(map[string]utils.Empty) // reset mapy użytych rodziców
				num_used_parents = 0
			}
		}

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
