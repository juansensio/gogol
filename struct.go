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
var ITS int = 100

type Board struct {
	cells  [][]bool
	_cells [][]bool
	w      int
	h      int
}

// implement String interface

func (b *Board) String() string {
	s := "\033[H\033[2J"
	for i := range b.cells {
		for j := range b.cells[i] {
			if b.cells[i][j] {
				s += "@"
			} else {
				s += " "
			}
		}
		s += "\n"
	}
	return s
}

// implement image.Image interface

func (b *Board) ColorModel() color.Model {
	return color.GrayModel
}

func (b *Board) Bounds() image.Rectangle {
	return image.Rect(0, 0, b.w-2, b.h-2)
}

func (b *Board) At(x, y int) color.Color {
	if b.cells[y+1][x+1] {
		return color.Gray{255}
	}
	return color.Gray{0}
}

// methods

func (b *Board) init() {
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
	for i := 1; i <= Ny; i++ {
		for j := 1; j <= Nx; j++ {
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
	for i := 1; i <= Ny+1; i++ {
		for j := 1; j <= Nx+1; j++ {
			b._cells[i][j] = b.cells[i][j]
		}
	}
}

func (b *Board) save(i int) {
	f, err := os.Create(fmt.Sprintf("out/%d.png", i))
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer f.Close()
	if err := png.Encode(f, b); err != nil { // this works because Board implements image.Image
		fmt.Println("Error encoding image:", err)
	}
}

func main() {

	// delete out folder and create it again
	os.RemoveAll("out")
	os.Mkdir("out", 0755)

	// initialize cells
	b := Board{}
	b.init()

	// iterate
	for i := 0; i < ITS; i++ {
		b.update()
		time.Sleep(100 * time.Millisecond)
		fmt.Println(&b) // this works because Board implements String()
		b.save(i)
	}
}
