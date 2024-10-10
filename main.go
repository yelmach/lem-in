package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"lemin/colony"
	"lemin/pathfinder"
	"lemin/simulator"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("usage: go run . <filename>")
	}
	if filepath.Ext(os.Args[1]) != ".txt" {
		log.Fatalln("usage: go run . <filename>.txt")
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	colony, err := colony.MakeColony(data)
	if err != nil {
		log.Fatalln(err)
	}

	allPaths := colony.FindAllPaths()
	pathfinder.SortPaths(&allPaths)
	fmt.Println(allPaths)

	optimizedPaths := pathfinder.FindOptimalPaths(allPaths, colony.Start.Key, colony.End.Key, colony.AntCount)
	fmt.Println(optimizedPaths)

	// Print ant movements
	simulator.PrintAntMovements(optimizedPaths, colony.AntCount)
}