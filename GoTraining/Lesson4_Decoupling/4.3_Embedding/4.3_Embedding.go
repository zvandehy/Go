package main

import "fmt"

type user struct {
	name string
	email string
}

//user type creates notify method
func (u *user) notify() {
	fmt.Printf("My name is %q and my email is %q\n", u.name, u.email)
}

type admin struct {
	user //embed the user struct into admin
	priority string
}

func main() {
	a := admin {
		user: user{ //create user
			name: "Zeke",
			email: "zvandehy@gmail.com",
		},
		priority: "High",
	}

	a.notify() //a can call the notify method because user is embedded

	a.user.notify() //also works to access the user

	//if admin implemented the notify method, then the user's implementation would not be promoted (would use admin's version)

}