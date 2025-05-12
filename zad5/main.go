package main

import (
	"IMO/reader"
	"IMO/solver"
	"IMO/utils"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"
)

type Solution struct {
	Result        [][]int       `json:"result"`
	Worst_Order   [][]int       `json:"worst order"`
	Best_Order    [][]int       `json:"best order"`
	Nodes         []reader.Node `json:"unordered nodes"`
	Times         []float64     `json:"times"`
	Longest_Time  float64       `json:"longest time"`
	Shortest_Time float64       `json:"shortest time"`
	Iter          []int         `json:"iterations"`
}

// użycie: go run main.go <ścieżka_do_instancji> [algorytm]
func main() {
	var (
		heuristic_algorithm    string = "rand"
		local_search_algorithm string = "se"
		time_limit             int
		use_local_search       bool
		population_size        int = 20
	)
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("usage: go run main.go <path_to_instance> [greedy heuristic] [time limit (ms)] [local search algorithm]")
		return
	}
	if len(args) > 1 {
		heuristic_algorithm = args[1]
	}
	if len(args) > 2 {
		time_limit_str := args[2]
		time_limit_float, err := strconv.Atoi(time_limit_str)
		if err != nil {
			panic(err)
		}
		time_limit = time_limit_float
	}
	if len(args) > 3 {
		local_search_algorithm = args[3]
		use_local_search = true
	}
	nodes, headers, err := reader.ReadInstance(args[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(nodes)
	fmt.Println(headers)

	var (
		distance_matrix [][]int = make([][]int, len(nodes))
		results         [][]int = make([][]int, 2)
		times           []time.Duration
		times_seconds   []float64
	)
	num_of_rep := 10
	for i := range distance_matrix {
		distance_matrix[i] = make([]int, len(nodes))
		for j := range distance_matrix[i] {
			distance_matrix[i][j] = solver.EucDist(nodes[i], nodes[j])
		}
	}
	results[0] = make([]int, num_of_rep)
	results[1] = make([]int, num_of_rep)
	var (
		iter          int
		iterations    []int
		best_order    [][]int
		worst_order   [][]int
		order         [][]int
		longest_time  time.Duration = time.Duration(0)
		shortest_time time.Duration = time.Duration(math.MaxInt64)
		start_time    time.Time
		elapsed       time.Duration
	)
	best_score := -1
	worst_score := -1

	for i := 0; i < num_of_rep; i++ {
		fmt.Printf("Trial: %d\n", i+1)
		start_time = time.Now()
		order, iter, err = solver.HAE(nodes, distance_matrix, time_limit, heuristic_algorithm, local_search_algorithm, use_local_search, population_size)
		elapsed = time.Since(start_time)
		if err != nil {
			fmt.Println(err)
			return
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		err = solver.ValidateOrder(order, nodes)
		if err != nil {
			fmt.Println(err)
			return
		}
		results[0][i] = utils.CalculateCycleLen(order[0], distance_matrix)
		results[1][i] = utils.CalculateCycleLen(order[1], distance_matrix)
		if results[0][i]+results[1][i] > worst_score {
			worst_score = results[0][i] + results[1][i]
			worst_order = append(order[:0:0], order...)
		}
		if best_score == -1 || results[0][i]+results[1][i] < best_score {
			best_score = results[0][i] + results[1][i]
			best_order = append(order[:0:0], order...)
		}
		if elapsed > longest_time {
			longest_time = elapsed
		}
		if elapsed < shortest_time {
			shortest_time = elapsed
		}
		iterations = append(iterations, iter)
		times = append(times, elapsed)
		times_seconds = append(times_seconds, elapsed.Seconds())
	}

	fmt.Printf("Times Duration: %v\n", times)
	fmt.Printf("Longest time millis: %v\n", longest_time.Milliseconds())
	fmt.Printf("Shortest time seconds: %v\n", shortest_time.Seconds()) // chyba to najlepiej - dodane do Solution

	solution := Solution{Iter: iterations, Result: results, Worst_Order: worst_order, Best_Order: best_order, Nodes: nodes, Times: times_seconds, Longest_Time: longest_time.Seconds(), Shortest_Time: shortest_time.Seconds()}

	finalJson, _ := json.MarshalIndent(solution, "", "\t")

	os.WriteFile("Res_RAND_ILS_KroB200.json", finalJson, 0644)
}
