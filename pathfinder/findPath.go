package pathfinder

import (
	"lemin/colony"
)

// FindAllPaths finds all possible paths from the start room to the end room in the colony.
// It returns a slice of slices string, where each inner slice represents a path.
func FindAllPaths(c *colony.Colony) [][]string {
	var allPaths [][]string
	DFS(c.Start, c.End, []string{}, &allPaths)
	return allPaths
}

// DFS performs a depth-first search to find all paths from start to end.
// It uses recursion to explore all possible routes.
func DFS(start *colony.Room, end *colony.Room, path []string, allPaths *[][]string) {
	// Limit the number of paths to prevent excessive memory usage
	if len(*allPaths) > 10000 {
		return
	}
	start.Visited = true
	path = append(path, start.Key)

	if start == end {
		// If we've reached the end, add the current path to allPaths
		*allPaths = append(*allPaths, append([]string{}, path...))
	} else {
		// else Continue searching through connected rooms
		for _, next := range start.Tunnels {
			if !next.Visited {
				DFS(next, end, path, allPaths)
			}
		}
	}

	// Backtrack: mark the room as unvisited when leaving
	start.Visited = false
}

// SortPaths sorts the paths by length in ascending order.
// It uses a simple bubble sort algorithm.
func SortPaths(paths *[][]string) {
	for i := 0; i < len(*paths); i++ {
		for j := i + 1; j < len(*paths); j++ {
			if len((*paths)[i]) > len((*paths)[j]) {
				(*paths)[i], (*paths)[j] = (*paths)[j], (*paths)[i]
			}
		}
	}
}
