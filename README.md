# godzilla

[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com)

- A powerful go web framework for highly scalable and resource efficient web application

# Installation:
```
go get -u github.com/godzillaframework/godzilla
```

# Examples:
```golang
package main

import (
    "github.com/godzillaframework/godzilla"
)

func main() {
    godz := godzilla.New()

    wb.Get("/hello", func(ctx godz.Context) {
        ctx.SendString("Helo Friends!!")
    })

    wb.Start(":3000")
}
```

- for more tutorials visit the [docs](https://github.com/godzillaframework/godzilla/blob/master/docs/learnwebfr.md)