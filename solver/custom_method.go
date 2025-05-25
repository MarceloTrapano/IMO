package solver

import (
	"IMO/reader"
	"IMO/utils"
	"fmt"
	"time"
)

func EqualSplitWithLS(distance_matrix [][]int, order [][]int, nodes []reader.Node, time_limit int, heuristic_algorithm string, local_search_algorithm string, parameters ...int) (int, error) {
	// split nodes into 50/50 and use heuristic algorithm to make cycle from each half.
	node1, node2, err := PickRandomFarthest(distance_matrix, nodes)
	// node1, node2, err := PickRandomNodes(nodes)
	if err != nil {
		return 0, err
	}
	// fmt.Println("node1:", node1, "node2:", node2)
	res, err := FindBestSplit(distance_matrix, node1, node2)
	if err != nil {
		return 0, err
	}
	for i := range res {
		SingleGreedyCycle(distance_matrix, res[i])
	}

	// local search
	res, err = Local_search(res, local_search_algorithm, distance_matrix)
	if err != nil {
		return 0, err
	}

	utils.CopyCycles(order, res)

	return 0, nil
}

func EqualSplitWithoutLS(distance_matrix [][]int, order [][]int, nodes []reader.Node, time_limit int, heuristic_algorithm string, local_search_algorithm string, parameters ...int) (int, error) {
	// split nodes into 50/50 and use heuristic algorithm to make cycle from each half.
	node1, node2, err := PickRandomFarthest(distance_matrix, nodes)
	// node1, node2, err := PickRandomNodes(nodes)
	if err != nil {
		return 0, err
	}
	// fmt.Println("node1:", node1, "node2:", node2)
	res, err := FindBestSplit(distance_matrix, node1, node2)
	if err != nil {
		return 0, err
	}
	for i := range res {
		SingleGreedyCycle(distance_matrix, res[i])
	}
	utils.CopyCycles(order, res)

	return 0, nil
}

func FindBestSplit(distance_matrix [][]int, node1, node2 int) ([][]int, error) {
	var (
		order [][]int = make([][]int, NumCycles)
		diffs [][]int = make([][]int, NumCycles)
	)
	for i := range distance_matrix {
		if i == node1 || i == node2 { // patrzymy pozostałe wierzchołki
			continue
		}

		closer, diff, err := utils.CloserNode(distance_matrix, node1, node2, i)
		if err != nil {
			return nil, err
		}
		ind := utils.IndexBetterInSortedArray(diffs[closer], diff)
		if ind == -1 {
			order[closer] = append(order[closer], i)
			diffs[closer] = append(diffs[closer], diff)
		} else {
			order[closer] = utils.Insert(order[closer], ind, i)
			diffs[closer] = utils.Insert(diffs[closer], ind, diff)
		}
	}

	// len1 := len(order[0])
	// len2 := len(order[1])
	// fmt.Println("len1:", len1, "len2:", len2)
	// fmt.Println("order[0]:", order[0])
	// fmt.Println("order[1]:", order[1])
	// fmt.Println("diffs[0]:", diffs[0])
	// fmt.Println("diffs[1]:", diffs[1])
	utils.EvenCycles(order)
	// len1 = len(order[0])
	// len2 = len(order[1])
	// fmt.Println("len1:", len1, "len2:", len2)
	order[0] = append(order[0], node1)
	order[1] = append(order[1], node2)

	return order, nil
}

