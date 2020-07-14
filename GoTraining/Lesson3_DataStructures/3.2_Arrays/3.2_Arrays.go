package main

import (
	"fmt"
)

func main() {
	//To reduce latency, pre-fetchers predict which data is needed next
	//"predictable access patterns"
	//Contiguous allocation of memory and iterate over it => predictable access pattern
	//data-oriented design is important, especially in Go
	//how you access your data can have very important effects
	//for this reason understanding your data is critical
	//understand input, output, and transformation
	//when practicable and reasonable, we should be using slices (for predictable access patterns)

	var strings [5]string
	strings[0] = "Apple"
	strings[1] = "Orange"
	strings[2] = "Banana"
	strings[3] = "Grape"
	strings[4] = "Plum"

	//fruit is a copy of the value in strings[i]
	//fmt.Println uses another copy of fruit
	//by copying, only the backing array is on the heap (which reduces pressure on the garbage collector)
	//"huge efficencies through our pointer sharing/pointer semantics"
	for i, fruit := range strings {
		fmt.Println(i, fruit)
	}
	
	numbers := [4]int{10,20,30,40}

	//it would be better to use range here...
	for i := 0; i<len(numbers);i++ {
		fmt.Println(i, numbers[i]) 
	}


	//cannot assign array of different sizes because their type is dependent on their size
	//[4]int is a different type from [5]int
	//arrays have a known size at compile time => cannot use a value/variable to determine size
}