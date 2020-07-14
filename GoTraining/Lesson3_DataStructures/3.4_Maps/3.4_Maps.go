package main

import (
	"fmt"
)

type user struct {
	name string
	surname string
}

func main() {
	//key:string, value:user
	users := make(map[string]user)

	users["Roy"] = user{"Rob","Roy"}
	users["Ford"] = user{"Henry", "Ford"}
	users["Mouse"] = user{"Mickey", "Mouse"}
	users["Jackson"] = user{"Michael", "Jackson"}

	//cannot guarantee the same order!
	for key, value := range users {
		fmt.Println(key, value)
	}

	//same map
	users2 := map[string]user{
		"Roy":{"Rob","Roy"},
		"Ford":{"Henry", "Ford"},
		"Mouse":{"Mickey", "Mouse"},
		"Jackson":{"Michael", "Jackson"},
	}

	fmt.Println(users2)


	//delete Roy key
	delete(users, "Roy")

	//cannot find Roy

	u, found := users["Roy"]
	fmt.Println("Roy",found,u)



}