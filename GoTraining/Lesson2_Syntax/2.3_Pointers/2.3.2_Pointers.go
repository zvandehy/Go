package main

import ("fmt")

func main() {

	count := 10

	//display value and address
	fmt.Println("Before: ",count,&count)

	incrementWithAddress(&count) //passing an address VALUE across program boundary
	
	fmt.Println("After : ",count,&count)
}

func incrementWithAddress(inc *int) {
	//receive the address value across program boundary
	//*int purpose is store store ADDRESS value (32 or 64 bit address)
	//*int is the address of count => &count
	*inc++//"*inc" gets the value at the address saved in "inc"
	//manipulates memory outside of this program boundary (outside of the allocated memory space in the stack for this function)

	//value that the pointer points to, memory address of inc, memory address that it points to
	fmt.Println("Inc   : ",*inc, &inc, inc)
}

//escape analysis decides if something is stored on the stack or the heap
//stack memory cannot be shared between two go routines
//pacing algorithm runs garbage collector to maintain a small heap
//go has made GC incredibly efficient
