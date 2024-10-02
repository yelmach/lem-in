package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
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

func (c *Colony) SelectOptimalPaths(paths [][]string) [][]string {
	var optimalPaths [][]string
	bestLength := 0

	for i := 0; i < len(paths); i++ {
		validPaths := [][]string{}
		usedRooms := make(map[string]bool)
		for j := i; j < len(paths); j++ {
			isValid := true
			for _, room := range paths[j][1 : len(paths[j])-1] {
				if usedRooms[room] {
					isValid = false
					break
				}
				usedRooms[room] = true
			}
			if isValid {
				validPaths = append(validPaths, paths[j])
			} else {
				// Reset used rooms for the next path
				for _, room := range paths[j][1 : len(paths[j])-1] {
					usedRooms[room] = false
				}
			}
		}

		antsDistribution := c.DistributeAnts(validPaths)
		maxSteps := antsDistribution[0] + len(validPaths[0]) - 1

		if len(optimalPaths) == 0 || maxSteps < bestLength {
			optimalPaths = validPaths
			bestLength = maxSteps
		}
	}

	return optimalPaths
}

func (c *Colony) DistributeAnts(paths [][]string) []int {
	antsDistribution := make([]int, len(paths))
	remainingAnts := c.AntCount
	round := 0

	for remainingAnts > 0 {
		for i := range paths {
			if i+1 > len(paths) || len(paths[i])+round >= len(paths[i+1]) {
				antsDistribution[i]++
				remainingAnts--
				if remainingAnts == 0 {
					break
				}
			}
		}
		round++
	}

	return antsDistribution
}

