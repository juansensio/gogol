package main

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"time"
)

// types after variable names
// can use , if types are the same

func add(a, b int) int {
	return a + b
}

// can return multiple values

func swap(a, b string) (string, string) {
	return b, a
}

// can return named values

func multiply(a, b int) (x int) {
	x = a * b
	return // will return x
}

// can define variables at package level
// can use , if types are the same

var x, y bool
var (
	a int
	b bool
)

func main() {
	fmt.Println("Hello, World!", time.Now())
	fmt.Println(rand.Intn(10), math.Pi)

	// functions
	fmt.Println(add(1, 2))
	fmt.Println(swap("Hello", "World"))
	fmt.Println(multiply(2, 3))

	// variables
	var i int // 0, false, "" by default
	fmt.Println(i, x, y)
	var j int = 1          // initialize with value
	var z = 2              // type inferred
	var c, d = 1.2, "hola" // can declare multiple variables of different types
	fmt.Println(j, z, c, d)
	// short variable declaration
	k := 3 // type inferred (will not work outside of functions)
	fmt.Println(k)
	// type conversion
	w := float32(k) + 1.2
	fmt.Println(w)
	// constants
	const X = 1
	fmt.Println(X)
	// X = 2 // cannot change constant

	// for loops

	sum := 0
	for i := 0; i < 10; i++ {
		sum += i
	}
	fmt.Println(sum)

	sum = 0
	for i := range 10 {
		sum += i
	}
	fmt.Println(sum)

	// for continued (while)
	sum = 1
	for sum < 1000 {
		sum += sum
	}
	fmt.Println(sum)

	// inifinite
	for {
		fmt.Println("loop")
		break
	}

	// conditionals
	var x int = -1
	if x > 0 {
		fmt.Println("x is positive")
	} else if x < 0 {
		fmt.Println("x is negative")
	} else {
		fmt.Println("x is zero")
	}

	if x := 1; x > 0 { // x is scoped to if
		fmt.Println("x is positive")
	}
	fmt.Println(x)

	// switch (no break needed, go puts break after each case automatically)
	fmt.Print("Go runs on ")
	switch os := runtime.GOOS; os {
	case "darwin":
		fmt.Println("OS X.")
	case "linux":
		fmt.Println("Linux.")
	default:
		fmt.Printf("%s.\n", os)
	}
	t := time.Now()
	switch { // no need for expression (is True)
	case t.Hour() < 12:
		fmt.Println("Good morning!")
	case t.Hour() < 17:
		fmt.Println("Good afternoon.")
	default:
		fmt.Println("Good evening.")
	}

	// defer (execute at end of function, multiple defers are executed in reverse order)
	defer fmt.Println("bye")
	fmt.Println("hello")

}
