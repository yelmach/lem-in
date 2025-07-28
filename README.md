# lem-in - Digital Ant Farm Pathfinding

A Go implementation of an ant colony simulation that finds the optimal paths for ants to traverse from a start room to an end room through a network of tunnels and rooms.

## Overview

**lem-in** reads a colony description from a file and simulates ants moving through tunnels between rooms. The program uses Depth-First Search (DFS) algorithms to find the quickest way to get all ants from the `##start` room to the `##end` room with the minimum number of moves.

## Features

- **Pathfinding Algorithm**: Uses Depth-First Search (DFS) to explore possible paths
- **Multiple Path Support**: Finds and utilizes multiple non-conflicting paths when beneficial
- **Traffic Management**: Prevents ants from colliding by ensuring only one ant per room (except start/end)
- **Tunnel Usage Optimization**: Each tunnel can only be used once per turn
- **Error Handling**: Comprehensive validation of input data format
- **Standard Go Packages Only**: No external dependencies

## How It Works

1. **Parse Input**: Read colony file containing number of ants, rooms, and tunnels
2. **Build Graph**: Create internal representation of the colony network
3. **Find Paths**: Use DFS to discover all possible paths from start to end
4. **Optimize Routes**: Select the combination of paths that minimizes total moves
5. **Simulate Movement**: Move ants turn by turn, displaying each move

### Algorithm Details

The program uses **Depth-First Search (DFS)** for pathfinding:
- Explores paths recursively from the start room
- Backtracks when reaching dead ends or visited nodes
- Finds multiple independent paths to optimize ant distribution
- Selects the best combination of paths based on total move count

## Input Format

The input file should contain:

1. **Number of ants** (first line)
2. **Rooms** in format: `name coord_x coord_y`
   - `##start` marker before the start room
   - `##end` marker before the end room
3. **Tunnels** in format: `room1-room2`
4. **Comments** starting with `#` (ignored)

### Sample Input

```
3
##start
1 23 3
2 16 7
3 16 3
4 16 5
5 9 3
6 1 5
7 4 8
##end
0 9 5
0-4
0-6
1-3
4-3
5-2
3-5
4-2
2-1
7-6
7-2
7-4
6-5
```

## Output Format

The program outputs:
1. **Complete input file content** (number of ants, rooms, and tunnels)
2. **Movement sequence** showing ants moving between rooms

### Sample Output

```
3
##start
1 23 3
2 16 7
3 16 3
4 16 5
5 9 3
6 1 5
7 4 8
##end
0 9 5
0-4
0-6
1-3
4-3
5-2
3-5
4-2
2-1
7-6
7-2
7-4
6-5

L1-3 L2-2
L1-4 L2-5 L3-3
L1-0 L2-6 L3-4
L2-0 L3-0
```

Where `Lx-y` means ant number `x` moves to room `y`.

## Rules and Constraints

- **Room Naming**: Rooms cannot start with 'L' or '#' and must contain no spaces
- **Room Capacity**: Each room holds maximum one ant (except start/end rooms)
- **Tunnel Usage**: Each tunnel can only be used once per turn
- **Movement**: Ants can only move through tunnels to adjacent empty rooms
- **Coordinates**: Room coordinates must be integers

## Error Handling

The program validates input and returns appropriate error messages:

- `ERROR: invalid data format` - General format errors
- More specific errors can include:
  - Invalid number of ants
  - Missing start or end room
  - Duplicate rooms
  - Invalid coordinates
  - Links to unknown rooms
  - Self-linking rooms
