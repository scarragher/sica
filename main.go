package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var _nodes []Node
var lastRefreshTime time.Time

const (
	_configurationFilePath = `.\config.json`
	_cacheDuration         = 30 // cache duration in minutes
)

func main() {

	err := loadConfiguration()

	if err != nil {
		log.Fatalf("\nFailed to read configuration file: %v", err)
	}

	http.HandleFunc("/getStatus", getStatus)
	http.HandleFunc("/addNode", addNode)
	http.HandleFunc("/addChild", addChild)

	log.Fatal(http.ListenAndServe(":1234", nil))
}

func addNode(w http.ResponseWriter, r *http.Request) {

	var node Node

	data, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(data, &node)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mutex := sync.Mutex{}

	mutex.Lock()

	_nodes = append(_nodes, node)

	err = saveConfiguration()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("\nFailed to save configuration: %v", err)
	}

	mutex.Unlock()

	fmt.Println("New node added ", node)
}

func addChild(w http.ResponseWriter, r *http.Request) {

}

func getStatus(w http.ResponseWriter, r *http.Request) {

	if time.Since(lastRefreshTime).Minutes() >= _cacheDuration {
		for i := range _nodes {
			node := &_nodes[i]
			node.updateTree()
		}

		lastRefreshTime = time.Now()
	}

	jsonData, err := json.Marshal(_nodes)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}

func loadConfiguration() error {
	file := ".\\config.json"

	if _, err := os.Stat(file); os.IsNotExist(err) {
		return nil
	}

	data, err := ioutil.ReadFile(file)

	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &_nodes)

	if err != nil {
		return err
	}

	return nil
}

func saveConfiguration() error {
	data, err := json.Marshal(_nodes)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(_configurationFilePath, data, os.ModePerm)

	return nil
}
