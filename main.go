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

	// Read the contents of the input file
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	// Create the ant colony based on the input data
	colony, err := colony.MakeColony(data)
	if err != nil {
		log.Fatalln(err)
	}

	// Find all possible paths through the colony
	allPaths := pathfinder.FindAllPaths(colony)
	if len(allPaths) == 0 {
		log.Fatalln("no path found")
	}

	// Sort the paths by length
	pathfinder.SortPaths(&allPaths)

	// Find the optimal set of paths for the ants to take
	bestPaths := pathfinder.FindOptimalPaths(allPaths, colony.Start.Key, colony.End.Key, colony.AntCount)

	// Print the original input data
	fmt.Println(string(data))
	fmt.Println()

	// Simulate and print the ant movements
	simulator.PrintAntMovements(bestPaths, colony.AntCount)
}
