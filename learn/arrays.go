package main

import "fmt"

func main() {
	// arrays
	var a [2]string // array of strings
	a[0] = "Hello"
	a[1] = "World"
	fmt.Println(a[0], a[1])
	fmt.Println(a)

	primes := [6]int{2, 3, 5, 7, 11, 13} // array of ints
	fmt.Println(primes)

	// slices
	// arrays have a fixed size, slices are dynamic
	// does not store any data, it just describes a section
	var s []int = primes[1:4]
	fmt.Println(s)
	// Changing the elements of a slice modifies the corresponding elements of its underlying array.
	s[0] = 100
	fmt.Println(primes)
	// slice literals
	s2 := []struct {
		i int
		b bool
	}{
		{2, true},
		{3, false},
		{5, true},
		{7, true},
		{11, false},
		{13, true},
	}
	fmt.Println(s2)

	var s3 []int
	fmt.Println(s3, len(s3), cap(s3))
	if s3 == nil {
		fmt.Println("nil!")
	}
	s3 = append(s3, 1)
	fmt.Println(s3)
	s3 = append(s3, 2, 3, 4)
	fmt.Println(s3)

	// iterating
	for i, v := range s3 { // return index and value
		fmt.Printf("2**%d = %d\n", i, v)
	}
	for i := range s3 { // only index
		fmt.Println(i)
	}
	for _, v := range s3 { // omit index
		fmt.Println(v)
	}
}
