package utils

import (
	"bytes"
	"fmt"
	"math"
	"slices"
	"zad1/reader"
)

type Edge struct {
	Nr     int   // edge number
	From   int   // node 1 number
	To     int   // node 2 number
	Prev   *Edge // previous edge
	Next   *Edge // next edge
	Length int   // distance between From and To
}

type EdgeLinkedList struct {
	Edge  int
	Prev  *EdgeLinkedList
	Next  *EdgeLinkedList
	Value int
}

func (eLL *EdgeLinkedList) String() string {
	var buf bytes.Buffer
	buf.WriteString("{")
	for eLL != nil {
		buf.WriteString(fmt.Sprintf("%v: %v", eLL.Edge, eLL.Value))
		if eLL.Next != nil {
			buf.WriteString(", ")
		}
		eLL = eLL.Next
	}
	buf.WriteString("}")
	return buf.String()
}

func NewEdge(from int, to int, distance_matrix [][]int, next *Edge, prev *Edge) Edge {
	return Edge{
		From:   from,
		To:     to,
		Prev:   prev,
		Next:   next,
		Length: distance_matrix[from][to],
	}
}

// o ile zwiększy się cykl po dodaniu wierzchołka w miejsce krawędzi
func EdgeInsertValue(distance_matrix [][]int, node int, edge *Edge) int {
	return distance_matrix[node][edge.To] + distance_matrix[node][edge.From] - edge.Length
}

func UpdateDistances(eLL *EdgeLinkedList, distance_matrix [][]int, newEdges []EdgeLinkedList, newEdgesSorted bool) *EdgeLinkedList {
	if !newEdgesSorted {
		// sortuje rosnąco - chcemy malejąco (najlepsze na końcu) więc przeciwnie: j-i zamiast i-j
		slices.SortFunc(newEdges, func(i, j EdgeLinkedList) int {
			return j.Value - i.Value
		})
	}
	fmt.Println("")
	for eLL != nil && len(newEdges) > 0 {
		fmt.Println(eLL)
		for i := len(newEdges) - 1; i >= 0; i-- {
			edge := &newEdges[i]
			switch {
			case edge.Value < eLL.Value: // lepszy niż aktualny - dodaj przed (aktualizacja wskaźnika dla poprzedniego, teraz ten jest aktualnie rozważany)
				prev := eLL.Prev
				if prev != nil { // może być początek LListy
					prev.Next = edge
				}
				edge.Prev = prev
				edge.Next = eLL
				eLL.Prev = edge
				eLL = edge
				newEdges = newEdges[:i]
			case eLL.Next == nil: // koniec listy - dodaj na końcu; && edge.Value >= eLL.Value pominięte
				fmt.Println("koniec")
				edge.Prev = eLL
				eLL.Next = edge
				newEdges = newEdges[:i]
			default: // gorszy niż aktulany - szukaj dalej; edge.Value >= eLL.Value pominięte
				eLL = eLL.Next
			}
		}
		fmt.Println(eLL, eLL.Prev, eLL.Next)
	}
	for eLL.Prev != nil {
		eLL = eLL.Prev
	}
	return eLL // zwróć początek listy
}

func MatrixMax(matrix [][]int) (int, int, int) {
	max := math.MinInt64
	x := 0
	y := 0
	for i, row := range matrix {
		for j, value := range row {
			if value > max {
				max = value
				x = j
				y = i
			}
		}
	}
	return x, y, max
}

func Insert(array []int, i int, j int) []int {
	var new_arr []int
	if i < 0 || i >= len(array) {
		panic("Index out of range")
	}
	for idx, val := range array {
		if idx == i {
			new_arr = append(new_arr, j)
		}
		new_arr = append(new_arr, val)
	}
	return new_arr
}

func CalculateCycleLen(order []int, distance_matrix [][]int) int {
	cost := 0
	for i := range order {
		cost += distance_matrix[order[i]][order[(i+1)%len(order)]]
	}
	return cost
}

func FarthestNode(nodes []reader.Node, distance_matrix [][]int, node int, visited []bool) (farthest int, err error) {
	max := math.MinInt64
	for i := range distance_matrix[node] {
		if !visited[i] && i != node && distance_matrix[node][i] > max {
			max = distance_matrix[node][i]
			farthest = i
		}
	}
	if farthest == -1 {
		err = fmt.Errorf("no farthest node found for %v - invalid distance matrix", node)
	}
	return
}

func NearestNode(nodes []reader.Node, distance_matrix [][]int, node int, visited []bool) (nearest int, err error) {
	min := math.MaxInt64
	for i := range distance_matrix[node] {
		if !visited[i] && i != node && distance_matrix[node][i] < min {
			min = distance_matrix[node][i]
			nearest = i
		}
	}
	if nearest == -1 {
		err = fmt.Errorf("no nearest node found for %v - invalid distance matrix", node)
	}
	return
}
