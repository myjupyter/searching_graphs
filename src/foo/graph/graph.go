package graph

import (
	"math"
	"fmt"
	"foo/node"
)

type Graph struct {
	Name  string
	Nodes map[string]*node.Node
}

func (p Graph) GetNode(node_name string) (node.Node, error) {
	value, ok := p.Nodes[node_name]

	if ok {
		return *value, nil
	}

	return node.Node{}, fmt.Errorf("Graph doesn't contain the element %s", node_name)
}

func (p Graph) GetDistanceBtw(v1 string, v2 string) float64 {
	return p.Nodes[v1].GetDistance(v2)
}

func (p Graph) AddBinding(node_name_1, node_name_2 string, dist float64) {
	ptr := new(node.Node)
	*ptr = node.Node{node_name_1, node.Pair{"", math.Inf(1)}, make(map[string]float64, 0)}
	p.PushBack(ptr)
	ptr_2 := new(node.Node)
	*ptr_2 = node.Node{node_name_2, node.Pair{"", math.Inf(1)}, make(map[string]float64, 0)}
	p.PushBack(ptr_2)

	p.Nodes[node_name_1].AddNeighbor(node_name_2, dist)
	p.Nodes[node_name_2].AddNeighbor(node_name_1, dist)
}

func (p Graph) PushBack(nd *node.Node) bool {
	_, ok := p.Nodes[nd.Name]

	if ok {
		return false
	}

	p.Nodes[nd.Name] = nd
	return true
}

func getMinWayVertex(ptr *Graph, visited_vertex map[string]int) string {
	min_key := ""
	min_value := math.Inf(1)
	for key, _ := range ptr.Nodes {
		_, ok := visited_vertex[key]
		if ok {
			continue
		}
		if ptr.Nodes[key].Mark.Length < min_value {
			min_key = key
			min_value = ptr.Nodes[key].Mark.Length
		}
	}
	return min_key
}

func DijkstraSearch(start string, ptr *Graph) error {
	_, ok := ptr.Nodes[start]
	if !ok {
		return fmt.Errorf("Graph doesn't contain the element %s", start)
	}

	ptr.Nodes[start].Mark = node.Pair{"", 0.0}
	vert_count := len(ptr.Nodes)
	visited_vert := make(map[string]int, len(ptr.Nodes))

	for len(visited_vert) < vert_count {
		vertex := getMinWayVertex(ptr, visited_vert)
		visited_vert[vertex] = 1

		for neighbor, distance := range ptr.Nodes[vertex].Neighbors {
			way := ptr.Nodes[vertex].Mark.Length + distance
			if ptr.Nodes[neighbor].Mark.Length > way {
				ptr.Nodes[neighbor].Mark.Length = way

				if ptr.Nodes[vertex].Mark.Way == "" {
					ptr.Nodes[neighbor].Mark.Way = vertex + "|" + neighbor + "|"
				} else {
					ptr.Nodes[neighbor].Mark.Way = ptr.Nodes[vertex].Mark.Way + neighbor + "|"
				}
			}
		}
	}
	return nil
}
