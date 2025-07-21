# nbjobqueue - Non-blocking unbounded lock-free job queue for Golang
[![GoDoc](https://godoc.org/github.com/rrgmc/nbjobqueue?status.png)](https://godoc.org/github.com/rrgmc/nbjobqueue)

`nbjobqueue` is a non-blocking unbounded lock-free job queue for Golang.

As it is unbounded, care must be taken to avoid memory exhaustion if adding faster than reading. 

## Install

```shell
go get github.com/rrgmc/nbjobqueue
```

# License

MIT

### Author

Rangel Reale (rangelreale@gmail.com)
