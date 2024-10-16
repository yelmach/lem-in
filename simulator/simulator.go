package simulator

import (
	"fmt"
	"strings"
)

func PrintAntMovements(paths [][]string, numAnts int) {
	antPaths := make([][]int, len(paths))
	for ant := 1; ant <= numAnts; ant++ {
		shortestPath := findShortestPath(paths, antPaths)
		antPaths[shortestPath] = append(antPaths[shortestPath], ant)
	}

	moves := generateMoves(antPaths, paths)
	for _, moveSet := range moves {
		if len(moveSet) > 0 {
			fmt.Println(strings.Join(moveSet, " "))
		}
	}
}

func findShortestPath(paths [][]string, antPaths [][]int) int {
	shortestIndex, shortestLength := 0, len(paths[0])+len(antPaths[0])
	for i := 1; i < len(paths); i++ {
		currentLength := len(paths[i]) + len(antPaths[i])
		if currentLength < shortestLength {
			shortestIndex, shortestLength = i, currentLength
		}
	}
	return shortestIndex
}

func generateMoves(antPaths [][]int, paths [][]string) [][]string {
	var moves [][]string
	for pathIndex, ants := range antPaths {
		path := paths[pathIndex][1:] // Skip the start room
		for stepIndex, antID := range ants {
			for roomIndex, room := range path {
				moveIndex := roomIndex + stepIndex
				move := fmt.Sprintf("L%d-%s", antID, room)
				if moveIndex >= len(moves) {
					moves = append(moves, []string{move})
				} else {
					moves[moveIndex] = append(moves[moveIndex], move)
				}
			}
		}
	}
	return moves
}


