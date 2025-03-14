package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

func naive(Nx, Ny, ITS int) {
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

func pad(Nx, Ny, ITS int) {
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

func structed(Nx, Ny, ITS int) {
	b := Board{}
	b.init(Nx, Ny)
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

func matrix(Nx, Ny, ITS int) {
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
	neighbors := make([][]int, (Nx+2)*(Ny+2))
	for i := range neighbors {
		neighbors[i] = make([]int, (Nx+2)*(Ny+2))
	}
	for i := 1; i <= Ny; i++ {
		for j := 1; j <= Nx; j++ {
			for di := -1; di <= 1; di++ {
				for dj := -1; dj <= 1; dj++ {
					if di == 0 && dj == 0 {
						continue
					}
					neighbors[i*(Nx+2)+j][(i+di)*(Nx+2)+(j+dj)] = 1
				}
			}
		}
	}
	neighbors = toSparse(neighbors)
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

func run(Nx, Ny, ITS int, runs int, f func(Nx, Ny, ITS int)) int64 {
	var totalTime int64
	for i := 0; i < runs; i++ {
		fmt.Print(fmt.Sprintf("%d.", i+1))
		start := time.Now()
		f(Nx, Ny, ITS)
		totalTime += time.Since(start).Milliseconds()
	}
	fmt.Println()
	return totalTime / int64(runs)
}

func main() {
	N := []int{100, 500, 1000}
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

		// Time naive implementation
		fmt.Printf("Naive implementation\n")
		meanTime := run(n, n, ITS, runs, naive)
		ips := float64(ITS) / (float64(meanTime) / 1000.0)
		fmt.Printf("Naive implementation: %d ms (%.1f iterations/s)\n", meanTime, ips)
		// Write naive result
		f.WriteString(fmt.Sprintf("%d,naive,%d,%.1f\n", n, meanTime, ips))

		// Time padded implementation
		fmt.Printf("Padded implementation\n")
		meanPadTime := run(n, n, ITS, runs, pad)
		padIps := float64(ITS) / (float64(meanPadTime) / 1000.0)
		fmt.Printf("Padded implementation: %d ms (%.1f iterations/s)\n", meanPadTime, padIps)
		// Write padded result
		f.WriteString(fmt.Sprintf("%d,padded,%d,%.1f\n", n, meanPadTime, padIps))

		// Time structed implementation
		fmt.Printf("Structed implementation\n")
		meanStructedTime := run(n, n, ITS, runs, structed)
		structedIps := float64(ITS) / (float64(meanStructedTime) / 1000.0)
		fmt.Printf("Structed implementation: %d ms (%.1f iterations/s)\n", meanStructedTime, structedIps)
		// Write structed result
		f.WriteString(fmt.Sprintf("%d,structed,%d,%.1f\n", n, meanStructedTime, structedIps))

		// Time matrix implementation
		fmt.Printf("Matrix implementation\n")
		meanMatrixTime := run(n, n, ITS, runs, matrix)
		matrixIps := float64(ITS) / (float64(meanMatrixTime) / 1000.0)
		fmt.Printf("Matrix implementation: %d ms (%.1f iterations/s)\n", meanMatrixTime, matrixIps)
		// Write matrix result
		f.WriteString(fmt.Sprintf("%d,matrix,%d,%.1f\n", n, meanMatrixTime, matrixIps))
	}
}
