package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

var Nx int = 160
var Ny int = 40
var ITS int = 1000

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

func main() {

	// delete out folder and create it again
	os.RemoveAll("out")
	os.Mkdir("out", 0755)

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

	// create neighbors sparse matrix directly in CSR format
	// calculate number of non-zero elements first
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

	// print sparse matrix
	// fmt.Println("Sparse Matrix:")
	// fmt.Println("Row pointers:", neighbors[0])
	// fmt.Println("Column indices:", neighbors[1])
	// fmt.Println("Values:", neighbors[2])

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

		// wait
		time.Sleep(50 * time.Millisecond)

		// print cells
		printCellsTerminal(i, cells)
	}
}
