package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Room struct {
	Key     string
	Tunnels []*Room
	Visited bool
}

type Colony struct {
	Rooms    map[string]*Room
	Start    *Room
	End      *Room
	AntCount int
}

func NewColony() *Colony {
	return &Colony{
		Rooms: make(map[string]*Room),
	}
}

func (c *Colony) AddRoom(key string) *Room {
	if _, exists := c.Rooms[key]; !exists {
		c.Rooms[key] = &Room{Key: key}
	}
	return c.Rooms[key]
}

func (c *Colony) AddTunnel(key1, key2 string) {
	room1 := c.AddRoom(key1)
	room2 := c.AddRoom(key2)
	if !containsRoom(room1.Tunnels, room2) {
		room1.Tunnels = append(room1.Tunnels, room2)
		room2.Tunnels = append(room2.Tunnels, room1)
	}
}

func containsRoom(rooms []*Room, room *Room) bool {
	for _, r := range rooms {
		if r == room {
			return true
		}
	}
	return false
}

func (c *Colony) SetStart(key string) {
	c.Start = c.AddRoom(key)
}

func (c *Colony) SetEnd(key string) {
	c.End = c.AddRoom(key)
}

func (c *Colony) SetAntCount(count int) {
	c.AntCount = count
}

func (c *Colony) ResetVisited() {
	for _, room := range c.Rooms {
		room.Visited = false
	}
}

func (c *Colony) DFS(start *Room, end *Room, path []string, allPaths *[][]string) {
	start.Visited = true
	path = append(path, start.Key)

	if start == end {
		pathCopy := make([]string, len(path))
		copy(pathCopy, path)
		*allPaths = append(*allPaths, pathCopy)
	} else {
		for _, neighbor := range start.Tunnels {
			if !neighbor.Visited {
				c.DFS(neighbor, end, path, allPaths)
			}
		}
	}

	start.Visited = false
	path = path[:len(path)-1]
}

func (c *Colony) FindAllPaths() [][]string {
	var allPaths [][]string
	c.ResetVisited()
	c.DFS(c.Start, c.End, []string{}, &allPaths)
	return allPaths
}

func SortPaths(paths [][]string) {
	sort.Slice(paths, func(i, j int) bool {
		return len(paths[i]) < len(paths[j])
	})
}

func (c *Colony) SelectOptimalPaths(sortedPaths [][]string) [][]string {
	var optimalPaths [][]string
	usedRooms := make(map[string]bool)

	for _, path := range sortedPaths {
		isValid := true
		for _, room := range path[1 : len(path)-1] {
			if usedRooms[room] {
				isValid = false
				break
			}
		}
		if isValid {
			optimalPaths = append(optimalPaths, path)
			for _, room := range path[1 : len(path)-1] {
				usedRooms[room] = true
			}
		}
		if len(optimalPaths) >= c.AntCount {
			break
		}
	}
	return optimalPaths
}

func (c *Colony) DistributeAnts(paths [][]string) []int {
	antDistribution := make([]int, len(paths))
	remainingAnts := c.AntCount

	for remainingAnts > 0 {
		for i := range paths {
			if remainingAnts > 0 {
				antDistribution[i]++
				remainingAnts--
			} else {
				break
			}
		}
	}
	return antDistribution
}

func (c *Colony) SimulateAntMovement(paths [][]string, antDistribution []int) {
	antPositions := make([]int, c.AntCount)
	antPaths := make([]int, c.AntCount)

	currentAnt := 0
	for pathIndex, antCount := range antDistribution {
		for i := 0; i < antCount; i++ {
			antPaths[currentAnt] = pathIndex
			currentAnt++
		}
	}

	finished := false
	turn := 0
	for !finished {
		finished = true
		movements := []string{}

		for ant := 0; ant < c.AntCount; ant++ {
			if antPositions[ant] < len(paths[antPaths[ant]])-1 {
				finished = false
				antPositions[ant]++
				movements = append(movements, fmt.Sprintf("L%d-%s", ant+1, paths[antPaths[ant]][antPositions[ant]]))
			}
		}

		if len(movements) > 0 {
			fmt.Printf("Turn %d: %s\n", turn+1, strings.Join(movements, " "))
		}
		turn++
	}
}

func ProcessInputFile(filename string) (*Colony, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	colony := NewColony()
	var startSet, endSet bool

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if i == 0 {
			antCount, err := strconv.Atoi(line)
			if err != nil {
				return nil, fmt.Errorf("invalid ant count: %s", line)
			}
			colony.SetAntCount(antCount)
		} else if line == "##start" {
			i++
			startRoom := strings.Split(lines[i], " ")[0]
			colony.SetStart(startRoom)
			startSet = true
		} else if line == "##end" {
			i++
			endRoom := strings.Split(lines[i], " ")[0]
			colony.SetEnd(endRoom)
			endSet = true
		} else if strings.Contains(line, "-") {
			rooms := strings.Split(line, "-")
			if len(rooms) == 2 {
				colony.AddTunnel(rooms[0], rooms[1])
			}
		}
	}

	if !startSet || !endSet {
		return nil, fmt.Errorf("start or end room not set")
	}

	return colony, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run . <filename>")
		return
	}

	colony, err := ProcessInputFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error processing input file: %v\n", err)
		return
	}

	allPaths := colony.FindAllPaths()
	SortPaths(allPaths)
	optimalPaths := colony.SelectOptimalPaths(allPaths)
	antDistribution := colony.DistributeAnts(optimalPaths)

	fmt.Printf("Number of ants: %d\n", colony.AntCount)
	fmt.Println("Rooms:")
	for key := range colony.Rooms {
		fmt.Println(key)
	}
	fmt.Println("Tunnels:")
	for _, room := range colony.Rooms {
		for _, tunnel := range room.Tunnels {
			fmt.Printf("%s-%s\n", room.Key, tunnel.Key)
		}
	}
	fmt.Println("\nAnt movements:")
	colony.SimulateAntMovement(optimalPaths, antDistribution)
}
