package main

import ("fmt")

func main() {
	//pass by value => what you see is what you get
	//  something is always copied, either the value or its address
	//  java is pass by reference ("this")
	//pointers serve one purpose: "sharing"
	//sharing a value across a program boundary
	//  like between function calls or between Go routine (kind of like a thread)
	//if we are sharing a value we need a pointer

	count := 10

	//display value and address
	fmt.Println("Before: ",count,&count)

	increment(count) //passing a value across program boundary
	
	fmt.Println("After : ",count,&count)
}

func increment(inc int) {
	//receive the value across program boundary
	inc++
	//notice different memory address
	fmt.Println("Inc   : ",inc,&inc)
}
