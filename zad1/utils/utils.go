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

func (edge *Edge) String() string {
	var buf bytes.Buffer
	var first int = edge.From
	buf.WriteString("{")
	for edge != nil {
		buf.WriteString(fmt.Sprintf("(%v -> %v): %v", edge.From, edge.To, edge.Length))
		if edge.Next != nil && edge.To != first {
			buf.WriteString(", ")
		}
		if edge.To == first {
			break
		}
		edge = edge.Next
	}

	buf.WriteString("}")
	return buf.String()
}

func NewEdge(from int, to int, distance_matrix [][]int, prev *Edge, next *Edge) *Edge {
	return &Edge{
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

func EdgeToNodeCycle(edge *Edge) []int {
	var (
		cycle []int
		first int = edge.From
	)
	for {
		cycle = append(cycle, edge.From)
		edge = edge.Next
		if edge.From == first {
			break
		}
	}
	return cycle
}

func UpdateDistances(eLL *EdgeLinkedList, distance_matrix [][]int, delEdges []int, newEdges []EdgeLinkedList, newEdgesSorted bool) *EdgeLinkedList {
	var remainingDelete int = len(delEdges)
	if !newEdgesSorted {
		// sortuje rosnąco - chcemy malejąco (najlepsze na końcu) więc przeciwnie: j-i zamiast i-j
		slices.SortFunc(newEdges, func(i, j EdgeLinkedList) int {
			return j.Value - i.Value
		})
	}
	for eLL != nil && (len(newEdges) > 0 || remainingDelete > 0) {
		// usuń krawędzie
		for i, delEdge := range delEdges {
			if eLL.Edge == delEdge {
				remainingDelete--
				switch {
				case eLL.Prev == nil && eLL.Next == nil && len(newEdges) > 0: // jedyny element
					eLL = &newEdges[len(newEdges)-1]
					newEdges = newEdges[:len(newEdges)-1]
				case eLL.Prev == nil: // pierwszy element
					eLL = eLL.Next
					eLL.Prev = nil
				case eLL.Next == nil && len(newEdges) > 0: // ostatni element
					eLL.Prev.Next = &newEdges[len(newEdges)-1]
					newEdges = newEdges[:len(newEdges)-1]
					eLL = eLL.Prev.Next
				case eLL.Next != nil && eLL.Prev != nil: //  w środku
					eLL.Prev.Next = eLL.Next
					eLL.Next.Prev = eLL.Prev
					eLL = eLL.Next
				case eLL.Next == nil && eLL.Prev != nil: // przedostatni element
					eLL.Prev.Next = nil
					eLL = eLL.Prev
				default: // brak krawędzi do dodania w zamian
					return nil
				}
				delEdges[i] = -1
			}
		}

		if len(newEdges) > 0 {
			last := len(newEdges) - 1
			edge := &newEdges[last]
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
				newEdges = newEdges[:last]
			case eLL.Next == nil: // koniec listy - dodaj na końcu; && edge.Value >= eLL.Value pominięte
				edge.Prev = eLL
				eLL.Next = edge
				newEdges = newEdges[:last]
			default: // gorszy niż aktulany - szukaj dalej; edge.Value >= eLL.Value pominięte
				eLL = eLL.Next
			}
		} else if eLL != nil && eLL.Next != nil {
			eLL = eLL.Next
		} else {
			break
		}
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
func MaxOfArray(arr []int) (int, int, error) {
	idx := -1
	max := -1
	for i, value := range arr {
		if value > max {
			max = value
			idx = i
		}
	}
	return max, idx, nil
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

func IndexBefore(arr []int, i int) int {
	if i == 0 {
		return len(arr) - 1
	}
	return i - 1
}

func IndexAfter(arr []int, i int) int {
	if i == len(arr)-1 {
		return 0
	}
	return i + 1
}

func ElemBefore(arr []int, i int) int {
	if i == 0 {
		return arr[len(arr)-1]
	}
	return arr[i-1]
}

func ElemAfter(arr []int, i int) int {
	if i == len(arr)-1 {
		return arr[0]
	}
	return arr[i+1]
}
