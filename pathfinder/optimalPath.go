package pathfinder

// FindOptimalPaths finds the optimal set of paths for the ants to take.
// It takes all possible paths, start and end room names, and the number of ants as input.
// It returns the best combination of paths that minimizes the number of turns needed.
func FindOptimalPaths(allPaths [][]string, start, end string, antCount int) [][]string {
	bestPaths := [][]string{}
	bestLength := 0

	// Iterate through all paths to find the optimal combination
	for i := 0; i < len(allPaths); i++ {
		currentPaths := [][]string{}
		usedRooms := make(map[string]bool)

		// Check each path from the current index
		for _, path := range allPaths[i:] {
			if isValidPath(path, start, end, usedRooms) {
				currentPaths = append(currentPaths, path)
				markRoomsAsUsed(path, start, end, usedRooms)
			}
		}

		// Update bestPaths if the current combination is better
		if isNewPathSetBetter(currentPaths, bestPaths, antCount, &bestLength) {
			bestPaths = currentPaths
		}
	}

	return bestPaths
}

// isValidPath checks if a path is valid (doesn't use already used rooms).
// It returns true if the path is valid, false otherwise.
func isValidPath(path []string, start, end string, usedRooms map[string]bool) bool {
	for _, room := range path[1 : len(path)-1] {
		if room != start && room != end && usedRooms[room] {
			return false
		}
	}
	return true
}

// markRoomsAsUsed marks all rooms in a path as used, except for start and end rooms.
func markRoomsAsUsed(path []string, start, end string, usedRooms map[string]bool) {
	for _, room := range path[1 : len(path)-1] {
		if room != start && room != end {
			usedRooms[room] = true
		}
	}
}

// isNewPathSetBetter checks if the new set of paths is better than the current best.
// It updates the currentBestLength if the new set is better.
// Returns true if the new set is better, false otherwise.
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

// DistributeAnts distributes the ants among the available paths.
// It aims to minimize the total time by sending ants through shorter paths first.
// Returns a slice with the number of ants assigned to each path.
func DistributeAnts(paths [][]string, antCount int) []int {
	antDistribution := make([]int, len(paths))
	pathLengths := make([]int, len(paths))
	for i, path := range paths {
		pathLengths[i] = len(path) - 1 // Subtract 1 to not count the start room
	}

	// Distribute ants one by one
	for ant := 0; ant < antCount; ant++ {
		shortestPath := 0
		shortestTime := pathLengths[0] + antDistribution[0]

		// Find the path that will result in the shortest time for the current ant
		for i := 1; i < len(paths); i++ {
			currentTime := pathLengths[i] + antDistribution[i]
			if currentTime < shortestTime {
				shortestPath = i
				shortestTime = currentTime
			}
		}

		// Assign the ant to the shortest path
		antDistribution[shortestPath]++
	}

	return antDistribution
}
