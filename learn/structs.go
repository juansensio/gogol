package main

import "fmt"

type Vertex struct { // collection of values
	X int
	Y int
}

var (
	v1 = Vertex{1, 2}  // has type Vertex
	v2 = Vertex{X: 1}  // Y:0 is implicit
	v3 = Vertex{}      // X:0 and Y:0
	p  = &Vertex{1, 2} // has type *Vertex
)

func main() {
	fmt.Println(Vertex{1, 2})

	v := Vertex{1, 2}
	v.X = 4 // can access fields with a dot
	fmt.Println(v.X)

	p := &v   // p is a pointer to v
	p.X = 1e9 // is the same as (*p).X
	fmt.Println(v)
}
