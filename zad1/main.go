package main

import (
	"IMO/reader"
	"IMO/solver"
	"IMO/utils"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"time"
)

type Solution struct {
	Result      [][]int       `json:"result"`
	Worst_Order [][]int       `json:"worst order"`
	Best_Order  [][]int       `json:"best order"`
	Nodes       []reader.Node `json:"unordered nodes"`
}

// użycie: go run main.go <ścieżka_do_instancji> [algorytm]
func main() {
	var algorithm string
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("usage: go run main.go <path_to_instance> [algorithm]")
		return
	}
	if len(args) > 1 {
		algorithm = args[1]
	}
	nodes, headers, err := reader.ReadInstance(args[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(nodes)
	fmt.Println(headers)

	var (
		distance_matrix [][]int       = make([][]int, len(nodes))
		results         [][]int       = make([][]int, 2)
		longest_time    time.Duration = time.Duration(0)
		shortest_time   time.Duration = time.Duration(math.MaxInt64)
		start_time      time.Time
		elapsed         time.Duration
		times           []time.Duration
		times_milis     []float64
	)
	num_of_rep := 100
	for i := range distance_matrix {
		distance_matrix[i] = make([]int, len(nodes))
		for j := range distance_matrix[i] {
			distance_matrix[i][j] = solver.EucDist(nodes[i], nodes[j])
		}
	}
	results[0] = make([]int, num_of_rep)
	results[1] = make([]int, num_of_rep)
	var (
		best_order  [][]int
		worst_order [][]int
	)
	best_score := -1
	worst_score := -1

	for i := 0; i < num_of_rep; i++ {
		fmt.Printf("Trial: %d\n", i+1)
		start_time = time.Now()
		order, err := solver.Solve(nodes, algorithm, distance_matrix)
		elapsed = time.Since(start_time)
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
		times = append(times, elapsed)
		times_milis = append(times_milis, float64(elapsed.Milliseconds()))
	}
	sum := 0.0
	for _,i := range(times_milis){
		sum += i
	}
	mean := sum / float64(len(times_milis))
	fmt.Printf("Times Duration: %v\n", times)
	fmt.Printf("Mean time millis: %v\n", mean)
	fmt.Printf("Longest time millis: %v\n", longest_time.Milliseconds())
	fmt.Printf("Shortest time seconds: %v\n", shortest_time.Milliseconds())

	solution := Solution{Result: results, Worst_Order: worst_order, Best_Order: best_order, Nodes: nodes}

	finalJson, _ := json.MarshalIndent(solution, "", "\t")

	os.WriteFile("Res_GRE_kroa.json", finalJson, 0644)
}
