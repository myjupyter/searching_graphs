package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type Pair struct {
	way    string
	length float64
}

type Node struct {
	name      string
	mark      Pair
	neighbors map[string]float64
}

func (n Node) AddNeighbor(name string, dist float64) bool {
	_, ok := n.neighbors[name]

	if ok {
		return false
	}

	n.neighbors[name] = dist
	return true
}

func (n Node) GetDistance(vertex string) float64 {
	_, ok := n.neighbors[vertex]
	if !ok {
		return math.Inf(1)
	}
	return n.neighbors[vertex]
}

type Graph struct {
	name  string
	nodes map[string]*Node
}

func (p Graph) GetNode(node_name string) (Node, error) {
	value, ok := p.nodes[node_name]

	if ok {
		return *value, nil
	}

	return Node{}, fmt.Errorf("Graph doesn't contain the element %s", node_name)
}

func (p Graph) GetDistanceBtw(v1 string, v2 string) float64 {
	return p.nodes[v1].GetDistance(v2)
}

func (p Graph) AddBinding(node_name_1, node_name_2 string, dist float64) {
	ptr := new(Node)
	*ptr = Node{node_name_1, Pair{"", math.Inf(1)}, make(map[string]float64, 0)}
	p.PushBack(ptr)
	ptr_2 := new(Node)
	*ptr_2 = Node{node_name_2, Pair{"", math.Inf(1)}, make(map[string]float64, 0)}
	p.PushBack(ptr_2)

	p.nodes[node_name_1].AddNeighbor(node_name_2, dist)
	p.nodes[node_name_2].AddNeighbor(node_name_1, dist)
}

func (p Graph) PushBack(node *Node) bool {
	_, ok := p.nodes[node.name]

	if ok {
		return false
	}

	p.nodes[node.name] = node
	return true
}

var GraphType = map[string]int{
	"digraph": 1,
	"graph":   2,
}

var Binding = map[string]int{
	"->": 1,
	"--": 2,
}

var Parameters = map[string]int{
	"label": 1,
}

func ParseLabel(word string) float64 {
	var s, e int
	for i, char := range word {
		if string(char) == "=" {
			s = i + 1
		}
		if string(char) == "]" {
			e = i
		}
	}
	dist, _ := strconv.ParseFloat(word[s:e], 64)
	return dist
}

func MakeGraph(words []string) *Graph {

	mygraph := new(Graph)
	*mygraph = Graph{
		name:  "",
		nodes: make(map[string]*Node, 0),
	}

	for i, word := range words {
		_, ok := GraphType[word]
		if ok {
			mygraph.name = words[i+1]
		}
		_, ok = Binding[word]
		if ok {
			cleared_word := strings.ReplaceAll(words[i+1], ";", "")
			mygraph.AddBinding(words[i-1], cleared_word, ParseLabel(words[i+2]))
		}
	}
	return mygraph
}

func getMinWayVertex(ptr *Graph, visited_vertex map[string]int) string {
	min_key := ""
	min_value := math.Inf(1)
	for key, _ := range ptr.nodes {
		_, ok := visited_vertex[key]
		if ok {
			continue
		}
		if ptr.nodes[key].mark.length < min_value {
			min_key = key
			min_value = ptr.nodes[key].mark.length
		}
	}
	return min_key
}

func DijkstraSearch(start string, ptr *Graph) error {
	_, ok := ptr.nodes[start]
	if !ok {
		return fmt.Errorf("Graph doesn't contain the element %s", start)
	}

	ptr.nodes[start].mark = Pair{"", 0.0}
	vert_count := len(ptr.nodes)
	visited_vert := make(map[string]int, len(ptr.nodes))

	for len(visited_vert) < vert_count {
		vertex := getMinWayVertex(ptr, visited_vert)
		visited_vert[vertex] = 1

		for neighbor, distance := range ptr.nodes[vertex].neighbors {
			way := ptr.nodes[vertex].mark.length + distance
			if ptr.nodes[neighbor].mark.length > way {
				ptr.nodes[neighbor].mark.length = way

				if ptr.nodes[vertex].mark.way == "" {
					ptr.nodes[neighbor].mark.way = vertex + "|" + neighbor + "|"
				} else {
					ptr.nodes[neighbor].mark.way = ptr.nodes[vertex].mark.way + neighbor + "|"
				}
			}
		}
	}
	return nil
}

type Pair2 struct {
	v1 string
	v2 string
}

func main() {
	if len(os.Args) <  4 {
		log.Fatalln("Not enough parametrs!")
	}

	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalln("Some problem with the file:", err)
	}

	text := strings.ReplaceAll(string(data), "\n", " ")
	text = strings.ReplaceAll(text, "\t", " ")
	words := strings.Split(text, " ")

	ptr := MakeGraph(words)

	DijkstraSearch(os.Args[2], ptr)


	vert := strings.Split(ptr.nodes[os.Args[3]].mark.way, "|")

	var out string
	out = "graph " + ptr.name  +  " {\n"
	for i := 1; i < len(vert) - 1; i++ {
		dist := strconv.FormatFloat(ptr.GetDistanceBtw(vert[i-1], vert[i]), 'f', -1, 64)
		out += "\t" + vert[i-1] + " -- " + vert[i] + " [label=" + dist + "]" +" [color=red];\n"
	}
	outputed := make(map[Pair2]int, 0)
	for i := 1; i < len(vert); i++ {
		outputed[Pair2{vert[i-1], vert[i]}] = 1
		outputed[Pair2{vert[i], vert[i-1]}] = 1
	}

	for _, node := range ptr.nodes {
		for v, dist := range node.neighbors {
			_, ok1 := outputed[Pair2{node.name, v}]
			_, ok2 := outputed[Pair2{v, node.name}]
			if ok1 {continue}
			if ok2 {continue}
			formated_dist := strconv.FormatFloat(dist, 'f', -1, 64)
			out += "\t" + node.name  + " -- " + v + " [label=" + formated_dist + "];\n"

			outputed[Pair2{v, node.name}] = 1
			outputed[Pair2{node.name, v}] = 1
		}
	}
	out += "}\n"

	file, err := os.Create("./result.dot")
	file.WriteString(out)
	file.Close()
	file.Close()}
