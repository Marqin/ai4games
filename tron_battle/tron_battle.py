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

# currently at bronze

import sys  # for sys.stderr
import copy  # for deepcopy
import collections  # for deque

WIDTH = 30
HEIGHT = 20
MAX_DEPTH = 4

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

    def flood(self, startX, startY):
        floodMap = self.clone()
        TARGET = -1
        REPLACE = -2
        size = 0
        if floodMap.get(startX, startY) != TARGET:
            return size
        Q = collections.deque([(startX, startY)])
        while Q:
            x, y = Q.popleft()
            if floodMap.map[x][y] != TARGET:
                continue

            for i in range(self.__moveWest(x, y), self.__moveEast(x, y)):
                floodMap.map[i][y] = REPLACE
                size += 1
                if y < (HEIGHT-1) and floodMap.map[i][y+1] == TARGET:
                    Q.append((i, y+1))
                if y > 0 and floodMap.map[i][y-1] == TARGET:
                    Q.append((i, y-1))
        return size

################################################################################


def moveScore(gM, x, y):
    tmp = gameMap.get(x, y)
    gameMap.set(x, y, -1)
    score = gameMap.flood(x, y)
    gameMap.set(x, y, tmp)
    return score

def minmax(gameMap, depth, playerID, maximizing, pos, enPos, alpha, beta):
    enemyID = 0 if playerID > 0 else 1

    if maximizing:
        x, y = pos
        fn = gameMap.freeNeighbours(x, y)

        if depth == 0 or len(fn) == 0:
            return moveScore(gameMap, x, y)

        v = []
        for n_x, n_y in fn:
            gM = gameMap.clone()
            gM.set(n_x, n_y, playerID)
            alpha = max(alpha, minmax(gM, depth - 1, playerID, False, (n_x, n_y), enPos, alpha, beta))
            if alpha >= beta:
                break # cut beta
        return alpha
    else:
        x, y = enPos
        fn = gameMap.freeNeighbours(x, y)

        if depth == 0 or len(fn) == 0:
            x, y = pos
            return moveScore(gameMap, x, y)

        v = []
        for n_x, n_y in fn:
            gM = gameMap.clone()
            gM.set(n_x, n_y, enemyID)

            beta = min(beta, minmax(gM, depth - 1, playerID, True, pos, (n_x, n_y), alpha, beta))
            if alpha >= beta:
                break # cut alpha
        return beta

################################################################################


def playField(gameMap, playerID, pos, enPos):
    fn = gameMap.freeNeighbours(x, y)
    maxV = 0
    n = fn[0]
    alpha, beta = float("-inf"), float("inf")
    for n_x, n_y in fn:
        gM = gameMap.clone()
        gM.set(n_x, n_y, playerID)
        v = minmax(gM, MAX_DEPTH, playerID, False, (n_x, n_y), enPos, alpha, beta)

        if v > maxV:
            maxV = v
            n = (n_x, n_y)
    return n


################################################################################

gameMap = GameMap(WIDTH, HEIGHT)

# game loop
while True:
    # n: total number of players (2 to 4).
    # p: your player number (0 to 3).
    n, p = [int(i) for i in input().split()]
    x, y = -1, -1
    nbr = 0
    enX, enY = -1, -1
    for i in range(n):

        x0, y0, x1, y1 = [int(j) for j in input().split()]


        if x1 != -1:
            gameMap.set(x1, y1, i+5)

        if i == p:
            x, y = x1, y1
        else:
            enX, enY = x1, y1

    if (x,y) == (-1,-1):
        break

    go_to = playField(gameMap, p, (x, y), (enX, enY))

    if go_to[0] > x:
        print("RIGHT")
    elif go_to[0] < x:
        print("LEFT")
    elif go_to[1] > y:
        print("DOWN")
    else:
        print("UP")
    # To debug: print("Debug messages...", file=sys.stderr)
