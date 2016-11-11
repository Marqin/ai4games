# Copyright (c) 2016 Hubert Jarosz
#
# This software is provided 'as-is', without any express or implied
# warranty. In no event will the authors be held liable for any damages
# arising from the use of this software.
#
# Permission is granted to anyone to use this software for any purpose,
# including commercial applications, and to alter it and redistribute it
# freely, subject to the following restrictions:
#
# 1. The origin of this software must not be misrepresented; you must not
#    claim that you wrote the original software. If you use this software
#    in a product, an acknowledgement in the product documentation would be
#    appreciated but is not required.
# 2. Altered source versions must be plainly marked as such, and must not be
#    misrepresented as being the original software.
# 3. This notice may not be removed or altered from any source distribution.

# currently at gold

import sys  # for sys.stderr
import copy  # for deepcopy
import collections  # for deque
import math  # for sqrt
import operator  # for itemgetter

WIDTH = 30
HEIGHT = 20
MAX_DEPTH = 2

################################################################################

# checks if given coordinate is on map
def is_on_map(x, y):
    if x < 0 or x >= WIDTH or y < 0 or y >= HEIGHT:
        return False
    return True

# gives map field neighbours that are on map
def get_neighbours(x, y):
    first_set = {(x-1, y), (x+1, y), (x, y-1), (x, y+1)}

    return {(vx, vy) for (vx, vy) in first_set if is_on_map(vx, vy)}

# compute distance between two points
def distance(p0, p1):
    return math.sqrt((p0[0] - p1[0])**2 + (p0[1] - p1[1])**2)

#
def getEnemy(players, playerID):
    pl = list(players)
    return pl[0] if pl[0] != playerID else pl[1]

################################################################################

class GameMap:
    def __init__(self, width, height):
        self.width = width
        self.height = height
        self.map = [[-1]*self.height for n in range(self.width)]
        #self.map = [[-1] * self.height] * self.width   ## BAD!!!!
    def set(self, x, y, val=0):
        self.map[x][y] = val
    def get(self, x, y):
        return self.map[x][y]

    def freeNeighbours(self, x, y):
        neibs = get_neighbours(x,y)
        fn = [ (nx,ny) for (nx,ny) in neibs if self.get(nx,ny) < 0 ]
        return fn

    def clone(self):
        return copy.deepcopy(self)

    def __moveWest(self, w, y):
        TARGET = -1
        while w >= 0 and self.map[w][y] == TARGET:
            w -= 1
        return w+1

    def __moveEast(self, e, y):
        TARGET = -1
        while e < WIDTH and self.map[e][y] == TARGET:
            e += 1
        return e ## no -1 because we use it in range() - dirty speedhack

    def expand(self, start):

        ID = self.get(*start)
        fn = self.freeNeighbours(*start)
        for x, y in fn:
            self.set(x, y, ID)
        return set(fn)

    def fill(self, playerID, players):

        startingPositions = {p: [coord] for p, coord in players.items()}
        gameMap = self.clone()

        playerList = list(players)
        data = {pid: {} for pid in playerList}

        playerOrder = [playerID] + [p for p in playerList if p > playerID] + [p for p in playerList if p < playerID]

        turnNo = 1
        canExplore = True

        while canExplore:
            canExplore = False
            movesThisTurn = {}
            for player in playerOrder:
                for pos in startingPositions[player]:
                    for n in gameMap.freeNeighbours(*pos):
                        canExplore = True
                        gameMap.map[pos[0]][pos[1]] = player
                        movesThisTurn[n] = player
            #
            for move, player in movesThisTurn.items():
                data[player][move] = turnNo

            if canExplore:
                turnNo += 1
                startingPositions = {player: [move for move, p in movesThisTurn.items() if p == player] for player in playerList}
        #

        playerTiles = len(data[playerID])
        enemyTiles = sum([len(data[i]) for i in playerList]) - playerTiles

        return (playerTiles, enemyTiles)


    def clean(self, pid):
        for x in range(len(self.map)):
            for y in range(len(self.map[x])):
                if self.map[x][y] == pid:
                    self.map[x][y] = -1



