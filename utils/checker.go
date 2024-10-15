package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// CheckingCoord validates the coordinates of a room and checks for duplicates.
// It takes a line from the input and a map of existing coordinates.
// Returns an error if the coordinates are invalid or duplicated.
func CheckingCoord(line string, coordMap map[[2]int]bool) error {
	parts := strings.Fields(line)
	if len(parts) != 3 {
		return fmt.Errorf("ERROR: Invalid room format")
	}
	if strings.HasPrefix(line, "L") {
		return fmt.Errorf("ERROR: Invalid room name")
	}

	// Parse and validate X and Y coordinates
	x, errX := strconv.Atoi(parts[1])
	y, errY := strconv.Atoi(parts[2])
	if errX != nil || errY != nil || x < 0 || y < 0 {
		return fmt.Errorf("ERROR: invalid coordinates for room '%s'", parts[0])
	}

	// Check for duplicate coordinates
	coord := [2]int{x, y}
	if coordMap[coord] {
		return fmt.Errorf("ERROR: coordinates (%d, %d) are duplicated", x, y)
	}
	coordMap[coord] = true
	return nil
}
