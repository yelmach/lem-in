package pathfinder

import "fmt"

func FindDisjointPaths(allpaths [][]string, start, end string, antCount int) [][]string {
	bestPaths := [][]string{}
	bestTotalTime := 0

	for i := 0; i < len(allpaths); i++ {
		validPaths := [][]string{}
		usedRooms := make(map[string]bool)

		for j := i; j < len(allpaths); j++ {
			isValid := true
			for _, room := range allpaths[j][1 : len(allpaths[j])-1] {
				if usedRooms[room] {
					isValid = false
					break
				}
				usedRooms[room] = true
			}

			if isValid {
				validPaths = append(validPaths, allpaths[j])
			}
		}
		fmt.Println(validPaths)

		if len(validPaths) > 0 {
			antDistribution := DistributeAnts(validPaths, antCount)
			totalTime := 0
			for i, count := range antDistribution {
				if count > 0 {
					pathTime := len(validPaths[i]) - 1 + count
					if pathTime > totalTime {
						totalTime = pathTime
					}
				}
			}

			if bestTotalTime == 0 || totalTime < bestTotalTime {
				bestPaths = validPaths
				bestTotalTime = totalTime
			}
		}
	}

	return bestPaths
}

func DistributeAnts(paths [][]string, antCount int) []int {
	distribution := make([]int, len(paths))
	pathLengths := make([]int, len(paths))
	for i, path := range paths {
		pathLengths[i] = len(path) - 1 // Subtract 1 because we don't count the start room
	}

	for ant := 1; ant <= antCount; ant++ {
		bestPath := 0
		bestArrivalTime := pathLengths[0] + distribution[0] + 1

		for i := 1; i < len(paths); i++ {
			arrivalTime := pathLengths[i] + distribution[i] + 1
			if arrivalTime < bestArrivalTime {
				bestPath = i
				bestArrivalTime = arrivalTime
			}
		}

		distribution[bestPath]++
	}

	return distribution
}
