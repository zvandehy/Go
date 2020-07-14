package main

import (
	"fmt"
)

func main() {
	//make() only works with slices, maps, and channels
	//3 word data structure:
	//  pointer to the backing array of strings, length, capacity
	//if there is a value in the [] like [x] => array
	//if there isn't a value like [] => slice

	//length and capacity of 5
	slice := make([]string, 5)
	slice[0] = "Apple"
	slice[1] = "Orange"
	slice[2] = "Banana"
	slice[3] = "Grape"
	slice[4] = "Plum"

	//cannot access slice[5]

	//uses value of slice, so Print will copy
	fmt.Println("slice: ", slice)

	//length of 5, capacity of 8
	slice2 := make([]string, 5, 8)
	fmt.Println("make(): ", slice2)

	//zero value: backing array: nil, length: 0, capacity: 0
	var slice0 []string
	fmt.Println("zero value: ", slice0, "\n")

	slice3 := append(slice,"Added")

	slice4 := slice[:]

	fmt.Println("Before")
	fmt.Println("sliceA: ", slice)
	fmt.Println("sliceB: ", slice3)
	fmt.Println("sliceC: ", slice4, "\n")

	slice3[2] = "CHANGED_B"
	slice4[3] = "CHANGED_C"

	fmt.Println("After")
	fmt.Println("sliceA: ", slice)
	fmt.Println("sliceB: ", slice3)
	fmt.Println("sliceC: ", slice4, "\n")


	//notice that append (slice B) creates a copy of the backing array
	//  0th element has a different address from the other slices
	//every slice header has a unique address


	//%p of a slice prints the address of the 0th element
	//%p of a &slice prints the address of the slice header

	fmt.Println("Slice Header Addresses")
	fmt.Printf("sliceA: \t\t%p\n", &slice)
	fmt.Printf("sliceB: \t\t%p\n", &slice3)
	fmt.Printf("sliceC: \t\t%p\n\n", &slice4)
	fmt.Println("Slice 0th Element Address")
	fmt.Printf("sliceA 0th element: \t%p\n", slice)
	fmt.Printf("sliceB 0th element: \t%p\n", slice3)
	fmt.Printf("sliceC 0th element: \t%p\n", slice4)
}