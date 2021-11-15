# learn how to use webfr:

- webfr is webframework that runs on any devices

# Hello world:

- first program using webfr

```golang
// package should be main
package main

// import the module 
import (
    "github.com/krishpranav/webfr"
)

func main() {
    // new application
    wb := webfr.New()

    // path
    wb.Get("/hello", func(ctx webfr.Context) {
        ctx.SendString("Helo Friends!!")
    })

    // start the server
    wb.Start(":3000")
}
```

- command to run the server:
```
go run main.go
```

- navigate to ```127.0.0.1/hello```

# Api:

- New:

- creates a new instance of webfr with optional settings

```golang
New(settings ...*Settings) Webfr
```

- Settings:

- pass application settings while calling new:
```golang

package main

import (
    "github.com/krishpranav/webfr"
)

func main() {
    // Setup webfr
    wb := webfr.New(&webfr.Settings{
        CaseInSensitive: true
        ServerName:      "webfr"
    })


    // start the server
    g.Start(":3000")
}
```

# Http Methods:

## Get

- The GET Method Request:

```golang
Get(path string, handlers ...handlerFunc) *Route
```

## Head

- The Head Method:

```golang
Head(path string, handlers ...handlerFunc) *Route
```

## Post

- The Post Method:

```golang
Post(path string, handlers ...handlerFunc) *Route
```

## Put

- The Put Method:

```golang
Put(path string, handlers ...handlerFunc) *Route
```

## Delete

- The Delete Method:

```golang
Delete(path string, handlers ...handlerFunc) *Route
```

## Connect

- The Connect Method:

```golang
Connect(path string, handlers ...handlerFunc) *Route
```

## Options

- The Options Method:

```golang
Options(path string, handlers ...handlerFunc) *Route
```

## Trace

- The Trace Method:

```golang
Trace(path string, handlers ...handlerFunc) *Route
```

## Patch

- The Patch Method:

```golang
Patch(path string, handlers ...handlerFunc) *Route
```

## Not Found

- The NotFound Method:

```golang
NotFound(handlers ...handlerFunc)
```

## Use

- Middleware Use Method:

```golang
Use(middlewares ...handlerFunc)
```

## Static

- Servers static files in a root directory under specific prefix:

```golang
Static(prefix, root string)
```

- Example

```golang
// Serve files in assets directory for prefix static
// requests will be like 
// http://localhost:3000/static/test.png
wb.Static("/static", "./assets")
```

## Group

- Group registers routes under specific prefix

```golang
Group(prefix string, routes []*Route) []*Route
```

- Example:

```golang
// handler request for /account/id
wb.Group("/account",  []*webfr.Route{
    wb.Get("/id", func(ctx webfr.Context) {
        ctx.SendString("User X")
    })
})
```
