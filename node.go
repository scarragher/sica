package main

import (
	"net/http"
	"strings"
)

// Node encapsulates information about a specific service
type Node struct {
	Name        string
	Description string
	Address     string
	Children    []Node
	Status      string
	StatusCode  int
}

func (n *Node) addChild(child Node) error {

	n.Children = append(n.Children, child)

	return nil
}

func (n *Node) updateTree() error {

	if n.Address != "" {
		if !strings.HasPrefix(n.Address, "http://") {
			n.Address = "http://" + n.Address
		}

		response, err := http.Get(n.Address)

		if err != nil {
			n.StatusCode = -1
			n.Status = err.Error()
		} else {
			n.Status = response.Status
			n.StatusCode = response.StatusCode
		}
	}

	if len(n.Children) > 0 {
		for i := range n.Children {
			child := &n.Children[i]
			child.updateTree()
		}
	}

	return nil
}
