package main

import (
	"fmt"
)

// user defines a user in the program.
type user struct {
	name  string
	email string
}

//type assertions

//String() implements the fmt.Stringer interface
func (u *user) String() string {
	return fmt.Sprintf("My name is %q and my email is %q", u.name, u.email)
}

func main() {
	user := user{
		"bill","bill@gmail.com",
	}
	fmt.Println(user)//using value semantics which ISN'T how the interface defines the contract
	//  so it just prints it as normal
	fmt.Println(&user)//using pointer semantics, so the print is overrided
}