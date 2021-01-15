[![Actions Status](https://github.com/neilotoole/errgroup/workflows/Go/badge.svg)](https://github.com/neilotoole/errgroup/actions?query=workflow%3AGo)
[![Go Report Card](https://goreportcard.com/badge/StevenACoffman/errgroup)](https://goreportcard.com/report/StevenACoffman/errgroup)
[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://pkg.go.dev/github.com/StevenACoffman/errgroup)
[![license](https://img.shields.io/github/license/StevenACoffman/errgroup)](./LICENSE)

# StevenACoffman/errgroup
`StevenACoffman/errgroup` is a drop-in alternative to Go's wonderful
[`sync/errgroup`](https://pkg.go.dev/golang.org/x/sync/errgroup) but it converts goroutine panics to errors. 

While `net/http` installs a panic handler with each request-serving goroutine,
goroutines **do not** and **cannot** inherit panic handlers from parent goroutines,
so a `panic()` in one of the child goroutines will kill the whole program.

So whenever you use an `sync.errgroup`, with some discipline, you can always remember to add a
deferred `recover()` to every goroutine.  This library just avoids that boilerplate and does that for you.

You can [see it in use](https://play.golang.org/p/S8Gmr_sWZIi)

```go
package main

import (
	"fmt"

	"github.com/StevenACoffman/errgroup"
)

func main() {
	g := new(errgroup.Group)
	var urls = []string{
		"http://www.golang.org/",
		"http://www.google.com/",
		"http://www.somestupidname.com/",
	}
	for i := range urls {
		// Launch a goroutine to fetch the URL.
		i := i // https://golang.org/doc/faq#closures_and_goroutines
		g.Go(func() error {

			// deliberate index out of bounds triggered
			fmt.Println("Fetching:", i, urls[i+1])

			return nil
		})
	}
	// Wait for all HTTP fetches to complete.
	err := g.Wait()
	if err == nil {
		fmt.Println("Successfully fetched all URLs.")
	} else {
		fmt.Println(err)
	}
}
```

This work was done by my co-worker [Ben Kraft](https://github.com/benjaminjkraft), and, with his permission, I lightly modified it to
lift it out of our repository for Go community discussion.

### Counterpoint
There is [an interesting discussion](https://github.com/oklog/run/issues/10) which has an alternative view that,
with few exceptions, panics **should** crash your program.

### Prior Art
With only a cursory search, I found a few existing open source examples.

#### [Kratos](https://github.com/go-kratos/kratos errgroup 

Kratos Go framework for microservices has a similar [errgroup](https://github.com/go-kratos/kratos/blob/master/pkg/sync/errgroup/errgroup.go)
solution.

#### PanicGroup by Sergey Alexandrovich

In the article [Errors in Go:
From denial to acceptance](https://evilmartians.com/chronicles/errors-in-go-from-denial-to-acceptance), 
(which advocates panic based flow control 😱), they have a PanicGroup that's roughly equivalent:

```
type PanicGroup struct {
  wg      sync.WaitGroup
  errOnce sync.Once
  err     error
}

func (g *PanicGroup) Wait() error {
  g.wg.Wait()
  return g.err
}

func (g *PanicGroup) Go(f func()) {
  g.wg.Add(1)

  go func() {
    defer g.wg.Done()
    defer func(){
      if r := recover(); r != nil {
        if err, ok := r.(error); ok {
          // We need only the first error, sync.Once is useful here.
          g.errOnce.Do(func() {
            g.err = err
          })
        } else {
          panic(r)
        }
      }
    }()

    f()
  }()
}
```
