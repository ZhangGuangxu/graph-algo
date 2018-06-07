package main

import (
	"testing"
)

// TestMapdata test loading a.map
func TestMapdata(t *testing.T) {
	md := mapData{}
	err := md.load("../../bin/a.map")
	if err != nil {
		t.Errorf("load a.map got error %v", err)
	}
}
