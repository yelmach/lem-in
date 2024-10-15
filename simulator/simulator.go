package simulator

import (
	"fmt"
	"lemin/pathfinder"
)

// PrintAntMovements simulates and prints the movement of ants through the colony.
// It takes the optimal paths and the total number of ants as input.
func PrintAntMovements(paths [][]string, antCount int) {
	// Distribute ants among the paths
	antDistribution := pathfinder.DistributeAnts(paths, antCount)

	// AntMove represents a single move of an ant
	type AntMove struct {
		antNumber int
		room      string
	}

	// Initialize tracking variables
	antPosition := make([]int, antCount+1) // Current position of each ant in its path
	antPath := make([]int, antCount+1)     // Which path each ant is following
	currentAnt := 1

	// Assign paths to ants
	for pathIndex, antsInPath := range antDistribution {
		for i := 0; i < antsInPath; i++ {
			antPath[currentAnt] = pathIndex
			currentAnt++
		}
	}

	// Simulate ant movements
	for {
		var moves []AntMove
		allFinished := true
		roomOccupancy := make(map[string]bool)
		endReached := make(map[int]bool)

		// Move each ant if possible
		for ant := 1; ant <= antCount; ant++ {
			pathIndex := antPath[ant]
			if antPosition[ant] < len(paths[pathIndex]) {
				allFinished = false
				nextRoom := paths[pathIndex][antPosition[ant]]

				// Check if the ant can move to the next room
				if !roomOccupancy[nextRoom] || nextRoom == paths[pathIndex][0] || nextRoom == paths[pathIndex][len(paths[pathIndex])-1] {
					if nextRoom == paths[pathIndex][len(paths[pathIndex])-1] {
						// Handle end room
						if !endReached[pathIndex] {
							antPosition[ant]++
							if nextRoom != paths[pathIndex][0] {
								moves = append(moves, AntMove{ant, nextRoom})
							}
						}
						endReached[pathIndex] = true
					} else {
						// Handle regular room
						antPosition[ant]++
						if nextRoom != paths[pathIndex][0] {
							moves = append(moves, AntMove{ant, nextRoom})
						}
					}
				}
				// Mark room as occupied
				if nextRoom != paths[pathIndex][0] && nextRoom != paths[pathIndex][len(paths[pathIndex])-1] {
					roomOccupancy[nextRoom] = true
				}
			}
		}
		// Exit if all ants have reached the end
		if allFinished {
			break
		}

		// Print the moves for this turn
		for _, move := range moves {
			fmt.Printf("L%d-%s ", move.antNumber, move.room)
		}
		fmt.Println()
	}
}
