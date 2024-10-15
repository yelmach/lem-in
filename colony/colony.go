package colony

import (
	"fmt"
	"log"
	"strings"

	"lemin/utils"
)

// Room represents a single room in the ant colony.
type Room struct {
	Key     string  // name of the room
	Tunnels []*Room // Connections to other rooms
	Visited bool    // Used in path finding
}

// Colony represents the entire ant colony structure.
type Colony struct {
	Rooms    map[string]*Room // All rooms in the colony
	Start    *Room            // Starting room
	End      *Room            // Ending room
	AntCount int              // Total number of ants
}

// MakeColony creates a new Colony structure from the provided input data.
// It parses the input, creates rooms, and make links(tunnels) between rooms.
func MakeColony(data []byte) (*Colony, error) {
	colony := NewColony()
	lines := strings.Split(string(data), "\n")

	// split the input into 3 parts: ant count, rooms, and tunnels
	antCount, rooms, tunnels, err := utils.ParseInput(lines)
	if err != nil {
		return nil, err
	}
	colony.SetAntCount(antCount)

	// Map for checking coordinates
	coordMap := make(map[[2]int]bool)
	var startFlag, endFlag bool

	// Process each room in the input
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
				return nil, fmt.Errorf("ERROR: no start room found")
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
				return nil, fmt.Errorf("ERROR: No end room found")
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

	// add tunnels (connection) between rooms
	for i := 0; i < len(tunnels); i++ {
		if check := strings.Split(tunnels[i], "-"); len(check) == 2 {
			colony.AddTunnel(check[0], check[1])
		} else {
			return nil, fmt.Errorf("ERROR: invalid tunnel")
		}
	}
	return colony, nil
}

// NewColony creates and returns a new empty Colony structure.
func NewColony() *Colony {
	return &Colony{
		Rooms: make(map[string]*Room),
	}
}

// SetAntCount sets the total number of ants in the colony.
func (c *Colony) SetAntCount(count int) {
	c.AntCount = count
}

// AddRoom adds a new room to the colony or returns an existing one.
// It returns the room and a boolean indicating if the room already existed.
func (c *Colony) AddRoom(key string) (*Room, bool) {
	if _, exists := c.Rooms[key]; !exists {
		c.Rooms[key] = &Room{Key: key}
		return c.Rooms[key], false
	}
	return c.Rooms[key], true
}

// AddTunnel creates a connection between two rooms.
func (c *Colony) AddTunnel(key1, key2 string) {
	if key1 == key2 {
		log.Fatalf("ERROR: Adding link Failed, you give the same room: '%s'", key1)
	}
	from, exist := c.AddRoom(key1)
	to, exist1 := c.AddRoom(key2)
	if !exist || !exist1 {
		log.Fatalln("Adding link Failed: room Does Not Exist")
	}
	if !containsRoom(from.Tunnels, to) {
		from.Tunnels = append(from.Tunnels, to)
		to.Tunnels = append(to.Tunnels, from)
	} else {
		log.Fatalln("Adding link Failed: link already Exist")
	}
}

// containsRoom checks if a room is already in a slice of rooms.
func containsRoom(rooms []*Room, room *Room) bool {
	for _, r := range rooms {
		if r == room {
			return true
		}
	}
	return false
}

// SetStart sets the starting room for the colony.
func (c *Colony) SetStart(key string) {
	start, exist := c.AddRoom(key)
	if exist {
		log.Fatalln("ERROR: Adding Start Room Failed")
	}
	c.Start = start
}

// SetEnd sets the ending room for the colony.
func (c *Colony) SetEnd(key string) {
	end, exist := c.AddRoom(key)
	if exist {
		log.Fatalln("ERROR: Adding end Room Failed")
	}
	c.End = end
}