################################################################################


def moveScore(gM, playerID, players):
    p, e = gM.fill(playerID, players)
    score = 10000 * p - 10 * e - len(gM.freeNeighbours(*players[playerID])) # + acc

    #print(players[playerID], p, file=sys.stderr)
    return score


# def minmax(gameMap, depth, playerID, maximizing, players, alpha, beta):
#
#     enemyID = getEnemy(players, playerID)
#     #enemyID = 0 if playerID > 0 else 1
#
#     fn = gameMap.freeNeighbours(*players[playerID])
#
#     if depth == 0 or len(fn) == 0:
#         x, y = players[playerID]
#         return moveScore(gameMap, playerID, players)
#
#     if maximizing:
#         v = []
#         for n_x, n_y in fn:
#             gM = gameMap.clone()
#             p = copy.deepcopy(players)
#             gM.set(n_x, n_y, playerID)
#             p[playerID] = (n_x, n_y)
#             alpha = max(alpha, minmax(gM, depth - 1, playerID, True, p, alpha, beta))
#             if alpha >= beta:
#                 break # cut beta
#         return alpha
    # else:
    #     x, y = players[enemyID]
    #     fn = gameMap.freeNeighbours(x, y)
    #
    #     v = []
    #     for n_x, n_y in fn:
    #         gM = gameMap.clone()
    #         p = copy.deepcopy(players)
    #         gM.set(n_x, n_y, enemyID)
    #         p[enemyID] = (n_x, n_y)
    #         beta = min(beta, minmax(gM, depth - 1, playerID, True, p, alpha, beta))
    #         if alpha >= beta:
    #             break # cut alpha
    #     return beta

################################################################################


def playField(gameMap, playerID, players):
    #print("start:", players[playerID], file=sys.stderr)
    fn = gameMap.freeNeighbours(*players[playerID])
    #print("fn:", fn, file=sys.stderr)

    if len(fn) == 0:
        print("Dead.", file=sys.stderr)
        return (-1, -1)

    maxV = 0
    n = fn[0]
    alpha, beta = float("-inf"), float("inf")
    for n_x, n_y in fn:
        gM = gameMap.clone()
        p = copy.deepcopy(players)
        gM.set(n_x, n_y, playerID)

        p[playerID] = (n_x, n_y)

        v = moveScore(gM, playerID, p)

        #v = minmax(gM, MAX_DEPTH, playerID, True, p , alpha, beta)

        if v > maxV:
            maxV = v
            n = (n_x, n_y)
    return n


################################################################################

gameMap = GameMap(WIDTH, HEIGHT)

stillInGame = None
players = dict()

firstTurn = True

# game loop
while True:
    # n: total number of players (2 to 4).
    # p: your player number (0 to 3).
    n, p = [int(i) for i in input().split()]

    if stillInGame == None:
        stillInGame = [True] * n

    x, y = -1, -1
    nbr = 0
    enX, enY = -1, -1
    for i in range(n):

        x0, y0, x1, y1 = [int(j) for j in input().split()]

        if firstTurn:
            gameMap.set(x0, y0, i)

        if x1 != -1:
            players[i] = (x1, y1)
            gameMap.set(x1, y1, i)
        else:
            if stillInGame[i]:
                gameMap.clean(i)
                stillInGame[i] = False
                del players[i]

    firstTurn = False

    x, y = players[p]

    if (x,y) == (-1,-1):
        print("I'm dead.", file=sys.stderr)
        break


    go_to = playField(gameMap, p, players)
    print((x,y),"->", go_to, file=sys.stderr)

    if go_to[0] > x:
        print("RIGHT")
    elif go_to[0] < x:
        print("LEFT")
    elif go_to[1] > y:
        print("DOWN")
    else:
        print("UP")
    # To debug: print("Debug messages...", file=sys.stderr)
