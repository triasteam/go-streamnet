package main

import (
	"os"
	"fmt"
)

func main() {
	fmt.Println(os.Args)
	fmt.Printf("%T - %d\n", os.Args, len(os.Args))
	fmt.Println(os.Args[1:])
}
