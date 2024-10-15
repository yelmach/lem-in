package simulator

import (
	"fmt"

	"lemin/pathfinder"
)

// PrintAntMovements simulates and prints the movement of ants through the colony.
// It takes the optimal paths and the total number of ants as input.
// func PrintAntMovements(paths [][]string, antCount int) {
//     // Distribute ants among the paths using a function from the pathfinder package
//     antDistribution := pathfinder.DistributeAnts(paths, antCount)

//     // Define a struct to represent a single move of an ant
//     type AntMove struct {
//         antNumber int    // The number of the ant making the move
//         room      string // The room the ant is moving to
//     }

//     // Initialize slice to track the current position of each ant in its path
//     antPosition := make([]int, antCount+1)
//     // Initialize slice to track which path each ant is following
//     antPath := make([]int, antCount+1)
//     // Initialize a counter for assigning ants to paths
//     currentAnt := 1

//     // Assign paths to ants based on the calculated distribution
//     for pathIndex, antsInPath := range antDistribution {
//         for i := 0; i < antsInPath; i++ {
//             antPath[currentAnt] = pathIndex
//             currentAnt++
//         }
//     }

//     // Start the main simulation loop
//     for {
//         // Initialize a slice to store moves for this turn
//         var moves []AntMove
//         // Assume all ants have finished, will be set to false if any ant moves
//         allFinished := true
//         // Map to track which rooms are occupied this turn
//         roomOccupancy := make(map[string]bool)
//         // Map to track which paths have reached their end
//         endReached := make(map[int]bool)

//         // Iterate over all ants to determine their moves
//         for ant := 1; ant <= antCount; ant++ {
//             // Get the path index for this ant
//             pathIndex := antPath[ant]
//             // Check if the ant hasn't reached the end of its path
//             if antPosition[ant] < len(paths[pathIndex]) {
//                 // At least one ant hasn't finished, so set allFinished to false
//                 allFinished = false
//                 // Get the next room for this ant
//                 nextRoom := paths[pathIndex][antPosition[ant]]

//                 // Check if the ant can move to the next room
//                 if !roomOccupancy[nextRoom] || nextRoom == paths[pathIndex][0] || nextRoom == paths[pathIndex][len(paths[pathIndex])-1] {
//                     if nextRoom == paths[pathIndex][len(paths[pathIndex])-1] {
//                         // Handle end room
//                         if !endReached[pathIndex] {
//                             // Move the ant forward
//                             antPosition[ant]++
//                             // Add the move to the list if it's not the start room
//                             if nextRoom != paths[pathIndex][0] {
//                                 moves = append(moves, AntMove{ant, nextRoom})
//                             }
//                         }
//                         // Mark this path as having reached the end
//                         endReached[pathIndex] = true
//                     } else {
//                         // Handle regular room
//                         // Move the ant forward
//                         antPosition[ant]++
//                         // Add the move to the list if it's not the start room
//                         if nextRoom != paths[pathIndex][0] {
//                             moves = append(moves, AntMove{ant, nextRoom})
//                         }
//                     }
//                 }
//                 // Mark room as occupied if it's not the start or end room
//                 if nextRoom != paths[pathIndex][0] && nextRoom != paths[pathIndex][len(paths[pathIndex])-1] {
//                     roomOccupancy[nextRoom] = true
//                 }
//             }
//         }
//         // Exit the loop if all ants have reached the end
//         if allFinished {
//             break
//         }

//         // Print the moves for this turn
//         for _, move := range moves {
//             fmt.Printf("L%d-%s ", move.antNumber, move.room)
//         }
//         // Print a newline to separate turns
//         fmt.Println()
//     }
// }

// PrintAntMovements simulates and prints the movement of ants through the colony.
func PrintAntMovements(paths [][]string, antCount int) {
	antDistribution := pathfinder.DistributeAnts(paths, antCount)
	antPositions, antPaths := initializeAnts(antCount, antDistribution)

	for !allAntsFinished(antPositions, paths, antPaths) {
		moves := simulateTurn(antPositions, antPaths, paths)
		printMoves(moves)
	}
}

// AntMove represents a single move of an ant
type AntMove struct {
	antNumber int
	room      string
}

// initializeAnts sets up the initial state for all ants
func initializeAnts(antCount int, antDistribution []int) ([]int, []int) {
	antPositions := make([]int, antCount+1)
	antPaths := make([]int, antCount+1)
	currentAnt := 1

	for pathIndex, antsInPath := range antDistribution {
		for i := 0; i < antsInPath; i++ {
			antPaths[currentAnt] = pathIndex
			currentAnt++
		}
	}

	return antPositions, antPaths
}

// allAntsFinished checks if all ants have reached the end of their paths
func allAntsFinished(antPositions []int, paths [][]string, antPaths []int) bool {
	for ant := 1; ant < len(antPositions); ant++ {
		if antPositions[ant] < len(paths[antPaths[ant]])-1 {
			return false
		}
	}
	return true
}

// simulateTurn simulates one turn of ant movements
func simulateTurn(antPositions []int, antPaths []int, paths [][]string) []AntMove {
	var moves []AntMove
	roomOccupancy := make(map[string]bool)
	endReached := make(map[int]bool)

	for ant := 1; ant < len(antPositions); ant++ {
		move := tryMoveAnt(ant, antPositions, antPaths, paths, roomOccupancy, endReached)
		if move != nil {
			moves = append(moves, *move)
		}
	}

	return moves
}

// tryMoveAnt attempts to move a single ant
func tryMoveAnt(ant int, antPositions []int, antPaths []int, paths [][]string, roomOccupancy map[string]bool, endReached map[int]bool) *AntMove {
	pathIndex := antPaths[ant]
	currentPath := paths[pathIndex]

	if antPositions[ant] >= len(currentPath)-1 {
		return nil // Ant has already finished
	}

	nextRoom := currentPath[antPositions[ant]+1]

	if canMoveToRoom(nextRoom, currentPath, roomOccupancy) {
		return moveAnt(ant, nextRoom, antPositions, pathIndex, currentPath, roomOccupancy, endReached)
	}

	return nil
}

// canMoveToRoom checks if an ant can move to the next room
func canMoveToRoom(nextRoom string, path []string, roomOccupancy map[string]bool) bool {
	return !roomOccupancy[nextRoom] || nextRoom == path[0] || nextRoom == path[len(path)-1]
}

// moveAnt moves an ant to the next room and returns the move
func moveAnt(ant int, nextRoom string, antPositions []int, pathIndex int, path []string, roomOccupancy map[string]bool, endReached map[int]bool) *AntMove {
	antPositions[ant]++

	if nextRoom == path[len(path)-1] {
		if !endReached[pathIndex] {
			endReached[pathIndex] = true
			return &AntMove{ant, nextRoom}
		}
	} else if nextRoom != path[0] {
		roomOccupancy[nextRoom] = true
		return &AntMove{ant, nextRoom}
	}

	return nil
}

// printMoves prints all moves for a single turn
func printMoves(moves []AntMove) {
	for _, move := range moves {
		fmt.Printf("L%d-%s ", move.antNumber, move.room)
	}
	fmt.Println()
}