package main

import "fmt"


func main()  {
	//f()
	//fmt.Println("return normally from f")
	test()
}

func test()  {
	defer func() {
		fmt.Println("defer 1 in test")
	}()

	defer func() {
		fmt.Println("defer 2 in test")

		e := recover()

		fmt.Println(e)
	}()

	panic("panic in test")

	fmt.Println("print in test")
}

func f()  {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recover in f", r)
		}
	}()

	fmt.Println("calling g")
	g(0)
	fmt.Println("returned normally from g")
}

func g(i int)  {
	if i > 3 {
		fmt.Println("panicking")
		panic(fmt.Sprint("%v", i))
	}

	defer fmt.Println("defer in g", i)

	fmt.Println("printing in g", i)
	g(i + 1)
}

