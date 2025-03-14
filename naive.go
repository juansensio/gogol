package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"time"
)

var Nx int = 160
var Ny int = 40
var ITS int = 10

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

func saveImage(i int, cells [][]bool) {
	img := image.NewGray(image.Rect(0, 0, Nx, Ny))
	for i := range cells {
		for j := range cells[i] {
			if cells[i][j] {
				img.Set(j, i, color.White)
			} else {
				img.Set(j, i, color.Black)
			}
		}
	}
	f, err := os.Create(fmt.Sprintf("out/%d.png", i))
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		fmt.Println("Error encoding image:", err)
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
		// Live cell continues to live if it has 2 or 3 neighbors, otherwise it dies
		return liveNeighbors == 2 || liveNeighbors == 3
	} else {
		// Dead cell becomes live if it has exactly 3 neighbors, otherwise it stays dead
		return liveNeighbors == 3
	}
}

func main() {

	// delete out folder and create it again
	os.RemoveAll("out")
	os.Mkdir("out", 0755)

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

		// save image
		saveImage(i, cells)
	}
}
