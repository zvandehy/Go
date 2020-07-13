package main

import ("fmt")

type example struct {
	flag bool //1 byte
	counter int16 //2 bytes
	pi float32 //4 bytes
}

//only better in terms of memory profile!
//may not be better in terms of readability!!
//First organize by Correctness, this example is only better for memory!!
type otherExample struct {
	pi float32 //4 bytes
	counter int16 //2 bytes
	flag bool //1 byte
}

func main() {
	//Struct zero value
	var example1 example //represents 8 bytes (7 bytes + 1 padding byte)
	//padding is from alignment boundaries: "every x byte value must lie on an x byte boundary (2 byte lies on 2)"
	//to reduce excess padding, put largest size fields in struct first

	fmt.Printf("%+v\n",example1)

	//Struct literal
	example2 := example {
		flag: true,
		counter: 10,
		pi: 3.141592,
	}

	fmt.Println("Flag",example2.flag)
	fmt.Println("Counter",example2.counter)
	fmt.Println("Pi",example2.pi)

	//not everything has to be named (prevents pollution)
	//Anonymous type
	e := struct {
		value int
		flag bool
	} {
		value: 10,
		flag: true,
	}

	fmt.Println("Flag",e.flag)
	fmt.Println("Value",e.value)

	//must use explicit conversion with named types,
	//but not necessary with anonymous types

}