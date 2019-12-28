package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"foo/graph"
	"foo/parser"
	"strconv"
)

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

	ptr := parser.MakeGraph(words)

	err = graph.DijkstraSearch(os.Args[2], ptr)
	if err != nil {
		log.Fatalln("No such element like '" + os.Args[2] + "'")
	}

	_, ok := ptr.Nodes[os.Args[3]]
	if !ok {
		log.Fatalln("No such element like '" + os.Args[3] + "'")
	}
	vert := strings.Split(ptr.Nodes[os.Args[3]].Mark.Way, "|")

	var out string
	out = "graph " + ptr.Name  +  " {\n"
	for i := 1; i < len(vert) - 1; i++ {
		dist := strconv.FormatFloat(ptr.GetDistanceBtw(vert[i-1], vert[i]), 'f', -1, 64)
		out += "\t" + vert[i-1] + " -- " + vert[i] + " [label=" + dist + "]" +" [color=red];\n"
	}
	outputed := make(map[Pair2]int, 0)
	for i := 1; i < len(vert); i++ {
		outputed[Pair2{vert[i-1], vert[i]}] = 1
		outputed[Pair2{vert[i], vert[i-1]}] = 1
	}

	for _, node := range ptr.Nodes {
		for v, dist := range node.Neighbors {
			_, ok1 := outputed[Pair2{node.Name, v}]
			_, ok2 := outputed[Pair2{v, node.Name}]
			if ok1 {continue}
			if ok2 {continue}
			formated_dist := strconv.FormatFloat(dist, 'f', -1, 64)
			out += "\t" + node.Name  + " -- " + v + " [label=" + formated_dist + "];\n"

			outputed[Pair2{v, node.Name}] = 1
			outputed[Pair2{node.Name, v}] = 1
		}
	}
	out += "}\n"

	file, err := os.Create("./result.dot")
	file.WriteString(out)
	file.Close()
	file.Close()
}
