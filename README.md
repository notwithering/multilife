# multilife

a program to place multiple conway's-game-of-life-like cellular automata in a single ecosystem that can all interact with eachother

run the program and it outputs `output.mp4` file with the video replay of the simulation

https://www.youtube.com/watch?v=zpPAf_UoUCc

<img width="1920" height="1080" alt="Image" src="https://github.com/user-attachments/assets/ce6326d4-7f90-43ca-bd3a-73dedacb3753" />

<img width="1920" height="1080" alt="Image" src="https://github.com/user-attachments/assets/a8d23037-b348-4a47-be6e-d5472f48d846" />

## license

the programs code is licensed under the [MIT license](LICENSE)

all videos created using this software are licensed under [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/)

## credits

- Michaelangel007's [nanofont3x4](https://github.com/Michaelangel007/nanofont3x4)
- FFmpeg for video encoding
- cellular automaton rules from various authors found from [conwaylife.com](https://conwaylife.com/wiki/List_of_Life-like_rules) and [hatsyas catagolue](https://catagolue.hatsya.com/rules/lifelike)

## technical overview of engine

- each cell can belong to a species each with its own `B/S` rule
- the world is a wrapped torodial grid
- at every step, the engine computes the next state

### neighbor counting

for each cell:

- count total neighbors in `totalNeighbors`
- count per-species neighbors in `specieNeighbors`

### survival

if the current cell is alive, it checks its own species survival condition with its total neighbors of any specie, if the rule fails, the cell dies or could be later replaced by a birth of a different specie

### birth / takeover

each specie independently checks if it can birth at this location using its own neighbor count (`specieNeighbors[specie.Id]`). a specie is considered a candidate for birth if the current cell meets the birth condition AND:

- the current cell is alive, and the aggressing specie is not the same as the current cells specie, or
- the current cell is dead

```go
canCompete := shouldBirth &&
	((cellIsAlive && differentSpecie) ||
		(!cellIsAlive))
````

species can attack other species cells, so survival doesnt make a cell invincible

### handling conflicts

each candidate has a weight equal to the total number of neighbors of its own specie, the candidate with the most number of neighbors wins. if two or more share the highest weight, nothing happens and the current cell remains in its current state

### summary

- each cell checks all species neighbor counts
- each specie can birth into a cell if its birth rule is met
- the current cell can survive if its survival rule is met
- if multiple species can birth at the cell, the one with the highest local neighbor density (`ownSpecieNeighbors/(totalNeighbors + 1)`)
- if theres a tie, the cell is left untouched
- species can freely replace others cells
