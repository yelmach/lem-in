package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseInput(lines []string) (int, []string, []string, error) {
	var numberOfAnts int
	var rooms []string
	var tunnels []string
	state := "ants"

	for i := 0; i < len(lines); i++ {
		lines[i] = strings.TrimSpace(lines[i])

		if lines[i] == "" || (strings.HasPrefix(lines[i], "#") && lines[i] != "##start" && lines[i] != "##end") {
			continue
		}

		switch state {
		case "ants":
			var err error
			numberOfAnts, err = strconv.Atoi(lines[i])
			if err != nil || numberOfAnts <= 0 {
				return 0, nil, nil, fmt.Errorf("ERROR: invalid number of Ants")
			}
			state = "rooms"
		case "rooms":
			// Detect the start of tunnels
			if strings.Contains(lines[i], "-") {
				state = "tunnels"
				i--
				continue
			}
			rooms = append(rooms, lines[i])
		case "tunnels":
			// Ensure correct format for tunnels
			if !strings.Contains(lines[i], "-") {
				return 0, nil, nil, fmt.Errorf("ERROR: invalid data format")
			}
			tunnels = append(tunnels, lines[i])
		}
	}
	return numberOfAnts, rooms, tunnels, nil
}
