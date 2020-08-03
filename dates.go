package main

import (
	"time"
	"fmt"
)

func main() {
	current := time.Now()
	fmt.Println(current.Date())
	lastWeek := current.AddDate(0,0,-7)
	fmt.Println(lastWeek.Date())
}