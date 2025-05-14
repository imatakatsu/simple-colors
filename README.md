# simple-colors
simple lib for gradients in golang. support rgb colors, but convert them to hsl to prevent dirty gradients. also has flexmode to make gradients more beautiful

## Installation
```bash
go get github.com/imatakatsu/simple-colors
```
## Usage example

```go
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
```

## Features
support more than 2 colors
flex mode for better gradients
hsl gradients to prevent dirty gradients
out string with ANSI color codes

made by me
tg: https://t.me/tokyoddos
