/*
Copyright (c) 2016 Hubert Jarosz

This software is provided 'as-is', without any express or implied
warranty. In no event will the authors be held liable for any damages
arising from the use of this software.

Permission is granted to anyone to use this software for any purpose,
including commercial applications, and to alter it and redistribute it
freely, subject to the following restrictions:

1. The origin of this software must not be misrepresented; you must not
   claim that you wrote the original software. If you use this software
   in a product, an acknowledgement in the product documentation would be
   appreciated but is not required.
2. Altered source versions must be plainly marked as such, and must not be
   misrepresented as being the original software.
3. This notice may not be removed or altered from any source distribution.
*/

package main

import (
	"fmt"
	"math"
	"os"
	"sort"
)

const (
	WIDTH           = 30
	HEIGHT          = 20
	MAX_DEPTH       = 2
	MAX_DEPTH_ALONE = 4
	FREE            = -1
)

var myID int = -1

const maxInt = int(^uint(0) >> 1)
const minInt = -maxInt - 1

func max(a int, b int) int {
	if a < b {
		return b
	}
	return a
}

func min(a int, b int) int {
	if a > b {
		return b
	}
	return a
}

type gameMap_t [WIDTH][HEIGHT]int
type pCoords map[int]coordinate

type coordinate struct {
	x int
	y int
}

type pair struct {
	fst int
	snd int
}

type pairSort []pair

