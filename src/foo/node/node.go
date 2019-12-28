package node

import "math"

type Pair struct {
	Way    string
	Length float64
}

type Node struct {
	Name      string
	Mark      Pair
	Neighbors map[string]float64
}

func (n Node) AddNeighbor(name string, dist float64) bool {
	_, ok := n.Neighbors[name]

	if ok {
		return false
	}

	n.Neighbors[name] = dist
	return true
}

func (n Node) GetDistance(vertex string) float64 {
	_, ok := n.Neighbors[vertex]
	if !ok {
		return math.Inf(1)
	}
	return n.Neighbors[vertex]
}

func (n Node) GetWay() string {
	return n.Mark.Way
}
