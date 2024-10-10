package colony

import (
	"fmt"
	"log"
	"strings"

	"lemin/utils"
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