func (c *Colony) PrintAntMovement(paths [][]string, antsDistribution []int) {
	type AntMove struct {
		antNumber int
		room      string
	}

	antPosition := make([]int, c.AntCount+1)
	antPath := make([]int, c.AntCount+1)
	currentAnt := 1

	for pathIndex, antsInPath := range antsDistribution {
		for i := 0; i < antsInPath; i++ {
			antPath[currentAnt] = pathIndex
			currentAnt++
		}
	}

	for {
		var moves []AntMove
		allFinished := true
		antPresent := make(map[string]bool)

		for i := 1; i <= c.AntCount; i++ {
			pathIndex := antPath[i]

			if antPosition[i] < len(paths[pathIndex]) {
				allFinished = false
				nextRoom := paths[pathIndex][antPosition[i]]

				if !antPresent[nextRoom] || nextRoom == paths[pathIndex][0] || nextRoom == paths[pathIndex][len(paths[pathIndex])-1] {
					antPosition[i]++
					if nextRoom != paths[pathIndex][0] {
						moves = append(moves, AntMove{
							antNumber: i,
							room:      nextRoom,
						})
					}
				}

				if nextRoom != paths[pathIndex][0] && nextRoom != paths[pathIndex][len(paths[pathIndex])-1] {
					antPresent[nextRoom] = true
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run . <filename>")
		return
	}
	colony, err := MakeGraph(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	allPaths := colony.FindAllPaths()
	SortPaths(&allPaths)
	fmt.Println(len(allPaths))
	os.Exit(0)
	optimalPaths := colony.SelectOptimalPaths(allPaths)
	antDistribution := colony.DistributeAnts(optimalPaths)

	// Print the input data
	data, _ := os.ReadFile(os.Args[1])
	fmt.Println(string(data))

	// Print ant movements
	colony.PrintAntMovement(optimalPaths, antDistribution)
}

func MakeGraph(filename string) (*Colony, error) {
	colony := NewColony()
	if filepath.Ext(filename) != ".txt" {
		log.Fatalln("Usage: go run . <filename>.txt")
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	var startFlag, endFlag bool

	for i := 0; i < len(lines); i++ {
		lines[i] = strings.TrimSpace(lines[i])
		parts := strings.Fields(lines[i])
		switch {
		case i == 0:
			antCount, err := strconv.Atoi(lines[i])
			if err != nil || antCount <= 0 {
				return nil, fmt.Errorf("ERROR: invalid number of Ants")
			}
			colony.SetAntCount(antCount)
		case lines[i] == "":
			continue
		case strings.HasPrefix(lines[i], "#") && lines[i] != "##start" && lines[i] != "##end":
			continue
		case lines[i] == "##start":
			if startFlag {
				return nil, fmt.Errorf("ERROR: 2 start points found")
			}
			startFlag = true
			i++
			if strings.HasPrefix(lines[i], "#") || strings.HasPrefix(lines[0], "L") {
				return nil, fmt.Errorf("ERROR: Invalid room name")
			}
			colony.SetStart(strings.Split(lines[i], " ")[0])
		case lines[i] == "##end":
			if endFlag {
				return nil, fmt.Errorf("ERROR: 2 end points found")
			}
			endFlag = true
			i++
			if strings.HasPrefix(lines[i], "#") || strings.HasPrefix(lines[0], "L") {
				return nil, fmt.Errorf("ERROR: Invalid room name")
			}
			colony.SetEnd(strings.Split(lines[i], " ")[0])
		case len(parts) == 3:
			if strings.HasPrefix(parts[0], "L") {
				return nil, fmt.Errorf("ERROR: Invalid room name")
			}
			_, exist := colony.AddRoom(parts[0])
			if exist {
				return nil, fmt.Errorf("ERROR: room already Exist")
			}
			x, err1 := strconv.Atoi(parts[1])
			y, err2 := strconv.Atoi(parts[2])
			if err1 != nil || err2 != nil || x < 0 || y < 0 {
				return nil, fmt.Errorf("ERROR: invalid coordinates for room '%s'", parts[0])
			}
			if check := strings.Fields(lines[i-1]); len(check) != 3 && !strings.HasPrefix(lines[i-1], "#") {
				return nil, fmt.Errorf("ERROR: order disrespected")
			}
		case len(parts) == 1:
			if check := strings.Split(lines[i], "-"); len(check) == 2 {
				colony.AddTunnel(check[0], check[1])
			} else {
				return nil, fmt.Errorf("ERROR: invalid tunnel")
			}
		}
	}
	if !startFlag || !endFlag {
		return nil, fmt.Errorf("ERROR: Start or End room not found")
	}
	return colony, nil
}

func NewColony() *Colony {
	return &Colony{
		Rooms: make(map[string]*Room),
	}
}

func (c *Colony) SetAntCount(count int) {
	c.AntCount = count
}

func (c *Colony) AddRoom(key string) (*Room, bool) {
	if _, exists := c.Rooms[key]; !exists {
		c.Rooms[key] = &Room{Key: key}
		return c.Rooms[key], false
	}
	return c.Rooms[key], true
}

func (c *Colony) AddTunnel(key1, key2 string) {
	if key1 == key2 {
		log.Fatalln("Adding tunnel Failed: give different rooms")
	}
	from, exist := c.AddRoom(key1)
	to, exist1 := c.AddRoom(key2)
	if !exist || !exist1 {
		log.Fatalln("Adding tunnel Failed: rooms Does Not Exist")
	}
	if !containsRoom(from.Tunnels, to) {
		from.Tunnels = append(from.Tunnels, to)
		to.Tunnels = append(to.Tunnels, from)
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
	start, exist := c.AddRoom(key)
	if exist {
		log.Fatalln("Adding Start Room Failed")
	}
	c.Start = start
}

func (c *Colony) SetEnd(key string) {
	end, exist := c.AddRoom(key)
	if exist {
		log.Fatalln("Adding end Room Failed")
	}
	c.End = end
}

func (c *Colony) FindAllPaths() [][]string {
	var allPaths [][]string
	c.DFS(c.Start, c.End, []string{}, &allPaths)
	return allPaths
}

func (c *Colony) DFS(start *Room, end *Room, path []string, allPaths *[][]string) {
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
				c.DFS(next, end, path, allPaths)
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
