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

specie = singular for species
species = plural for species

- each cell can belong to a specie each with its own `B/S` rule
- the world is a wrapped torodial grid
- at every step, the engine computes the next state in parallel with multiple workers

### neighbor counting

for each cell:

- count total neighbors in `totalNeighbors`
- count per-species neighbors in `specieNeighbors`

### survival

if the current cell is alive:

- check its own specie survival condition using `totalNeighbors`
- if it fails, the cell dies and can later be replaced by a birth of a different specie

### birth / takeover

each specie independently checks if it can birth at the current cells location using its own neighbor count from `specieNeighbors[specie.Id]`

a specie is added as a candidate for birth if:

- the cell meets its birth condition AND
	- the cell is dead OR
	- the cell is alive but the specie is different from the current cells specie

```go
canCompete := shouldBirth && (!cellIsAlive || differentSpecie)
````

since species can replace other species cells, survival doesnt make a cell invincible

### handling conflicts

- each candidates weight = number of neighbors of its own specie (`specieNeighbors[candidateId]`)
- the candidate with the highest weight wins
- if multiple candidates tie for the highest weight, there is no change to the cell

### summary

- each cell counts neighbors per specie and total neighbors
- a cell can survive according to its species survival rule
- species can attempt to birth into empty cells or replace other species
- multiple candidates resolve with neighbor-weight; ties keep the cell unchanged
- the simulation supports multi-threading for faster computation
