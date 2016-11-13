package main

import (
	"fmt"
	"os"
	"sort"
)

const (
	WIDTH     = 30
	HEIGHT    = 20
	MAX_DEPTH = 2
	FREE      = -1
)

var myID int = -1

type gameMap_t [WIDTH][HEIGHT]int
type pCoords map[int]coordinate

type coordinate struct {
	x int
	y int
}

func dirToStr(from coordinate, to coordinate) string {
	switch {
	case from.x < to.x && from.y == to.y:
		return "RIGHT"
	case from.x > to.x && from.y == to.y:
		return "LEFT"
	case from.y < to.y && from.x == to.x:
		return "DOWN" // inverted y
	case from.y > to.y && from.x == to.x:
		return "UP" // inverted y
	}
	return "INVALID"
}

func getNeighbours(coord coordinate) []coordinate {
	neighbours := make([]coordinate, 0, 4)
	if coord.x-1 >= 0 {
		neighbours = append(neighbours, coordinate{coord.x - 1, coord.y})
	}
	if coord.x+1 < WIDTH {
		neighbours = append(neighbours, coordinate{coord.x + 1, coord.y})
	}
	if coord.y-1 >= 0 {
		neighbours = append(neighbours, coordinate{coord.x, coord.y - 1})
	}
	if coord.y+1 < HEIGHT {
		neighbours = append(neighbours, coordinate{coord.x, coord.y + 1})
	}
	return neighbours
}

func (g gameMap_t) getFreeNeighbours(coord coordinate) []coordinate {
	freeNeighbours := make([]coordinate, 0, 4)
	for _, n := range getNeighbours(coord) {
		if g.get(n) == FREE {
			freeNeighbours = append(freeNeighbours, n)
		}
	}
	return freeNeighbours
}

func (g gameMap_t) get(coord coordinate) int {
	return g[coord.x][coord.y]
}

func (g *gameMap_t) set(coord coordinate, val int) {
	g[coord.x][coord.y] = val
}

func (g *gameMap_t) boardValue(players pCoords) (int, bool) {

	wallhug := len(g.getFreeNeighbours(players[myID]))
	p, e, alone := g.getValues(players)
	score := 10000*p - 10*e - wallhug

	fmt.Fprintln(os.Stderr, players[myID], p, e, alone)

	return score, alone
}

func (g *gameMap_t) preparePlayers(players pCoords) ([]int, map[int][]coordinate) {
	startingPositions := make(map[int][]coordinate)

	pre := make([]int, 0, 3)
	order := make([]int, 0, 4)

	for p, pos := range players {
		if p < myID {
			pre = append(pre, p)
		}
		if p > myID {
			order = append(pre, p)
		}

		startingPositions[p] = append(startingPositions[p], pos)
	}
	sort.Ints(pre)
	sort.Ints(order)
	order = append(order, pre...)
	order = append(order, myID)

	return order, startingPositions
}

func (g gameMap_t) getValues(players pCoords) (int, int, bool) {
	var order []int
	var positions map[int][]coordinate
	order, positions = g.preparePlayers(players)

	pFields := 0
	eFields := 0
	canExplore := true
	alone := true

	for canExplore {
		canExplore = false
		movesThisTurn := make(map[coordinate]int)

		for _, player := range order {
			for _, coord := range positions[player] {
				neighbours := getNeighbours(coord)
				for _, n := range neighbours {
					value := g.get(n)
					if value == FREE {
						canExplore = true
						g.set(n, player+10)
						movesThisTurn[n] = player
					} else if value >= 10 && value != myID+10 && player == myID {
						alone = false
					}
				}
			}
		}

		positions = make(map[int][]coordinate)

		for move, player := range movesThisTurn {
			if player == myID {
				pFields++
			} else {
				eFields++
			}
			positions[player] = append(positions[player], move)
		}
	}
	return pFields, eFields, alone
}

func (p pCoords) getEnemy() int {
	for id, _ := range p {
		if id != myID {
			return id
		}
	}

	return 0
}

func getNextMove(g gameMap_t, players pCoords) coordinate {
	nextMove := coordinate{-1, -1}

	myPos := players[myID]
	fn := g.getFreeNeighbours(myPos)
	if len(fn) < 0 {
		return nextMove
	}

	maxV := 0

	for _, n := range fn {
		g.set(n, myID)
		players[myID] = n

		v, _ := g.boardValue(players)

		if v > maxV {
			maxV = v
			nextMove = n
		}

		g.set(n, FREE)
	}

	players[myID] = myPos

	return nextMove
}

func (g *gameMap_t) init() {
	for x := 0; x < WIDTH; x++ {
		for y := 0; y < HEIGHT; y++ {
			g[x][y] = FREE
		}
	}
}

func (g *gameMap_t) clean(id int) {
	for x := 0; x < WIDTH; x++ {
		for y := 0; y < HEIGHT; y++ {
			if g[x][y] == id {
				g[x][y] = FREE
			}
		}
	}
}

func main() {

	players := make(pCoords)
	stillInGame := [4]bool{true, true, true, true}
	var gameMap gameMap_t
	gameMap.init()

	for {
		// N: total number of players (2 to 4).
		// P: your player number (0 to 3).
		var N, P int
		fmt.Scan(&N, &P)
		var starting, current coordinate

		for i := 0; i < N; i++ {
			// X0: starting X coordinate of lightcycle (or -1)
			// Y0: starting Y coordinate of lightcycle (or -1)
			// X1: current X coordinate of lightcycle
			// Y1: current Y coordinate of lightcycle

			fmt.Scan(&starting.x, &starting.y, &current.x, &current.y)

			if myID < 0 {
				gameMap.set(starting, i)
			}

			if current.x == -1 {
				if stillInGame[i] {
					stillInGame[i] = false
					gameMap.clean(i)
					delete(players, i)
				}
			} else {
				players[i] = current
				gameMap.set(current, i)
			}

		}

		if myID < 0 {
			myID = P
		}

		nextMove := getNextMove(gameMap, players)

		fmt.Println(dirToStr(players[myID], nextMove))
	}
}
