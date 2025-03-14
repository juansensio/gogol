package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var Nx int = 160
var Ny int = 40
var ITS int = 1000
var WORKERS int = 4

func printCellsTerminal(i int, cells []int) {
	fmt.Print("\033[H\033[2J") // Clear terminal and move cursor to top-left
	fmt.Printf("Iteration: %d\n", i)
	for i := 1; i <= Ny; i++ {
		for j := 1; j <= Nx; j++ {
			if cells[i*(Nx+2)+j] == 1 {
				fmt.Print("@")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func matrixVectorMultiplicationSparseParallel(sparse [][]int, vector []int) []int {
	out := make([]int, len(vector))
	rowPtr := sparse[0]
	colIdx := sparse[1]
	values := sparse[2]

	var wg sync.WaitGroup
	wg.Add(WORKERS)

	// Calculate rows per worker
	totalRows := len(vector)
	rowsPerWorker := totalRows / WORKERS

	// Launch workers
	for w := 0; w < WORKERS; w++ {
		startRow := w * rowsPerWorker
		endRow := startRow + rowsPerWorker
		if w == WORKERS-1 {
			endRow = totalRows // Make sure the last worker processes any remaining rows
		}

		go func(start, end int) {
			defer wg.Done()
			// For each row in this worker's range
			for i := start; i < end; i++ {
				// Get range of non-zero elements for this row
				rowStart := rowPtr[i]
				rowEnd := rowPtr[i+1]

				// Multiply and sum non-zero elements
				for j := rowStart; j < rowEnd; j++ {
					col := colIdx[j]
					val := values[j]
					out[i] += val * vector[col]
				}
			}
		}(startRow, endRow)
	}

	wg.Wait()
	return out
}

func applyRulesParallel(cells []int, alive []int) {
	var wg sync.WaitGroup
	wg.Add(WORKERS)

	// Calculate cells per worker
	totalCells := len(cells)
	cellsPerWorker := totalCells / WORKERS

	// Launch workers
	for w := 0; w < WORKERS; w++ {
		startCell := w * cellsPerWorker
		endCell := startCell + cellsPerWorker
		if w == WORKERS-1 {
			endCell = totalCells // Make sure the last worker processes any remaining cells
		}

		go func(start, end int) {
			defer wg.Done()
			// Apply rules to this worker's range of cells
			for i := start; i < end; i++ {
				if cells[i] == 1 {
					if alive[i] == 2 || alive[i] == 3 {
						cells[i] = 1
					} else {
						cells[i] = 0
					}
				} else {
					if alive[i] == 3 {
						cells[i] = 1
					} else {
						cells[i] = 0
					}
				}
			}
		}(startCell, endCell)
	}

	wg.Wait()
}

func main() {
	// initialize cells
	cells := make([]int, (Nx+2)*(Ny+2))
	for i := 1; i <= Ny; i++ {
		for j := 1; j <= Nx; j++ {
			if rand.Float32() < 0.3 { // 30% chance of being alive
				cells[i*(Nx+2)+j] = 1
			}
		}
	}

	// Create matrix of neighbors
	nnz := 8 * Nx * Ny // each cell has 8 neighbors
	// Create CSR format arrays
	neighbors := make([][]int, 3)
	neighbors[0] = make([]int, (Nx+2)*(Ny+2)+1) // row pointers
	neighbors[1] = make([]int, nnz)             // column indices
	neighbors[2] = make([]int, nnz)             // values
	// Fill arrays directly
	idx := 0
	// Initialize all row pointers to their correct starting positions
	for i := 0; i < (Nx+2)*(Ny+2); i++ {
		neighbors[0][i] = idx
		// Only add neighbors for cells in the active grid
		row := i / (Nx + 2)
		col := i % (Nx + 2)
		if row >= 1 && row <= Ny && col >= 1 && col <= Nx {
			for di := -1; di <= 1; di++ {
				for dj := -1; dj <= 1; dj++ {
					if di == 0 && dj == 0 {
						continue
					}
					ni := row + di
					nj := col + dj
					if ni >= 0 && ni < Ny+2 && nj >= 0 && nj < Nx+2 {
						neighbors[1][idx] = ni*(Nx+2) + nj
						neighbors[2][idx] = 1
						idx++
					}
				}
			}
		}
	}
	neighbors[0][(Nx+2)*(Ny+2)] = idx

	// iterate
	for i := 0; i < ITS; i++ {
		// matrix-vector multiplication (parallel version)
		alive := matrixVectorMultiplicationSparseParallel(neighbors, cells)

		// apply rules (parallel version)
		applyRulesParallel(cells, alive)

		// wait
		time.Sleep(50 * time.Millisecond)

		// print cells
		printCellsTerminal(i, cells)
	}
}
