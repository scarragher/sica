package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestNodes(t *testing.T) {

	parent := Node{
		Name:       "A parent",
		Status:     "200 - OK",
		StatusCode: 200,
		Children: []Node{
			Node{
				Name:       "Child A",
				Status:     "200 - OK",
				StatusCode: 200,
			},
			Node{
				Name:       "Child B",
				Status:     "404 - Not Found",
				StatusCode: 404,
			},
		},
	}

	fmt.Printf("\nParent: %s, %s\n", parent.Name, parent.Status)

	for _, child := range parent.Children {
		fmt.Printf("\nChild: %s, %s\n", child.Name, child.Status)
	}

	jsonData, err := json.Marshal(parent)

	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(jsonData))

	var unmarshaledNode Node

	err = json.Unmarshal(jsonData, &unmarshaledNode)

	if err != nil {
		t.Error(err)
	}
}
