package main

import (
	"fmt"
	"os"

	"zad1/reader"
)

// użycie: go run main.go <ścieżka_do_instancji>
func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("usage: go run main.go <path_to_instance>")
		return
	}
	nodes, headers, err := reader.ReadInstance(args[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(nodes)
	fmt.Println(headers)
}
