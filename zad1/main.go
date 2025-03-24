package main

import (
	"encoding/json"
	"fmt"
	"os"
	"zad1/reader"
	"zad1/solver"
)

type Soltution struct {
	Order [][]int         `json:"order"`
	Nodes []reader.Node `json:"unordered nodes"`
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

	order, err := solver.Solve(nodes, algorithm)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(order)

	solution := Soltution{Order: order, Nodes: nodes}

	finalJson, _ := json.MarshalIndent(solution, "", "\t")

	os.WriteFile("sol_kroA200.json", finalJson, 0644)



}
