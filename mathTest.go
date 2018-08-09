package main

import "fmt"
import (
	"math/rand"
)

type Status uint

func main() {
	var a [4][4]int

	a = [4][4]int{{1, 2, 3, 4}}

	fmt.Println(rand.Intn(len(a)))

	fmt.Println(5 % 2)

	fmt.Println(2 << 1)


	const (
		Win  Status = iota
		Lose
		Add
		Max  = 2048
	)

	fmt.Println(Win, Lose, Add, Max)
}
