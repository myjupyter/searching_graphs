package parser

import (
	"strconv"
	"strings"
	"foo/node"
	"foo/graph"
)

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

func MakeGraph(words []string) *graph.Graph {

	mygraph := new(graph.Graph)
	*mygraph = graph.Graph{
		Name:  "",
		Nodes: make(map[string]*node.Node, 0),
	}

	for i, word := range words {
		_, ok := GraphType[word]
		if ok {
			mygraph.Name = words[i+1]
		}
		_, ok = Binding[word]
		if ok {
			cleared_word := strings.ReplaceAll(words[i+1], ";", "")
			mygraph.AddBinding(words[i-1], cleared_word, ParseLabel(words[i+2]))
		}
	}
	return mygraph
}
