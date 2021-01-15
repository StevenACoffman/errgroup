[![Actions Status](https://github.com/neilotoole/errgroup/workflows/Go/badge.svg)](https://github.com/neilotoole/errgroup/actions?query=workflow%3AGo)
[![Go Report Card](https://goreportcard.com/badge/StevenACoffman/errgroup)](https://goreportcard.com/report/StevenACoffman/errgroup)
[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://pkg.go.dev/github.com/StevenACoffman/errgroup)
[![license](https://img.shields.io/github/license/StevenACoffman/errgroup)](./LICENSE)

# StevenACoffman/errgroup
`StevenACoffman/errgroup` is a drop-in alternative to Go's wonderful
[`sync/errgroup`](https://pkg.go.dev/golang.org/x/sync/errgroup) but it converts goroutine panics to errors.

This work was done by my co-worker Ben Kraft, and, with his permission, I lightly modified it to
lift it out of our repository for Go community discussion.

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
