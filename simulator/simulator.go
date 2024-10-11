package simulator

import (
	"fmt"
	"lemin/pathfinder"
)

func PrintAntMovements(paths [][]string, antCount int) {
	antDistribution := pathfinder.DistributeAnts(paths, antCount)

	type AntMove struct {
		antNumber int
		room      string
	}

	antPosition := make([]int, antCount+1)
	antPath := make([]int, antCount+1)
	currentAnt := 1

	for pathIndex, antsInPath := range antDistribution {
		for i := 0; i < antsInPath; i++ {
			antPath[currentAnt] = pathIndex
			currentAnt++
		}
	}

	for {
		var moves []AntMove
		allFinished := true
		roomOccupancy := make(map[string]bool)
		endReached := make(map[int]bool)

		for ant := 1; ant <= antCount; ant++ {
			pathIndex := antPath[ant]
			if antPosition[ant] < len(paths[pathIndex]) {
				allFinished = false
				nextRoom := paths[pathIndex][antPosition[ant]]

				if !roomOccupancy[nextRoom] || nextRoom == paths[pathIndex][0] || nextRoom == paths[pathIndex][len(paths[pathIndex])-1] {
					if nextRoom == paths[pathIndex][len(paths[pathIndex])-1] {
						if !endReached[pathIndex] {
							antPosition[ant]++
							if nextRoom != paths[pathIndex][0] {
								moves = append(moves, AntMove{ant, nextRoom})
							}
						}
						endReached[pathIndex] = true
					} else {
						antPosition[ant]++
						if nextRoom != paths[pathIndex][0] {
							moves = append(moves, AntMove{ant, nextRoom})
						}
					}
				}

				if nextRoom != paths[pathIndex][0] && nextRoom != paths[pathIndex][len(paths[pathIndex])-1] {
					roomOccupancy[nextRoom] = true
				}
			}
		}

		if allFinished {
			break
		}

		for _, move := range moves {
			fmt.Printf("L%d-%s ", move.antNumber, move.room)
		}
		fmt.Println()
	}
}
