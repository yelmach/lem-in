package simulator

import (
	"fmt"
	"strings"

	"lemin/pathfinder"
)

type AntMove struct {
	antNumber int
	room      string
}

func PrintAntMovements(paths [][]string, numAnts int) {
	antsInEachPath := pathfinder.DistributeAnts(paths, numAnts)
	antPositions := make([]int, numAnts+1)
	antPaths := make([]int, numAnts+1)
	currentAnt := 1

	for pathIndex, antsInPath := range antsInEachPath {
		for i := 0; i < antsInPath; i++ {
			antPaths[currentAnt] = pathIndex
			currentAnt++
		}
	}

	var result strings.Builder
	for {
		moves := make([]AntMove, 0, numAnts)
		allFinished := true
		antPresent := make(map[string]bool)

		for ant := 1; ant <= numAnts; ant++ {
			pathIndex := antPaths[ant]
			path := paths[pathIndex]

			if antPositions[ant] < len(path) {
				allFinished = false
				nextRoom := path[antPositions[ant]]

				if canMoveToRoom(nextRoom, path, antPresent) {
					antPositions[ant]++
					if nextRoom != path[0] {
						moves = append(moves, AntMove{ant, nextRoom})
						if nextRoom != path[len(path)-1] {
							antPresent[nextRoom] = true
						}
					}
				}
			}
		}

		if allFinished {
			break
		}

		printMoves(&result, moves)
		fmt.Println(result.String())
		result.Reset()
	}
}

func canMoveToRoom(room string, path []string, antPresent map[string]bool) bool {
	return !antPresent[room] || room == path[0] || room == path[len(path)-1]
}

func printMoves(result *strings.Builder, moves []AntMove) {
	for _, move := range moves {
		result.WriteString(fmt.Sprintf("L%d-%s ", move.antNumber, move.room))
	}
}