func EqualSplitStartHAE(distance_matrix [][]int, order [][]int, nodes []reader.Node, time_limit int, heuristic_algorithm string, local_search_algorithm string, parameters ...int) (int, error) {
	var (
		iter                  int                        // wykonane iteracje
		population            [][][]int                  // eltarna
		population_cycles_len []int                      // długości cykli
		start_time            time.Time     = time.Now() // czas rozpoczęcia algorytmu
		elapsed               time.Duration              // czas od rozpoczęnia algorytmu
		time_limit_reached    bool
		population_size       int = 20 // rozmiar populacji
		equal_split_size      int = 10 // część populacji wygenerowana przez EqualSplit
		// equal_split_gc_size   int = 6                      // część populacji wygenerowana przez EqualSplit z Greedy Cycle
		cycle_nodes int = len(nodes) / NumCycles // liczba wierzchołków w cyklu
	)
	// parametry jeśli podano
	if len(parameters) > 0 {
		// fmt.Println("Użyte parametry:", parameters)
		population_size = parameters[0]
	}

	if len(parameters) > 1 {
		equal_split_size = parameters[1]
	}

	if equal_split_size > population_size {
		return iter, fmt.Errorf("błąd: equal_split_size (%d) nie może być większe niż population_size (%d)", equal_split_size, population_size)
	}

	// 1. Stworzenie populacji elitarnej
	if equal_split_size == population_size {
		population = make([][][]int, 0, population_size)
		population_cycles_len = make([]int, 0, population_size)
	} else {
		population, population_cycles_len = CreateStartPopulation(distance_matrix, nodes, population_size-equal_split_size, heuristic_algorithm, local_search_algorithm)
	}

	// część populacji wygenerowana przez EqualSplit
	for i := 0; i < equal_split_size; i++ {
		res := make([][]int, NumCycles)
		for j := range res {
			res[j] = make([]int, cycle_nodes)
		}

		iter, err := EqualSplitWithLS(distance_matrix, res, nodes, time_limit, heuristic_algorithm, local_search_algorithm)
		if err != nil {
			return iter, err
		}
		len_order := utils.CalculateCycleLen(res[0], distance_matrix) + utils.CalculateCycleLen(res[1], distance_matrix)
		ind := utils.IndexBetterInSortedArray(population_cycles_len, len_order)
		if prev_index := ind - 1; prev_index >= 0 && len_order == population_cycles_len[prev_index] {
			i--
			continue // jeśli taki cykl już istnieje to nie dodawaj go ponownie
		}
		if ind == -1 || (len(population) == 0 && ind == 0) {
			population = append(population, res)
			population_cycles_len = append(population_cycles_len, len_order)
			continue
		}
		population_cycles_len = utils.Insert(population_cycles_len, ind, len_order)
		population = utils.Insert(population, ind, res)
	}

	// // część populacji wygenerowana przez EqualSplit
	// for i := 0; i < equal_split_gc_size; i++ {
	// 	res := make([][]int, NumCycles)
	// 	for j := range res {
	// 		res[j] = make([]int, cycle_nodes)
	// 	}

	// 	iter, err := EqualSplitWithLSGC(distance_matrix, res, nodes, time_limit, heuristic_algorithm, local_search_algorithm)
	// 	if err != nil {
	// 		return iter, err
	// 	}
	// 	len_order := utils.CalculateCycleLen(res[0], distance_matrix) + utils.CalculateCycleLen(res[1], distance_matrix)
	// 	ind := utils.IndexBetterInSortedArray(population_cycles_len, len_order)
	// 	if prev_index := ind - 1; prev_index >= 0 && len_order == population_cycles_len[prev_index] {
	// 		i--
	// 		continue // jeśli taki cykl już istnieje to nie dodawaj go ponownie
	// 	}
	// 	if ind == -1 || (len(population) == 0 && ind == 0) {
	// 		population = append(population, res)
	// 		population_cycles_len = append(population_cycles_len, len_order)
	// 		continue
	// 	}
	// 	population_cycles_len = utils.Insert(population_cycles_len, ind, len_order)
	// 	population = utils.Insert(population, ind, res)
	// }

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
				// if index_better == 0 {
				// 	fmt.Println("Znaleziono lepsze rozwiązanie:", population_cycles_len, "iteracja:", iter)
				// } else {
				// 	fmt.Println("Nowe rozwiązanie:")
				// }
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

func EqualSplitStartHAEMutation(distance_matrix [][]int, order [][]int, nodes []reader.Node, time_limit int, heuristic_algorithm string, local_search_algorithm string, parameters ...int) (int, error) {
	var (
		iter                  int                        // wykonane iteracje
		population            [][][]int                  // eltarna
		population_cycles_len []int                      // długości cykli
		start_time            time.Time     = time.Now() // czas rozpoczęcia algorytmu
		elapsed               time.Duration              // czas od rozpoczęnia algorytmu
		time_limit_reached    bool
		population_size       int = 20 // rozmiar populacji
		equal_split_size      int = 10 // część populacji wygenerowana przez EqualSplit
		// equal_split_gc_size   int = 0                      // część populacji wygenerowana przez EqualSplit z Greedy Cycle
		cycle_nodes int = len(nodes) / NumCycles // liczba wierzchołków w cyklu
	)
	// parametry jeśli podano
	if len(parameters) > 0 {
		// fmt.Println("Użyte parametry:", parameters)
		population_size = parameters[0]
	}

	if len(parameters) > 1 {
		equal_split_size = parameters[1]
	}

	if equal_split_size > population_size {
		return iter, fmt.Errorf("błąd: equal_split_size (%d) nie może być większe niż population_size (%d)", equal_split_size, population_size)
	}

	// 1. Stworzenie populacji elitarnej
	if equal_split_size == population_size {
		population = make([][][]int, 0, population_size)
		population_cycles_len = make([]int, 0, population_size)
	} else {
		population, population_cycles_len = CreateStartPopulation(distance_matrix, nodes, population_size-equal_split_size, heuristic_algorithm, local_search_algorithm)
	}

	// część populacji wygenerowana przez EqualSplit
	for i := 0; i < equal_split_size; i++ {
		res := make([][]int, NumCycles)
		for j := range res {
			res[j] = make([]int, cycle_nodes)
		}

		iter, err := EqualSplitWithLS(distance_matrix, res, nodes, time_limit, heuristic_algorithm, local_search_algorithm)
		if err != nil {
			return iter, err
		}
		len_order := utils.CalculateCycleLen(res[0], distance_matrix) + utils.CalculateCycleLen(res[1], distance_matrix)
		ind := utils.IndexBetterInSortedArray(population_cycles_len, len_order)
		if prev_index := ind - 1; prev_index >= 0 && len_order == population_cycles_len[prev_index] {
			i--
			continue // jeśli taki cykl już istnieje to nie dodawaj go ponownie
		}
		if ind == -1 || (len(population) == 0 && ind == 0) {
			population = append(population, res)
			population_cycles_len = append(population_cycles_len, len_order)
			continue
		}
		population_cycles_len = utils.Insert(population_cycles_len, ind, len_order)
		population = utils.Insert(population, ind, res)
	}

	// // część populacji wygenerowana przez EqualSplit
	// for i := 0; i < equal_split_gc_size; i++ {
	// 	res := make([][]int, NumCycles)
	// 	for j := range res {
	// 		res[j] = make([]int, cycle_nodes)
	// 	}

	// 	iter, err := EqualSplitWithLSGC(distance_matrix, res, nodes, time_limit, heuristic_algorithm, local_search_algorithm)
	// 	if err != nil {
	// 		return iter, err
	// 	}
	// 	len_order := utils.CalculateCycleLen(res[0], distance_matrix) + utils.CalculateCycleLen(res[1], distance_matrix)
	// 	ind := utils.IndexBetterInSortedArray(population_cycles_len, len_order)
	// 	if prev_index := ind - 1; prev_index >= 0 && len_order == population_cycles_len[prev_index] {
	// 		i--
	// 		continue // jeśli taki cykl już istnieje to nie dodawaj go ponownie
	// 	}
	// 	if ind == -1 || (len(population) == 0 && ind == 0) {
	// 		population = append(population, res)
	// 		population_cycles_len = append(population_cycles_len, len_order)
	// 		continue
	// 	}
	// 	population_cycles_len = utils.Insert(population_cycles_len, ind, len_order)
	// 	population = utils.Insert(population, ind, res)
	// }

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
		// mutation
		new_order, err = Mutation(new_order, distance_matrix, nodes)
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
				// if index_better == 0 {
				// 	fmt.Println("Znaleziono lepsze rozwiązanie:", population_cycles_len, "iteracja:", iter)
				// } else {
				// 	fmt.Println("Nowe rozwiązanie:")
				// }
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

func Mutation(order [][]int, distance_matrix [][]int, nodes []reader.Node) ([][]int, error) {
	var (
		perturbation_ratio float32 = 0.2 // współczynnik perturbacji
		// destroy_ratio      float32 = 0.1  // współczynnik niszczenia
		err error
	)
	err = Perturbarion(order, perturbation_ratio) // nałożenie perturbacji
	if err != nil {
		return nil, fmt.Errorf("błąd podczas perturbacji: %w", err)
	}
	// err = Destroy(order, destroy_ratio) // niszczymy jakiś procent wierzchołków
	// if err != nil {
	// 	panic("Error")
	// }
	// err = Repair(order, distance_matrix, nodes) // naprawiamy szkody przy pomocy greedy cycle (tylko ta metoda działa dla naszej implementacji)
	// if err != nil {
	// 	panic("Error")
	// }

	return order, nil
}
