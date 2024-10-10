package pathfinder

import "log"

func SortPaths(paths *[][]string) {
	for i := 0; i < len(*paths); i++ {
		for j := i + 1; j < len(*paths); j++ {
			if len((*paths)[i]) > len((*paths)[j]) {
				(*paths)[i], (*paths)[j] = (*paths)[j], (*paths)[i]
			}
		}
	}
}

func FindOptimalPaths(allPaths [][]string, start, end string, antCount int) [][]string {
	bestPaths := [][]string{}
	bestLength := 0

	for i := 0; i < len(allPaths); i++ {
		currentPaths := [][]string{}
		usedRooms := make(map[string]bool)

		for _, path := range allPaths[i:] {
			if isValidPath(path, start, end, usedRooms) {
				currentPaths = append(currentPaths, path)
				markRoomsAsUsed(path, start, end, usedRooms)
			}
		}

		if isNewPathSetBetter(currentPaths, bestPaths, antCount, &bestLength) {
			bestPaths = currentPaths
		}
	}

	return bestPaths
}

func isValidPath(path []string, start, end string, usedRooms map[string]bool) bool {
	for _, room := range path[1 : len(path)-1] {
		if room != start && room != end && usedRooms[room] {
			return false
		}
	}
	return true
}

func markRoomsAsUsed(path []string, start, end string, usedRooms map[string]bool) {
	for _, room := range path[1 : len(path)-1] {
		if room != start && room != end {
			usedRooms[room] = true
		}
	}
}

func isNewPathSetBetter(validPaths, currentBestPaths [][]string, ants int, currentBestLength *int) bool {
	if currentBestPaths == nil {
		return true
	}
	if len(validPaths) == 0 {
		return false
	}
	antsPerPath := DistributeAnts(validPaths, ants)
	if *currentBestLength == 0 || antsPerPath[0]+len(validPaths[0]) < *currentBestLength {
		*currentBestLength = antsPerPath[0] + len(validPaths[0])
		return true
	}
	return false
}

func DistributeAnts(paths [][]string, antCount int) []int {
	if len(paths) == 0 {
		log.Fatal("No paths available for ant distribution")
	}

	antDistribution := make([]int, len(paths))
	pathLengths := make([]int, len(paths))
	for i, path := range paths {
		pathLengths[i] = len(path) - 1 // Subtract 1 to not count the start room
	}

	for ant := 0; ant < antCount; ant++ {
		shortestPath := 0
		shortestTime := pathLengths[0] + antDistribution[0]

		for i := 1; i < len(paths); i++ {
			currentTime := pathLengths[i] + antDistribution[i]
			if currentTime < shortestTime {
				shortestPath = i
				shortestTime = currentTime
			}
		}

		antDistribution[shortestPath]++
	}

	return antDistribution
}
