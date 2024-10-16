package simulator

import "fmt"

// PrintAntMovements simulates the movement of ants through the best paths.
// It takes the best paths and the total number of ants as input, calculates
// the optimal distribution of ants across paths, and then prints the movements.
func PrintAntMovements(bestPath [][]string, antCount int) {
	// pathCosts holds the cost (length) of each path minus the start and end rooms.
	pathCosts := make([]int, len(bestPath))
	// antsOnPath keeps track of how many ants are assigned to each path.
	antsOnPath := make([]int, len(bestPath))

	// Initialize path costs and assign at least one ant to each path.
	for i, path := range bestPath {
		pathCosts[i] = len(path) - 2 // Exclude start and end rooms
		antsOnPath[i] = 1
		antCount-- // Deduct one ant as it's assigned to this path.
	}

	// Distribute the remaining ants to paths, preferring the shorter ones.
	for antCount > 0 {
		minIndex := 0
		// Find the path that currently offers the least total cost (path cost + ants).
		for i := 1; i < len(bestPath); i++ {
			if pathCosts[i]+antsOnPath[i] < pathCosts[minIndex]+antsOnPath[minIndex] {
				minIndex = i
			}
		}
		// Assign the ant to the path with the minimum cost.
		antsOnPath[minIndex]++
		antCount--
	}

	// Print the movements of ants based on the distribution across paths.
	printMoves(bestPath, antsOnPath)
}

// printMoves simulates and prints the actual movement of ants through the paths.
// It takes the paths and the number of ants assigned to each path as input.
// The function continues to move ants and print their positions until all ants
// have reached the end of their respective paths.
func printMoves(paths [][]string, antsPerPath []int) {
	// antPositions keeps track of the current positions of ants on each path.
	antPositions := make([][]string, len(paths))
	for i, path := range paths {
		antPositions[i] = make([]string, len(path))
	}

	antID := 1 // Track the ID of the current ant being moved.
	for {
		// Move ants one step forward along their respective paths.
		move(antPositions)

		// Place new ants at the start of their paths as long as there are remaining ants.
		for pathIndex, remainingAnts := range antsPerPath {
			if remainingAnts > 0 {
				antPositions[pathIndex][1] = fmt.Sprintf("L%d", antID) // Place the ant in the first room after start
				antID++
				antsPerPath[pathIndex]-- // Decrease the count of remaining ants for this path.
			}
		}

		// Print the current positions of all ants.
		printCurrentPositions(antPositions, paths)

		// If all ants have reached the end of their paths, exit the loop.
		if allAntsReachedEnd(antPositions) {
			break
		}
		fmt.Println() // Print a newline to separate the moves.
	}
}

// printCurrentPositions prints the current positions of ants in their respective rooms.
func printCurrentPositions(antPositions [][]string, paths [][]string) {
	// Iterate through each path and print the position of ants in reverse order.
	for pathIndex, path := range antPositions {
		for roomIndex := len(path) - 1; roomIndex > 0; roomIndex-- {
			if path[roomIndex] != "" {
				// Print the ant ID along with the room it is in.
				fmt.Printf("%s-%s ", path[roomIndex], paths[pathIndex][roomIndex])
			}
		}
	}
}

// allAntsReachedEnd checks if all ants have reached the end of their paths.
// It returns true if all ants have reached the end, otherwise returns false.
func allAntsReachedEnd(antPositions [][]string) bool {
	// Check each path to see if any ant is still in transit (i.e., not at the end).
	for _, path := range antPositions {
		for roomIndex := 1; roomIndex < len(path); roomIndex++ {
			if path[roomIndex] != "" {
				return false // Ant found still in transit, so not all ants have reached the end.
			}
		}
	}
	return true // All ants have reached the end of their paths.
}

// move simulates the movement of ants along their paths.
// It shifts each ant one step forward along its path,
// emptying the start room of each path in the process.
func move(antPositions [][]string) {
	// Move ants along each path by shifting them forward by one room.
	for pathIndex, path := range antPositions {
		for roomIndex := len(path) - 1; roomIndex > 0; roomIndex-- {
			// Shift the ant from the previous room to the current room.
			antPositions[pathIndex][roomIndex] = antPositions[pathIndex][roomIndex-1]
		}
		antPositions[pathIndex][0] = "" // Clear the start room after moving the ants.
	}
}
