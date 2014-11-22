package main

import (
	_ "fmt"
	"log"

	"code.google.com/p/gcfg"
)

// Struct to hold the indexes and queries once they have been loaded in
type cfg struct {
	Index map[string]*struct {
		Query    []string
		QueryDef []string
	}
}

// Load indexes and queries from cfg file
func Load() cfg {
	//Create a new cfg struct
	var c cfg
	//Read the index.hbit file
	err := gcfg.ReadFileInto(&c, "index.hbit")
	if err != nil {
		log.Fatalf("Failed to parse gcfg data: %s", err)
	}
	return c
}

// Get all indices
func getIndices(c *cfg) []string {
	//The return array
	r := []string{}
	for k, _ := range c.Index {
		r = append(r, k)
	}
	return r
}
