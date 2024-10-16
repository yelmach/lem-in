package simulator

import (
	"fmt"
	"strings"
)

// PrintAntMovements simulates and prints the movement of ants through the colony.
// It takes the optimal paths and total number of ants as input, then prints each
// turn's movements in the format "L<ant_number>-<room_name>".
func PrintAntMovements(bestPaths [][]string, totalAnts int) {
	// Initialize a slice to store which ants are assigned to each path
	antPath := make([][]int, len(bestPaths))

	// Assign each ant to the path that will result in the shortest completion time
	for ant := 1; ant <= totalAnts; ant++ {
		bestPathIndex := findShortestPathForAnt(bestPaths, antPath)
		antPath[bestPathIndex] = append(antPath[bestPathIndex], ant)
	}

	// Generate and print all movements for each turn
	moves := simulateAntMovements(antPath, bestPaths)
	for _, turn := range moves {
		if len(turn) > 0 {
			fmt.Println(strings.Join(turn, " "))
		}
	}
}

// findShortestPathForAnt determines which path would be best for the next ant.
// It calculates the completion time for each path considering its length and
// how many ants are already assigned to it.
// and will return The index of the path that will result in the shortest completion time
func findShortestPathForAnt(paths [][]string, antPath [][]int) int {
	bestPathIndex := 0
	minTime := len(paths[0])+len(antPath[0])
	// Check each path's completion time
	for pathIndex := 1; pathIndex < len(paths); pathIndex++ {
		time := len(paths[pathIndex]) + len(antPath[pathIndex])
		if time < minTime {
			bestPathIndex = pathIndex
			minTime = time
		}
	}
	return bestPathIndex
}

// simulateAntMovements generates all ant movements through their assigned paths.
// It creates a timeline of movements where each ant follows its path while
// avoiding collisions with other ants.
// Returns: A slice where each element represents one turn of ant movements
func simulateAntMovements(antPath [][]int, paths [][]string) [][]string {
	var movements [][]string

	// Process movements for each path
	for pathIndex, antsOnPath := range antPath {
		// removing start room
		path := paths[pathIndex][1:]

		// Generate moves for each ant on the current path
		for antOrder, ant := range antsOnPath {
			// Move ant through each room in the path
			for roomIndex, roomName := range path {
				// Calculate when this movement should occur
				timeStep := roomIndex + antOrder
				move := fmt.Sprintf("L%d-%s", ant, roomName)

				// Add the movement to the appropriate time step
				if timeStep >= len(movements) {
					movements = append(movements, []string{move})
				} else {
					movements[timeStep] = append(movements[timeStep], move)
				}
			}
		}
	}
	return movements
}