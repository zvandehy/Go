package main

import (
	"fmt"
)

// user defines a user in the program.
type user struct {
	name  string
	email string
}

//these are methods because of the (u user) syntax (I believe this is called "receiver")
//methods operate on the data provided (u user)
//functions do not have that, they simply have the function name followed by the return type
//functions don't hide the cost of things as much, so they are generally preferred but not always

//notify implements a method with a value receiver.
//  the user struct is copied
func (u user) notify() {
	fmt.Printf("Sending User Email To %s<%s>\n",
		u.name,
		u.email)
}

//changeEmail implements a method with a pointer receiver.
//  only the pointer is copied (cheap)
func (u *user) changeEmail(email string) {
	u.email = email
}


func main() {
	//when creating a type, we must decide if it should use pointer semantics or value semantics
	//  if the type should be shared, then it should use pointer semantics
	//value semantics always pass a copy of the data
	//pointer semantics pass a copy of the pointer

	//general guideline: if using built-in types like numerics, string, bool,
	//  then use value semantics
	//if using reference types like slices, maps, channels, 
	//  then value semantics (expections on marshal functions)

	//if not sure, safer to use pointer semantics

	bill := user{"bill", "bill@gmail.com"}
	fmt.Println(bill)
	bill.notify()
	bill.changeEmail("bill_new@gmail.com") //go is smart enough to work with pointers or values
	bill.notify()

	fmt.Println("==========")

	f1 := bill.notify 
	//f1 is a pointer to the method
	//f1 has a copy of bill BECAUSE notify uses value semantics
	bill.name = "lisa"
	f1() //sends to bill NOT lisa
	bill.name = "bill" //change name back to bill

	f2 := bill.changeEmail
	//f2 is a pointer to the method
	//f2 points to bill, so it isn't using a copy
	bill.name = "charles"
	f2("bill@gmail.com") //will send to charles with this email
	bill.notify()
}