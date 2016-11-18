# Tron Battle

In this game you drive a [Tron](https://en.wikipedia.org/wiki/Tron_(video_game))-like cycle on rectangular map and have to outlive (you die when crashing into your own/someones trail or into walls) your opponents.

## Bronzy start

At first I've made a simple greedy algorithm in Python that went to the
field that had the smallest number of free neighbours (but bigger than 1).

## Silver war(e)

I was going to write Min-Max for this game, but for it to succeed I needed some
good board value function. My first one was basic "flood fill" - I was flooding
map from my position to count how much fields I can go and choose direction with
more fields. Then I've implemented Min-Max with alpha-beta pruning... And of
course it was not enough to get into silver! I was looking for ideas, and then
I've noticed that my value is ignoring enemies. That's when I've found about
Voronoi, and I've used it as my value: I count fields closer to given player
than to others, then I count board value as:
```
BIG_NUMBER*myFields - SMALL_NUMBER*enemyFields - freeNeighbours
```
It was enough, to get to silver.

## Highway to Gold

After adding support for more than one enemy and fixing ugly bugs, I've landed
on top of Silver. That's when I've noticed I'm getting many timeouts, so I tried
disabling Min-Max (left just one-turn Voronoi for board value) - surprisingly it
was enough to get to Gold.

## Legendary struggle
I've tried optimizing, disabled Min-Max for case when I'm closed alone in
a space, and even fixed some bugs. It was not enough even for #2 in Gold. I had enough too.

I've decided to rewrite my code to some faster language.  
And I choose `Go`, because it's quite fast and I've never used it (but heard a lot about it).

It took me a while to fully rewrite it to Go, fix many bugs, optimize it, wrote test, and I'd finally got my #2 place Gold. But I was still worse than Boss. I've dug into all those match replays and I got the pattern - my bot often died by closing itself inside self! So I've wrote a condition - if I'm close alone I go in such way that I leave path to enemies that are worse (in term of board value) than me. That way I can wait in my closing until they die, then go out, take their place and win! I got to Legendary with this.
