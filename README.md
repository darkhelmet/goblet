# goblet

Embedded assets for Go

## Usage

Basics

    go get github.com/darkhelmet/goblet/gobletize

    gobletize -h
    gobletize -input app -gobletize public

For HTTP stuff

    import H "github.com/darkhelmet/goblet/http"

    // Serve the /static/ tree, fallback to the public directory
    // for ease of development.
    http.Handle("/static/", H.LoadHandlerWithFallback("/static", "public"))

Non-HTTP

    import "github.com/darkhelmet/goblet"

    g, err := goblet.Load()
    // handle err

    asset := g.Get("path.txt")
    asset = g.Get("non-existent") // asset == nil 
