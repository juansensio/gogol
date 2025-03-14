package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var Nx int = 160
var Ny int = 40
var ITS int = 100
var WORKERS int = 4

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

// Create a struct to hold both the result and the worker ID
type WorkerResult struct {
	cells    [][]bool
	workerID int
}

func updateCellParallel(cells [][]bool, nextCells [][]bool, startRow, endRow int, workerID int, wg *sync.WaitGroup) {
	defer wg.Done()

	// Process only the assigned rows
	for i := startRow; i < endRow; i++ {
		for j := range cells[i] {
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
				// Live cell continues to live if it has 2 or 3 neighbors, otherwise it dies
				nextCells[i][j] = liveNeighbors == 2 || liveNeighbors == 3
			} else {
				// Dead cell becomes live if it has exactly 3 neighbors, otherwise it stays dead
				nextCells[i][j] = liveNeighbors == 3
			}
		}
	}
}

func updateCellsConcurrent(cells [][]bool) [][]bool {
	// Create a new grid for the next state
	nextCells := make([][]bool, Ny)
	for i := range nextCells {
		nextCells[i] = make([]bool, Nx)
	}

	// Use WaitGroup for synchronization
	var wg sync.WaitGroup
	wg.Add(WORKERS)

	// Calculate rows per worker
	rowsPerWorker := Ny / WORKERS

	// Launch workers
	for w := 0; w < WORKERS; w++ {
		startRow := w * rowsPerWorker
		endRow := startRow + rowsPerWorker
		if w == WORKERS-1 {
			endRow = Ny // Make sure the last worker processes any remaining rows
		}
		go updateCellParallel(cells, nextCells, startRow, endRow, w, &wg)
	}

	// Wait for all workers to complete
	wg.Wait()

	return nextCells
}

func main() {
	// Add import for sync package at the top of the file
	// import "sync"

	// Set random seed for more varied patterns
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
		cells = updateCellsConcurrent(cells)
		time.Sleep(100 * time.Millisecond)
		printCellsTerminal(i, cells)
	}
}
