package main

import (
	"fmt"
	"math/rand"
	"time"
)

var Nx int = 160
var Ny int = 40
var ITS int = 1000

func printCellsTerminal(i int, cells [][]bool) {
	fmt.Print("\033[H\033[2J") // Clear terminal and move cursor to top-left
	fmt.Printf("Iteration: %d\n", i)
	for i := range cells {
		for j := range cells[i] {
			if cells[i][j] {
				fmt.Print("@")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func updateCell(cells [][]bool, i int, j int) bool {
	// Count live neighbors, handling edges carefully
	liveNeighbors := 0
	for di := -1; di <= 1; di++ {
		for dj := -1; dj <= 1; dj++ {
			// Skip the cell itself
			if di == 0 && dj == 0 {
				continue
			}
			// Calculate neighbor coordinates
			ni := i + di
			nj := j + dj
			// Check bounds
			if ni >= 0 && ni < Ny && nj >= 0 && nj < Nx {
				if cells[ni][nj] {
					liveNeighbors++
				}
			}
		}
	}
	// Apply Game of Life rules
	if cells[i][j] {
		// Live cell
		return liveNeighbors == 2 || liveNeighbors == 3
	} else {
		// Dead cell
		return liveNeighbors == 3
	}
}

func main() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// initialize cells
	cells := make([][]bool, Ny)
	for i := range cells {
		cells[i] = make([]bool, Nx)
		for j := range cells[i] {
			if rand.Float32() < 0.3 { // 30% chance of being alive
				cells[i][j] = true
			}
		}
	}

	// iterate
	for i := 0; i < ITS; i++ {
		// Create a new grid for the next state
		nextCells := make([][]bool, Ny)
		for y := range nextCells {
			nextCells[y] = make([]bool, Nx)
		}

		// Calculate the next state based on the current state
		for y := range cells {
			for x := range cells[y] {
				nextCells[y][x] = updateCell(cells, y, x)
			}
		}

		// Replace the old grid with the new one
		cells = nextCells

		// wait
		time.Sleep(100 * time.Millisecond)

		// print cells
		printCellsTerminal(i, cells)
	}
}
