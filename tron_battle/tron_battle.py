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
            if floodMap.get(x,y) != TARGET:
                continue
            size += 1
            #print(size, file=sys.stderr)
            floodMap.set(x, y, REPLACE)
            n = [fn for fn in floodMap.freeNeighbours(x, y) if floodMap.get(fn[0], fn[1]) == TARGET]
            Q.extend(n)
        return size

################################################################################

gameMap = GameMap(WIDTH, HEIGHT)

# game loop
while True:
    # n: total number of players (2 to 4).
    # p: your player number (0 to 3).
    n, p = [int(i) for i in input().split()]
    x, y = -1, -1
    nbr = 0
    for i in range(n):
        x0, y0, x1, y1 = [int(j) for j in input().split()]

        #print(x0,y0,x1,y1, nbr, file=sys.stderr)

        if x1 != -1:
            gameMap.set(x1, y1, i+5)

        if i == p:
            x, y = x1, y1

    if (x,y) == (-1,-1):
        break

    fn = gameMap.freeNeighbours(x,y)

    if len(fn) == 0:
        print("FAILURE!!!", file=sys.stderr)
        print("DOWN")
        continue

    max_n = 0
    maxScore = 0

    for i in range(len(fn)):
        score = gameMap.flood(*fn[i])
        if score > maxScore:
            max_n = i
            maxScore = score

    go_to = fn[max_n]

    if go_to[0] > x:
        print("RIGHT")
    elif go_to[0] < x:
        print("LEFT")
    elif go_to[1] > y:
        print("DOWN")
    else:
        print("UP")
    # To debug: print("Debug messages...", file=sys.stderr)
