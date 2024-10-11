package pathfinder

import (
	"lemin/colony"
)

func FindAllPaths(c *colony.Colony) [][]string {
	var allPaths [][]string
	DFS(c.Start, c.End, []string{}, &allPaths)
	return allPaths
}

func DFS(start *colony.Room, end *colony.Room, path []string, allPaths *[][]string) {
	if len(*allPaths) > 10000 {
		return
	}
	start.Visited = true
	path = append(path, start.Key)

	if start == end {
		*allPaths = append(*allPaths, append([]string{}, path...))
	} else {
		for _, next := range start.Tunnels {
			if !next.Visited {
				DFS(next, end, path, allPaths)
			}
		}
	}

	start.Visited = false
}

func SortPaths(paths *[][]string) {
	for i := 0; i < len(*paths); i++ {
		for j := i + 1; j < len(*paths); j++ {
			if len((*paths)[i]) > len((*paths)[j]) {
				(*paths)[i], (*paths)[j] = (*paths)[j], (*paths)[i]
			}
		}
	}
}
