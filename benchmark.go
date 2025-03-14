package main

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

func naive(cells [][]bool, Nx, Ny, ITS int) {
	// iterate
	for i := 0; i < ITS; i++ {
		// Create a new grid for the next state
		nextCells := make([][]bool, Ny)
		for y := range nextCells {
			nextCells[y] = make([]bool, Nx)
		}
		// Calculate the next state based on the current state
		for i := range cells {
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
		// Replace the old grid with the new one
		cells = nextCells
	}
}

func pad(cells [][]bool, Nx, Ny, ITS int) {
	// iterate
	for i := 0; i < ITS; i++ {
		// Create a new grid for the next state
		nextCells := make([][]bool, Ny+2)
		for y := range nextCells {
			nextCells[y] = make([]bool, Nx+2)
		}
		// Calculate the next state based on the current state
		for i := 1; i <= Ny; i++ {
			for j := 1; j <= Nx; j++ {
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
						if cells[ni][nj] {
							liveNeighbors++
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
		// Replace the old grid with the new one
		cells = nextCells
	}
}

type Board struct {
	cells  [][]bool
	_cells [][]bool
	w      int
	h      int
}

// methods

func (b *Board) init(Nx, Ny int) {
	b.cells = make([][]bool, Ny+2)
	b._cells = make([][]bool, Ny+2)
	for i := range b.cells {
		b.cells[i] = make([]bool, Nx+2)
		b._cells[i] = make([]bool, Nx+2)
	}
	for i := 1; i <= Ny; i++ {
		for j := 1; j <= Nx; j++ {
			if rand.Float32() < 0.3 { // 30% chance of being alive
				b.cells[i][j] = true
				b._cells[i][j] = true
			}
		}
	}
	b.w = Nx + 2
	b.h = Ny + 2
}

func (b *Board) update() {
	for i := 1; i <= b.h-2; i++ {
		for j := 1; j <= b.w-2; j++ {
			liveNeighbors := 0
			for di := -1; di <= 1; di++ {
				for dj := -1; dj <= 1; dj++ {
					if di == 0 && dj == 0 {
						continue
					}
					ni := i + di
					nj := j + dj
					if b._cells[ni][nj] {
						liveNeighbors++
					}
				}
			}
			// Apply Game of Life rules
			if b._cells[i][j] {
				// Live cell continues to live if it has 2 or 3 neighbors, otherwise it dies
				b.cells[i][j] = liveNeighbors == 2 || liveNeighbors == 3
			} else {
				// Dead cell becomes live if it has exactly 3 neighbors, otherwise it stays dead
				b.cells[i][j] = liveNeighbors == 3
			}
		}
	}
	for i := 1; i <= b.h-2; i++ {
		for j := 1; j <= b.w-2; j++ {
			b._cells[i][j] = b.cells[i][j]
		}
	}
}

func structed(b *Board, ITS int) {
	for i := 0; i < ITS; i++ {
		b.update()
	}
}

func matrixVectorMultiplication(matrix [][]int, vector []int) []int {
	out := make([]int, len(vector))
	for i := range matrix {
		for j := range matrix[i] {
			if matrix[i][j] != 0 {
				out[i] += matrix[i][j] * vector[j]
			}
		}
	}
	return out
}

func toSparse(matrix [][]int) [][]int {
	// Count non-zero elements
	nnz := 0
	for i := range matrix {
		for j := range matrix[i] {
			if matrix[i][j] != 0 {
				nnz++
			}
		}
	}

	// Create CSR format arrays
	sparse := make([][]int, 3)
	sparse[0] = make([]int, len(matrix)+1) // row pointers
	sparse[1] = make([]int, nnz)           // column indices
	sparse[2] = make([]int, nnz)           // values

	// Fill CSR arrays
	idx := 0
	for i := range matrix {
		sparse[0][i] = idx
		for j := range matrix[i] {
			if matrix[i][j] != 0 {
				sparse[1][idx] = j
				sparse[2][idx] = matrix[i][j]
				idx++
			}
		}
	}
	sparse[0][len(matrix)] = nnz

	return sparse
}

func matrixVectorMultiplicationSparse(sparse [][]int, vector []int) []int {
	out := make([]int, len(vector))
	rowPtr := sparse[0]
	colIdx := sparse[1]
	values := sparse[2]

	// For each row
	for i := 0; i < len(vector); i++ {
		// Get range of non-zero elements for this row
		start := rowPtr[i]
		end := rowPtr[i+1]

		// Multiply and sum non-zero elements
		for j := start; j < end; j++ {
			col := colIdx[j]
			val := values[j]
			out[i] += val * vector[col]
		}
	}

	return out
}

func matrix(cells []int, neighbors [][]int, Nx, Ny, ITS int) {
	// iterate
	for i := 0; i < ITS; i++ {
		// matrix-vector multiplication
		// alive := matrixVectorMultiplication(neighbors, cells)
		alive := matrixVectorMultiplicationSparse(neighbors, cells)
		// apply rules
		for i := range cells {
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
	}
}

type WorkerResult struct {
	cells    [][]bool
	workerID int
}

func updateCellParallel(cells [][]bool, startRow, endRow int, workerID int, resultChan chan WorkerResult, Nx, Ny int) {
	// Create a new grid section for the next state
	nextCells := make([][]bool, endRow-startRow)
	for y := range nextCells {
		nextCells[y] = make([]bool, Nx)
	}

	// Process only the assigned rows
	for i := startRow; i < endRow; i++ {
		localI := i - startRow // Local index for our section
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
				nextCells[localI][j] = liveNeighbors == 2 || liveNeighbors == 3
			} else {
				// Dead cell becomes live if it has exactly 3 neighbors, otherwise it stays dead
				nextCells[localI][j] = liveNeighbors == 3
			}
		}
	}

	// Send the result back through the channel with the worker ID
	resultChan <- WorkerResult{cells: nextCells, workerID: workerID}
}

func updateCellParallel2(cells [][]bool, nextCells [][]bool, startRow, endRow int, workerID int, wg *sync.WaitGroup, Nx, Ny int) {
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

func parallel(cells [][]bool, Nx, Ny, ITS, WORKERS int) {
	// iterate
	for i := 0; i < ITS; i++ {
		// Create a channel to receive results
		resultChan := make(chan WorkerResult, WORKERS)

		// Calculate rows per worker
		rowsPerWorker := Ny / WORKERS

		// Launch workers
		for w := 0; w < WORKERS; w++ {
			startRow := w * rowsPerWorker
			endRow := startRow + rowsPerWorker
			if w == WORKERS-1 {
				endRow = Ny // Make sure the last worker processes any remaining rows
			}
			go updateCellParallel(cells, startRow, endRow, w, resultChan, Nx, Ny)
		}

		// Create a new grid for the combined result
		nextCells := make([][]bool, Ny)
		for i := range nextCells {
			nextCells[i] = make([]bool, Nx)
		}

		// Collect results from all workers
		for w := 0; w < WORKERS; w++ {
			result := <-resultChan
			workerID := result.workerID
			partialResult := result.cells
			startRow := workerID * rowsPerWorker

			// Copy the partial result to the combined grid
			for i := range partialResult {
				copy(nextCells[startRow+i], partialResult[i])
			}
		}

	}
}

func parallel2(cells [][]bool, Nx, Ny, ITS, WORKERS int) {
	// iterate
	for i := 0; i < ITS; i++ {
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
			go updateCellParallel2(cells, nextCells, startRow, endRow, w, &wg, Nx, Ny)
		}

		// Wait for all workers to complete
		wg.Wait()
	}
}

func runNaive(Nx, Ny, ITS int, runs int) int64 {
	var totalTime int64
	for i := 0; i < runs; i++ {
		fmt.Print(fmt.Sprintf("%d.", i+1))
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
		start := time.Now()
		naive(cells, Nx, Ny, ITS)
		totalTime += time.Since(start).Milliseconds()
	}
	fmt.Println()
	return totalTime / int64(runs)
}

func runPad(Nx, Ny, ITS int, runs int) int64 {
	var totalTime int64
	for i := 0; i < runs; i++ {
		fmt.Print(fmt.Sprintf("%d.", i+1))
		// initialize cells
		cells := make([][]bool, Ny+2)
		for i := range cells {
			cells[i] = make([]bool, Nx+2)
		}
		for i := 1; i <= Ny; i++ {
			for j := 1; j <= Nx; j++ {
				if rand.Float32() < 0.3 { // 30% chance of being alive
					cells[i][j] = true
				}
			}
		}
		start := time.Now()
		pad(cells, Nx, Ny, ITS)
		totalTime += time.Since(start).Milliseconds()
	}
	fmt.Println()
	return totalTime / int64(runs)
}

func runStruct(Nx, Ny, ITS int, runs int) int64 {
	var totalTime int64
	for i := 0; i < runs; i++ {
		fmt.Print(fmt.Sprintf("%d.", i+1))
		// initialize cells
		b := Board{}
		b.init(Nx, Ny)
		start := time.Now()
		structed(&b, ITS)
		totalTime += time.Since(start).Milliseconds()
	}
	fmt.Println()
	return totalTime / int64(runs)
}

func runMatrix(Nx, Ny, ITS int, runs int) int64 {
	var totalTime int64
	// Create matrix of neighbors
	// neighbors := make([][]int, (Nx+2)*(Ny+2))
	// for i := range neighbors {
	// 	neighbors[i] = make([]int, (Nx+2)*(Ny+2))
	// }
	// for i := 1; i <= Ny; i++ {
	// 	for j := 1; j <= Nx; j++ {
	// 		for di := -1; di <= 1; di++ {
	// 			for dj := -1; dj <= 1; dj++ {
	// 				if di == 0 && dj == 0 {
	// 					continue
	// 				}
	// 				neighbors[i*(Nx+2)+j][(i+di)*(Nx+2)+(j+dj)] = 1
	// 			}
	// 		}
	// 	}
	// }
	// neighbors = toSparse(neighbors)
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
	for i := 0; i < runs; i++ {
		fmt.Print(fmt.Sprintf("%d.", i+1))
		// initialize cells
		cells := make([]int, (Nx+2)*(Ny+2))
		for i := 1; i <= Ny; i++ {
			for j := 1; j <= Nx; j++ {
				if rand.Float32() < 0.3 { // 30% chance of being alive
					cells[i*(Nx+2)+j] = 1
				}
			}
		}
		start := time.Now()
		matrix(cells, neighbors, Nx, Ny, ITS)
		totalTime += time.Since(start).Milliseconds()
	}
	fmt.Println()
	return totalTime / int64(runs)
}

func runParallel(Nx, Ny, ITS, runs, WORKERS int) int64 {
	var totalTime int64
	for i := 0; i < runs; i++ {
		fmt.Print(fmt.Sprintf("%d.", i+1))
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
		start := time.Now()
		parallel(cells, Nx, Ny, ITS, WORKERS)
		totalTime += time.Since(start).Milliseconds()
	}
	fmt.Println()
	return totalTime / int64(runs)
}

func runParallel2(Nx, Ny, ITS, runs, WORKERS int) int64 {
	var totalTime int64
	for i := 0; i < runs; i++ {
		fmt.Print(fmt.Sprintf("%d.", i+1))
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
		start := time.Now()
		parallel2(cells, Nx, Ny, ITS, WORKERS)
		totalTime += time.Since(start).Milliseconds()
	}
	fmt.Println()
	return totalTime / int64(runs)
}

func main() {
	N := []int{100, 200, 500, 1000}
	ITS := 100
	runs := 10

	// Create/open CSV file
	f, err := os.Create("benchmark.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Write CSV header
	f.WriteString("size,algorithm,time_ms,iterations_per_second\n")

	for _, n := range N {
		fmt.Printf("Benchmarking size %d\n", n)

		// // Time naive implementation
		// fmt.Printf("Naive implementation\n")
		// meanTime := runNaive(n, n, ITS, runs)
		// ips := float64(ITS) / (float64(meanTime) / 1000.0)
		// fmt.Printf("Naive implementation: %d ms (%.1f iterations/s)\n", meanTime, ips)
		// // Write naive result
		// f.WriteString(fmt.Sprintf("%d,naive,%d,%.1f\n", n, meanTime, ips))

		// // Time padded implementation
		// fmt.Printf("Padded implementation\n")
		// meanPadTime := runPad(n, n, ITS, runs)
		// padIps := float64(ITS) / (float64(meanPadTime) / 1000.0)
		// fmt.Printf("Padded implementation: %d ms (%.1f iterations/s)\n", meanPadTime, padIps)
		// // Write padded result
		// f.WriteString(fmt.Sprintf("%d,padded,%d,%.1f\n", n, meanPadTime, padIps))

		// // Time structed implementation
		// fmt.Printf("Structed implementation\n")
		// meanStructedTime := runStruct(n, n, ITS, runs)
		// structedIps := float64(ITS) / (float64(meanStructedTime) / 1000.0)
		// fmt.Printf("Structed implementation: %d ms (%.1f iterations/s)\n", meanStructedTime, structedIps)
		// // Write structed result
		// f.WriteString(fmt.Sprintf("%d,structed,%d,%.1f\n", n, meanStructedTime, structedIps))

		// // Time matrix implementation (only for small grids)
		// fmt.Printf("Matrix implementation\n")
		// meanMatrixTime := runMatrix(n, n, ITS, runs)
		// matrixIps := float64(ITS) / (float64(meanMatrixTime) / 1000.0)
		// fmt.Printf("Matrix implementation: %d ms (%.1f iterations/s)\n", meanMatrixTime, matrixIps)
		// // Write matrix result
		// f.WriteString(fmt.Sprintf("%d,matrix,%d,%.1f\n", n, meanMatrixTime, matrixIps))

		// Time parallel implementation
		workers := []int{2, 4, 8}
		for _, w := range workers {
			fmt.Printf("Parallel implementation with %d workers\n", w)
			meanParallelTime := runParallel(n, n, ITS, runs, w)
			parallelIps := float64(ITS) / (float64(meanParallelTime) / 1000.0)
			fmt.Printf("Parallel implementation: %d ms (%.1f iterations/s)\n", meanParallelTime, parallelIps)
			// Write parallel result
			f.WriteString(fmt.Sprintf("%d,parallel_%d,%d,%.1f\n", n, w, meanParallelTime, parallelIps))
		}

		// Time parallel2 implementation
		for _, w := range workers {
			fmt.Printf("Parallel2 implementation with %d workers\n", w)
			meanParallelTime := runParallel2(n, n, ITS, runs, w)
			parallelIps := float64(ITS) / (float64(meanParallelTime) / 1000.0)
			fmt.Printf("Parallel2 implementation: %d ms (%.1f iterations/s)\n", meanParallelTime, parallelIps)
			// Write parallel result
			f.WriteString(fmt.Sprintf("%d,parallel2_%d,%d,%.1f\n", n, w, meanParallelTime, parallelIps))
		}
	}
}
