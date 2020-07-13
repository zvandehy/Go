package main

import (
	"fmt"
)

func main() {
	//constants only exist at compile time
	//considered mathematically exact

	//untyped constants can be implicitly converted
	const ui = 12345 //kind: integer
	const uf = 3.141592 //kind: floating-point

	//typed constants cannot be implicitly converted
	const ti int = 12345
	const tf float64 = 3.141592

	//untyped constants allow this
	const answer = 3 * 0.333 //answer will be floating point

	const third = 1 / 3.0 //third will be floating point

	const zero = 1 / 3 //will be int (0)

	const (
		A = iota
		B
		C
	)

	fmt.Println("iota: ", A, B, C)




}



