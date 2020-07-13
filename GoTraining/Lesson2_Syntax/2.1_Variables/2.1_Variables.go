package main

import (
	"fmt"
)

func main() {
	
	//go word size matches address size (and thus pointer size)
	//we are on a 64 bit OS so address size is 64 bit

	//zero value: every value must be initialized, so if not specified it is set to the 0 value
	
	//"go has too many ways to declare variables" haha

	//when creating a variable that will be set to zero value, use "var"
	var a int
	//strings are "two word" datastructures: backing array, length
	//zero value of string: nil, 0: "empty"
	var b string
	var c float64
	var d bool

	fmt.Printf("var a int \t %T [%v]\n",a,a)
	fmt.Printf("var b string \t %T [%v]\n",b,b)
	fmt.Printf("var c float64 \t %T [%v]\n",c,c)
	fmt.Printf("var d bool \t %T [%v]\n",d,d)

	//Declare and initialize (not zero value)
	aa := 10
	//"two word" => backing array: [hello], length: 5
	bb := "hello"
	cc := 3.14159
	dd := true

	fmt.Printf("aa := 10 \t %T [%v]\n",aa,aa)
	fmt.Printf("bb := \"hello\" \t %T [%v]\n",bb,bb)
	fmt.Printf("cc := 3.14159 \t %T [%v]\n",cc,cc)
	fmt.Printf("dd := true \t %T [%v]\n",dd,dd)

	//go doesn't have casting
	//instead it has conversion

	aaa := int32(10)
	fmt.Printf("aaa := int32(10) \t %T [%v]\n",aaa,aaa)
}