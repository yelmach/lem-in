package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"src/utils"
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

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("usage: go run . <filename>")
	}
	if filepath.Ext(os.Args[1]) != ".txt" {
		log.Fatalln("usage: go run . <filename>.txt")
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	colony, err := MakeColony(data)
	if err != nil {
		log.Fatalln(err)
	}

	allPaths := colony.FindAllPaths()
	SortPaths(&allPaths)
	fmt.Println(allPaths)

	optimizedPaths := cleanPaths(allPaths, colony.Start.Key, colony.End.Key, colony.AntCount)
	fmt.Println(optimizedPaths)

	// Print ant movements
	PrintAntMovements(optimizedPaths, colony.AntCount)
}

// checking file and adding room to graph and links between rooms
func MakeColony(data []byte) (*Colony, error) {
	colony := NewColony()
	lines := strings.Split(string(data), "\n")

	antCount, rooms, tunnels, err := utils.ParseInput(lines)
	if err != nil {
		return nil, err
	}
	colony.SetAntCount(antCount)

	// get name of rooms
	coordMap := make(map[[2]int]bool)
	var startFlag, endFlag bool
	for i := 0; i < len(rooms); i++ {
		switch {
		case rooms[i] == "##start":
			if startFlag {
				return nil, fmt.Errorf("ERROR: 2 start points found")
			}
			startFlag = true
			if i+1 < len(rooms) {
				i++
				err := utils.CheckingCoord(rooms[i], coordMap)
				if err != nil {
					return nil, err
				}
				colony.SetStart(strings.Split(rooms[i], " ")[0])
			} else {
				return nil, fmt.Errorf("ERROR: No start point found")
			}
		case rooms[i] == "##end":
			if endFlag {
				return nil, fmt.Errorf("ERROR: 2 end points found")
			}
			endFlag = true
			if i+1 < len(rooms) {
				i++
				err := utils.CheckingCoord(rooms[i], coordMap)
				if err != nil {
					return nil, err
				}
				colony.SetEnd(strings.Split(rooms[i], " ")[0])
			} else {
				return nil, fmt.Errorf("ERROR: No end point found")
			}
		case len(strings.Fields(rooms[i])) == 3:
			err := utils.CheckingCoord(rooms[i], coordMap)
			if err != nil {
				return nil, err
			}
			// rooms
			_, exist := colony.AddRoom(strings.Split(rooms[i], " ")[0])
			if exist {
				return nil, fmt.Errorf("ERROR: room already Exist")
			}
		default:
			return nil, fmt.Errorf("ERROR: invalid data format")
		}
	}
	if !startFlag || !endFlag {
		return nil, fmt.Errorf("ERROR: Start or End room not found")
	}

	// add links between rooms
	for i := 0; i < len(tunnels); i++ {
		if check := strings.Split(tunnels[i], "-"); len(check) == 2 {
			colony.AddTunnel(check[0], check[1])
		} else {
			return nil, fmt.Errorf("ERROR: invalid tunnel")
		}
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
		log.Fatalf("ERROR: Adding link Failed, you give same room: '%s'", key1)
	}
	from, exist := c.AddRoom(key1)
	to, exist1 := c.AddRoom(key2)
	if !exist || !exist1 {
		log.Fatalln("Adding link Failed: rooms Does Not Exist")
	}
	if !containsRoom(from.Tunnels, to) {
		from.Tunnels = append(from.Tunnels, to)
		to.Tunnels = append(to.Tunnels, from)
	} else {
		log.Fatalln("Adding link Failed: link already Exist")
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

// end

func (c *Colony) FindAllPaths() [][]string {
	var allPaths [][]string
	DFS(c.Start, c.End, []string{}, &allPaths)
	return allPaths
}

func DFS(start *Room, end *Room, path []string, allPaths *[][]string) {
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
// ----------------------------------------------------------------

/*
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
			antDistribution := distributeAnts(validPaths, antCount)
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

func distributeAnts(paths [][]string, antCount int) []int {
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

func PrintAntMovements(paths [][]string, antCount int) {
	antDistribution := distributeAnts(paths, antCount)

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
}*/


func cleanPaths(pathes [][]string, start, end string, naml int) [][]string {
	finalpath := [][]string{}
	l := 0
	for i := 0; i < len(pathes); i++ {
		validPaths := [][]string{}
		nodeMap := make(map[string]int)
		for j := i; j < len(pathes); j++ {
			doz := true
			s := 0
			for _, node := range pathes[j] {
				if node != end && node != start {
					nodeMap[node]++ //=nodemap[A]==1  nodemap[B]==2 nodemap[C]==1
					s++
				}
				if nodeMap[node] > 1 {
					doz = false
					for _, node := range pathes[j] {
						s--
						nodeMap[node]--
						if s == 0 {
							break
						}
					}
					break
				}
			}
			if doz {
				validPaths = append(validPaths, pathes[j])
			}
		}
		if comparisonpath(validPaths, finalpath, naml, &l) {
			finalpath = validPaths
		}
	}
	return finalpath
}

func comparisonpath(validPaths, finalpath [][]string, naml int, l *int) bool {
	if finalpath == nil {
		return true
	}
	if len(validPaths) == 0 {
		return false
	}
	arypaths := make([]int, len(validPaths))
	arypaths = HowManyAntsInEachPath(validPaths, naml)
	if *l == 0 || arypaths[0]+len(validPaths[0]) < *l {
		*l = arypaths[0] + len(validPaths[0])
		return true
	}
	return false
}

func HowManyAntsInEachPath(pathes [][]string, naml int) []int {
	arypaths := make([]int, len(pathes))
	doz := 0
	k := 0
	if len(pathes) == 0 {
		os.Exit(1)
	}
	for naml > 0 {
		if len(pathes) > doz+1 && len(pathes[doz])+k >= len(pathes[doz+1]) {
			k = 0
			doz++
			continue
		}
		for i := 0; i <= doz; i++ {
			arypaths[i]++
			naml--
			if naml == 0 {
				break
			}
			k++
		}
	}
	return arypaths
}

func PrintAntMovements(paths [][]string, numAnts int) {
	// Get the number of ants assigned to each path
	antsInEachPath := HowManyAntsInEachPath(paths, numAnts)

	// Struct to represent an ant's movement with its number and the room it's moving to
	type AntMove struct {
		antNumber int
		room      string
	}

	// Arrays to track each ant's current position in their path and which path they are on
	antPosition := make([]int, numAnts+1)
	antPath := make([]int, numAnts+1)
	currentAnt := 1

	// Distribute ants across the paths based on the number of ants assigned to each path
	for pathIndex, antsInPath := range antsInEachPath { //[0,0,0]
		for i := 0; i < antsInPath; i++ {
			antPath[currentAnt] = pathIndex // Assign the ant to a path
			currentAnt++
		}
	}

	for {
		var moves []AntMove
		allFinished := true
		antPresent := make(map[string]bool)
		finishedInPath := make(map[int]bool)

		for i := 1; i <= numAnts; i++ {
			pathIndex := antPath[i]

			// Check if the ant hasn't reached the end of its path
			if antPosition[i] < len(paths[pathIndex]) {
				allFinished = false
				nextRoom := paths[pathIndex][antPosition[i]]

				if !antPresent[nextRoom] || nextRoom == paths[pathIndex][0] || nextRoom == paths[pathIndex][len(paths[pathIndex])-1] {
					if nextRoom == paths[pathIndex][len(paths[pathIndex])-1] {
						if !finishedInPath[pathIndex] {
							antPosition[i]++

							if nextRoom != paths[pathIndex][0] {
								moves = append(moves, AntMove{
									antNumber: i,
									room:      nextRoom,
								})
							}
						}
						finishedInPath[pathIndex] = true
					} else {

						antPosition[i]++
						if nextRoom != paths[pathIndex][0] {
							moves = append(moves, AntMove{
								antNumber: i,
								room:      nextRoom,
							})
						}
					}
				}
				// If the ant is in an intermediate room, mark it as present
				if nextRoom != paths[pathIndex][0] && nextRoom != paths[pathIndex][len(paths[pathIndex])-1] {
					antPresent[nextRoom] = true
				}
			}
		}

		// If all ants have finished moving, exit the loop
		if allFinished {
			break
		}

		// Print all the movements of ants in this round
		for _, move := range moves {
			fmt.Printf("L%d-%s ", move.antNumber, move.room)
		}
		fmt.Println()
	}
}