func (p pairSort) Len() int {
	return len(p)
}
func (p pairSort) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p pairSort) Less(i, j int) bool {
	return p[i].fst < p[j].fst
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

func (g *gameMap_t) get(coord coordinate) int {
	return g[coord.x][coord.y]
}

func (g *gameMap_t) set(coord coordinate, val int) {
	g[coord.x][coord.y] = val
}

func (g *gameMap_t) boardValue(players pCoords) (int, bool) {

	wallhug := len(g.getFreeNeighbours(players[myID]))
	fields, alone := g.getValues(players)
	e := 0

	for p, _ := range players {
		if p != myID {
			e += fields[p]
		}
	}

	score := 10000*fields[myID] - 10*e - wallhug
	//fmt.Fprintln(os.Stderr, "normal:", players[myID], fields[myID], fields, score)

	return score, alone
}

func (g *gameMap_t) boardValueAlone(players pCoords) int {

	wallhug := len(g.getFreeNeighbours(players[myID]))
	fields, _ := g.getValues(players)

	worse := make([]pair, 0, 4)

	for p, _ := range players {
		if p != myID {
			if fields[p] < fields[myID] {
				worse = append(worse, pair{fields[p], p})
			}
		}
	}

	sort.Sort(pairSort(worse))

	fmt.Fprintln(os.Stderr, players[myID], worse)

	score := fields[myID] - wallhug

	for _, worse := range worse {
		if g.pathTo(players, worse.snd) {
			score = 10000*fields[myID] - wallhug
			break
		}
	}
	//fmt.Fprintln(os.Stderr, "normal:", score, worse)
	return score
}

func (g *gameMap_t) preparePlayers(players pCoords) ([]int, [4][]coordinate) {
	var startingPositions [4][]coordinate

	pre := make([]int, 0, 3)
	order := make([]int, 0, 4)

	for p, pos := range players {
		if p < myID {
			pre = append(pre, p)
		}
		if p > myID {
			order = append(order, p)
		}

		startingPositions[p] = []coordinate{pos}
	}
	sort.Ints(pre)
	sort.Ints(order)
	order = append(order, pre...)
	order = append(order, myID)

	return order, startingPositions
}

func (g gameMap_t) getValues(players pCoords) ([4]int, bool) {
	order, positions := g.preparePlayers(players)

	fields := [4]int{0, 0, 0, 0}

	canExplore := true
	alone := true

	for canExplore {
		canExplore = false
		var movesThisTurn [4][]coordinate

		for _, player := range order {
			for _, coord := range positions[player] {
				neighbours := getNeighbours(coord)
				for _, n := range neighbours {
					value := g.get(n)
					if value == FREE {
						canExplore = true
						g.set(n, player+10)

						fields[player]++

						movesThisTurn[player] = append(movesThisTurn[player], n)
					} else if value >= 10 && value != myID+10 && player == myID {
						alone = false
					}
				}
			}
		}

		positions = movesThisTurn

	}
	return fields, alone
}

func (g gameMap_t) pathTo(players pCoords, id int) bool {

	moves := getNeighbours(players[myID])
	var newMoves []coordinate

	for len(moves) > 0 {
		for _, move := range moves {

			if g.get(move) == id {
				return true
			}

			g.set(move, myID)

			for _, n := range getNeighbours(move) {
				if g.get(move) == id {
					return true
				} else if g.get(move) == FREE {
					newMoves = append(newMoves, n)
				}
			}
		}

		moves = newMoves
	}

	return false
}

func (p pCoords) getEnemy() int {

	d := (WIDTH + HEIGHT) * 10.0
	enemyID := -1

	for id, coord := range p {
		if id != myID {
			dist := distance(p[myID], coord)
			if dist < d {
				d = dist
				enemyID = id
			}
		}
	}

	return enemyID
}

func getNextMove(g gameMap_t, players pCoords) coordinate {
	nextMove := coordinate{-1, -1}

	myPos := players[myID]
	fn := g.getFreeNeighbours(myPos)
	if len(fn) <= 0 {
		fmt.Fprintln(os.Stderr, "No moves :(")
		return nextMove
	}

	maxV := minInt
	_, alone := g.boardValue(players)

	first := true

	for _, n := range fn {
		g.set(n, myID)
		players[myID] = n

		v := minInt

		if alone {
			v = lookAhead(g, players, MAX_DEPTH_ALONE)
			//v = g.boardValueAlone(players)
			//v, _ = g.boardValue(players)
		} else {
			depth := MAX_DEPTH
			if len(players) < 3 {
				depth = 4
			}
			v = minmax(g, players, depth, false, minInt, maxInt)
			//v, _ = g.boardValue(players)
		}

		if v > maxV || first {
			maxV = v
			nextMove = n
		}

		g.set(n, FREE)
		first = false
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

func distance(c0, c1 coordinate) float64 {
	x := c0.x - c0.x
	y := c0.y - c1.y
	return math.Sqrt(float64(x*x + y*y))
}

func runGame(in *os.File) {

	test := true

	if in == nil {
		test = false
		in = os.Stdin
	}

	players := make(pCoords)
	stillInGame := [4]bool{true, true, true, true}
	var gameMap gameMap_t
	gameMap.init()

	for {
		// N: total number of players (2 to 4).
		// P: your player number (0 to 3).
		var N, P int
		fmt.Fscan(in, &N, &P)
		var starting, current coordinate

		for i := 0; i < N; i++ {
			// X0: starting X coordinate of lightcycle (or -1)
			// Y0: starting Y coordinate of lightcycle (or -1)
			// X1: current X coordinate of lightcycle
			// Y1: current Y coordinate of lightcycle

			fmt.Fscan(in, &starting.x, &starting.y, &current.x, &current.y)

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

		if test {
			return
		} else {
			fmt.Fprintln(os.Stderr, players[myID], "->", nextMove)
			fmt.Println(dirToStr(players[myID], nextMove))
		}
	}
}

func main() {
	runGame(nil)
}

func minmax(g gameMap_t, players pCoords, depth int, maximizing bool, alpha int, beta int) int {
	fn := g.getFreeNeighbours(players[myID])

	if depth == 0 || len(fn) <= 0 {
		v, _ := g.boardValue(players)
		return v
	}

	if maximizing {
		for _, n := range fn {
			p := players.clone()
			g.set(n, myID)
			p[myID] = n

			alpha = max(alpha, minmax(g, p, depth-1, false, alpha, beta))
			g.set(n, FREE)

			if alpha >= beta {
				break
			}
		}
		return alpha
	} else {

		for enemyID, coord := range players {
			if enemyID != myID {
				fn = g.getFreeNeighbours(coord)
				for _, n := range fn {
					p := players.clone()

					g.set(n, enemyID)
					p[enemyID] = n

					beta = min(beta, minmax(g, p, depth-1, true, alpha, beta))
					g.set(n, FREE)

					if alpha >= beta {
						return beta
					}
				}
			}
		}
		return beta
	}
}

func (p pCoords) clone() pCoords {
	pClone := make(pCoords)
	for k, v := range p {
		pClone[k] = v
	}
	return pClone
}

func lookAhead(g gameMap_t, players pCoords, depth int) int {
	fn := g.getFreeNeighbours(players[myID])

	if depth == 0 || len(fn) <= 0 {
		return g.boardValueAlone(players)
	}

	maxV := minInt

	for _, n := range fn {
		p := players.clone()
		g.set(n, myID)
		p[myID] = n

		maxV = max(maxV, lookAhead(g, p, depth-1))
		g.set(n, FREE)

	}
	return maxV
}
