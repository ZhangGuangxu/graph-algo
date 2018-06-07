package main

import (
	"fmt"
	"github.com/json-iterator/go"
	"io/ioutil"
)

type edge struct {
	Cost float32 `json:"v"`
}

type edges map[int]edge

type edgesMap map[int]edges

type mapData struct {
	edgesMap edgesMap
}

var p = fmt.Println

func (d *mapData) load(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		p(err)
		return err
	}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(content, &d.edgesMap)
	if err != nil {
		p(err)
		return err
	}

	return nil
}

func (d *mapData) show() {
	p(d.edgesMap)
}

func (d *mapData) store() {

}
