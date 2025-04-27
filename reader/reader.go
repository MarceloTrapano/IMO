package reader

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Node struct {
	X int
	Y int
}

// wczytywanie instancji z pliku
// instancje TSPLIB kroA200, kroB200 z https://github.com/mastqe/tsplib
// oryginalne źródło: http://comopt.ifi.uni-heidelberg.de/software/TSPLIB95/tsp/
func ReadInstance(srcPath string) (nodes []Node, headers map[string]string, err error) {
	var num_nodes int                 // liczba wierzchołków
	headers = make(map[string]string) // nagłówki z pliku

	file, err := os.Open(srcPath)
	if err != nil {
		return
	}
	defer file.Close()

	// stworzenie skanera
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	// wczytywanie nagłówków
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		// rozpoczęcie sekcji z wierzchołkami
		if line == "NODE_COORD_SECTION" {
			break
		}
		// rozdzielenie nagłówka od wartości na dwukropku i usunięcie białych znaków
		fields := strings.Split(line, ":")
		if len(fields) < 2 {
			err = fmt.Errorf("invalid instance file format (reading %d-th header)", i+1)
			return
		}
		for field, _ := range fields {
			fields[field] = strings.TrimSpace(fields[field])
		}
		var h, val *string = &fields[0], &fields[1] // ćwiczonko wskaźników
		headers[*h] = *val
		// alternatywnie na wartościach
		// h, val := fields[0], fields[1]
		// headers[h] = val
	}
	// odczytanie liczby wierzchołków
	dim, ok := headers["DIMENSION"]
	num_nodes, err = strconv.Atoi(dim)
	// błąd konwersji lub brak nagłówka DIMENSION
	if err != nil || !ok {
		return
	}
	// wczytywanie wierzchołków
	for i := 0; scanner.Scan() && i < num_nodes; i++ {
		var (
			x_str, y_str string
			x, y         int
		)
		fields := strings.Fields(scanner.Text())
		if len(fields) < 3 {
			err = fmt.Errorf("invalid instance file format (reading %d-th node)", i+1)
			return
		}
		_, x_str, y_str = fields[0], fields[1], fields[2]
		x, err = strconv.Atoi(x_str)
		if err != nil {
			return
		}
		y, err = strconv.Atoi(y_str)
		if err != nil {
			return
		}
		nodes = append(nodes, Node{x, y})
	}

	return
}
