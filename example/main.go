package main

import (
	"fmt"

	sc "github.com/imatakatsu/simple-colors"
)

func main() {
	fmt.Println(sc.Gradient(false, "just beautiful example",
		sc.Rgb(195, 232, 14),
		sc.Rgb(68, 216, 72),
		sc.Rgb(10, 155, 162),
	))
}
