package main

//think of packages as "self-contained unit of code"

//when an identifier starts with a capital letter
//  it is considered "exported"

//packages are relative paths from your go path

import "Programming/Go/GoTraining/Lesson4_Decoupling/4.4_Exporting/users"

func main() {

	u := user{
		Name:"Zeke",
		Email:"zvandehy@gmail.com"
		//password: "qwerty1234" // not exported!
	}
}