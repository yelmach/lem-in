package main

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

// AddRoom adds a new room to the colony
func (c *Colony) AddRoom(key string) *Room {
	if _, exists := c.Rooms[key]; !exists {
		c.Rooms[key] = &Room{Key: key}
	}
	return c.Rooms[key]
}

// AddTunnel connects two rooms with a tunnel
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