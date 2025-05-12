package solver

import (
	"IMO/reader"
	"IMO/utils"
	"math"
	"math/rand"
	"time"
)

func MSLS(distance_matrix [][]int, order [][]int, nodes []reader.Node, num_iterations int) (int, error) {
	var (
		cost       int     = math.MaxInt              // koszt rozwiązania najlepszego
		length     int                                // długość aktualnych cykli
		best_order [][]int = make([][]int, NumCycles) // najlepsze cykle
	)
	for _ = range num_iterations { // pusta pętla
		err := Random(distance_matrix, order, nodes) // losu losu
		if err != nil {
			panic("Error")
		}
		err = SteepestEdge(distance_matrix, order) // najlepszy algorytm ls
		if err != nil {
			panic("Error")
		}
		length = utils.CalculateCycleLen(order[0], distance_matrix) + utils.CalculateCycleLen(order[1], distance_matrix)
		if length < cost {
			cost = length
			utils.CopyCycles(best_order, order)
		}
	}
	utils.CopyCycles(order, best_order)
	return num_iterations, nil
}

func ILS(distance_matrix [][]int, order [][]int, nodes []reader.Node, alg_time int) (int, error) {
	var (
		cost               int     = math.MaxInt              // koszt rozwiązania najlepszego
		length             int                                // długość aktualnych cykli
		best_order         [][]int = make([][]int, NumCycles) // najlepsze cykle
		perturbation_ratio float32 = 0.3                      // współczynnik perturbacji
		start_time         time.Time = time.Now() // czas rozpoczęcia algorytmu
		iter               int = 0 // liczba iteracji
	)
	err := Random(distance_matrix, order, nodes) // losu losu startowe
	if err != nil {
		panic("Error")
	}
	err = SteepestEdge(distance_matrix, order) // startowy local search
	if err != nil {
		panic("Error")
	}
	utils.CopyCycles(best_order, order)
	for int(time.Since(start_time).Seconds()) < alg_time { // pętla czasowa
		utils.CopyCycles(order, best_order)
		err = Perturbarion(order, perturbation_ratio) // nałożenie perturbacji
		if err != nil {
			panic("Error")
		}
		err = SteepestEdge(distance_matrix, order) // local search w celu poprawy jakości
		if err != nil {
			panic("Error")
		}
		length = utils.CalculateCycleLen(order[0], distance_matrix) + utils.CalculateCycleLen(order[1], distance_matrix)
		if length < cost { // warunek na poprawę rozwiązania
			cost = length
			utils.CopyCycles(best_order, order)
		}
		iter += 1
	}
	utils.CopyCycles(order, best_order)
	return iter, nil
}
func LNSWithLS(distance_matrix [][]int, order [][]int, nodes []reader.Node, alg_time int) (int, error) {
	var (
		cost          int     = math.MaxInt              // koszt rozwiązania najlepszego
		length        int                                // długość aktualnych cykli
		best_order    [][]int = make([][]int, NumCycles) // najlepsze cykle
		destroy_ratio float32 = 0.3                      // współczynnik niszczenia
		start_time         time.Time = time.Now() // czas rozpoczęcia algorytmu
		iter               int = 0 // liczba iteracji
	)
	err := Random(distance_matrix, order, nodes) // losu losu startowe
	if err != nil {
		panic("Error")
	}
	err = SteepestEdge(distance_matrix, order) // startowy local search
	if err != nil {
		panic("Error")
	}
	utils.CopyCycles(best_order, order)
	for int(time.Since(start_time).Seconds()) < alg_time { // pętla czasowa
		utils.CopyCycles(order, best_order)
		err = Destroy(order, destroy_ratio) // niszczymy jakiś procent wierzchołków
		if err != nil {
			panic("Error")
		}
		err = Repair(order, distance_matrix, nodes) // naprawiamy szkody przy pomocy greedy cycle (tylko ta metoda działa dla naszej implementacji)
		if err != nil {
			panic("Error")
		}
		err = SteepestEdge(distance_matrix, order) // dodatkowy local search
		if err != nil {
			panic("Error")
		}
		length = utils.CalculateCycleLen(order[0], distance_matrix) + utils.CalculateCycleLen(order[1], distance_matrix)
		if length < cost {
			cost = length
			utils.CopyCycles(best_order, order)
		}
		iter += 1
	}
	utils.CopyCycles(order, best_order)
	return iter, nil
}
func LNSWithoutLS(distance_matrix [][]int, order [][]int, nodes []reader.Node, alg_time int) (int, error) {
	var (
		cost          int     = math.MaxInt              // koszt rozwiązania najlepszego
		length        int                                // długość aktualnych cykli
		best_order    [][]int = make([][]int, NumCycles) // najlepsze cykle
		destroy_ratio float32 = 0.3                      // współczynnik niszczenia
		start_time         time.Time = time.Now() // czas rozpoczęcia algorytmu
		iter               int = 0 // liczba iteracji
	)
	err := Random(distance_matrix, order, nodes) // losu losu startowe
	if err != nil {
		panic("Error")
	}
	err = SteepestEdge(distance_matrix, order) // local search startowy
	if err != nil {
		panic("Error")
	}
	utils.CopyCycles(best_order, order)
	for int(time.Since(start_time).Seconds()) < alg_time { // pętla czasowa
		utils.CopyCycles(order, best_order)
		err = Destroy(order, destroy_ratio) // niszyczymy ileś wierzchołków
		if err != nil {
			panic("Error")
		}
		err = Repair(order, distance_matrix, nodes) // naprawa przy pomocy greedy cycle
		if err != nil {
			panic("Error")
		}
		length = utils.CalculateCycleLen(order[0], distance_matrix) + utils.CalculateCycleLen(order[1], distance_matrix)
		if length < cost {
			cost = length
			utils.CopyCycles(best_order, order)
		}
		iter += 1
	}
	utils.CopyCycles(order, best_order)
	return iter, nil
}
func Perturbarion(order [][]int, perturbation_ratio float32) error {
	var (
		num_of_max_perturbation int  = int(perturbation_ratio * float32(len(order[0]))) // maksymalna liczba przemieszań
		num_of_perturbation_c1  int  = 0                                                // liczba przemieszań dla cyklu pierwszego
		num_of_perturbation_c2  int  = 0                                                // liczba przemieszań dla cyklu drugiego
		sw1                     int  = -1                                               // indeks zamiany 1
		sw2                     int  = -1                                               // indeks zamiany 2
		rand_move               int                                                     // indeks losowego ruchu od 0-2
		move                    Move                                                    // wykonywany losowy ruch
	)
	if num_of_max_perturbation == 0 {
		panic("Za niski współczynnik ")
	}
	for num_of_perturbation_c1 == 0 {
		num_of_perturbation_c1 = rand.Intn(num_of_max_perturbation) // losu losu ale tak by nie wylosować zera
	}
	for num_of_perturbation_c2 == 0 {
		num_of_perturbation_c2 = rand.Intn(num_of_max_perturbation)
	}

	for i := range int(math.Max(float64(num_of_perturbation_c1), float64(num_of_perturbation_c2))) {
		rand_move = rand.Intn(3)
		switch rand_move {
		case 0:
			if i < num_of_perturbation_c1 {
				sw1 = rand.Intn(len(order[0]))
				sw2 = rand.Intn(len(order[0]))

				move = &MoveEdge{Cycle: 0, N1: sw1, N2: sw2, Delta: 0}
				move.ExecuteMove(order)
			}
			if i < num_of_perturbation_c2 {
				sw1 = rand.Intn(len(order[1]))
				sw2 = rand.Intn(len(order[1]))

				move = &MoveEdge{Cycle: 0, N1: sw1, N2: sw2, Delta: 0} // zamiana krawędzi
				move.ExecuteMove(order)
			}
		case 1:
			if i < num_of_perturbation_c1 {
				sw1 = rand.Intn(len(order[0]))
				sw2 = rand.Intn(len(order[0]))

				move = &MoveNode{Cycle: 0, N1: sw1, N2: sw2, Delta: 0}
				move.ExecuteMove(order)
			}
			if i < num_of_perturbation_c2 {
				sw1 = rand.Intn(len(order[1]))
				sw2 = rand.Intn(len(order[1]))

				move = &MoveNode{Cycle: 0, N1: sw1, N2: sw2, Delta: 0} // zamiana krawędzi
				move.ExecuteMove(order)
			}
		case 2:
			sw1 = rand.Intn(len(order[0]))
			sw2 = rand.Intn(len(order[1]))

			move = &SwapMove{N1: sw1, N2: sw2, Delta: 0}
			move.ExecuteMove(order)
		}
	}
	return nil
}
func Destroy(order [][]int, destroy_ratio float32) error {
	var (
		delete_c1 int = int(destroy_ratio * float32(len(order[0]))) // wyznaczenie liczby wierzchołków do zniknięcia
		delete_c2 int = int(destroy_ratio * float32(len(order[1])))
		del       int
	)
	for i := range int(math.Max(float64(delete_c1), float64(delete_c2))) {
		if i < delete_c1 {
			del = rand.Intn(len(order[0]))         // losu do usunięcia
			order[0] = utils.Remove(order[0], del) // usuwanie losowego wierzchołka
		}
		if i < delete_c2 {
			del = rand.Intn(len(order[1]))
			order[1] = utils.Remove(order[1], del)
		}
	}
	return nil
}
func Repair(order [][]int, distance_matrix [][]int, nodes []reader.Node) error {
	err := ContinueGreedyCycle(distance_matrix, order, nodes) // modyfikacja greedy cycle do kontunuuacji budowy cyklu
	if err != nil {
		panic("Error")
	}
	return nil
}
