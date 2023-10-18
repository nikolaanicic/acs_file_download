package parser

import "strings"

type Node struct {
	Name  string
	Url   string
	Nodes []*Node
}

func newNode(name string, url string) *Node {
	return &Node{
		Name:  name,
		Url:   url,
		Nodes: make([]*Node, 0),
	}
}

func (n *Node) AddNodes(nodes []*Node) {
	n.Nodes = append(n.Nodes, nodes...)
}

func (n *Node) IsFolder() bool {
	return !strings.Contains(n.Name, ".")
}
