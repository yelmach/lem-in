package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func CheckingCoord(line string, coordMap map[[2]int]bool) error {
	parts := strings.Fields(line)
	if len(parts) != 3 {
		return fmt.Errorf("ERROR: Invalid room format")
	}
	if strings.HasPrefix(line, "L") {
		return fmt.Errorf("ERROR: Invalid room name")
	}
	x, errX := strconv.Atoi(parts[1])
	y, errY := strconv.Atoi(parts[2])
	if errX != nil || errY != nil || x < 0 || y < 0 {
		return fmt.Errorf("ERROR: invalid coordinates for room '%s'", parts[0])
	}

	coord := [2]int{x, y}
	if coordMap[coord] {
		return fmt.Errorf("ERROR: coordinates (%d, %d) are duplicated", x, y)
	}
	coordMap[coord] = true
	return nil
}
