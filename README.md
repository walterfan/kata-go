# kata-go

Kata should be called "routine" in Chinese. The secret of practicing martial arts is to master various routines.

First, you must learn from the strengths of hundreds of schools and be familiar with the routines of various schools in the world. Only then can you integrate them and achieve the state of no moves being better than moves.

Dave Thomas - the author of "The Pragmatic Programmer", proposed the idea of ​​Code Kata. Dave also collected a small practice project on his website (http://codekata.com/).

As a professional programmer, I hope to practice some routines that can be often used in work, such as some small routines for file modification, image cutting, and network sending and receiving, so I will organize and collect some routines here.

## Cheat sheet of golang

* [Go cheatsheet 1](go-cheat-sheet.md)
* [Go cheatsheet 2 ](https://devhints.io/go)
* [Go cheatsheet 3](https://quickref.me/go.html)

## Go tutorial

* [Go tutorial](https://tour.golang.org/welcome/1)
* [Go by example](https://gobyexample.com/)
* [Go Style Guide](https://google.github.io/styleguide/go/guide)
* [Go Style Decisions](https://google.github.io/styleguide/go/decisions)
* [Go Style Best Practices](https://google.github.io/styleguide/go/best-practices)

## Go Tools
### build and run
```shell
go build xxx.go
go run xxx.go
```
### check dependency
go list, go get, go mod xxx

```
go mod init my_project
go mod tidy
```
### format code
go fmt, gofmt

### Debug
```bash
dlv debug main.go

```
### documentation
go doc, godoc

```shell
go install github.com/swaggo/swag/cmd/swag@latest
swag init

```
### ****unit test
go test
### static analysis
go vet, golangci-lint

e.g.
  ```bash
  brew install golangci-lint
  golangci-lint run ./...
  ```
### performance profile
go tool pprof, go tool trace

```bash
go tool pprof http://localhost:6060/debug/pprof/profile

go run main.go
go tool trace trace.out
```
### upgrade
go tool fix
### bug report
go bug


## example

* [unix_socket](./kata/unix_socket)
* [cron-service.go](./kata/cron)
* [list_files.go](./kata/files/list_files.go)
* [links.go](./kata/http/links.go)

## go with vscode
* install go extensions and [delve](https://github.com/go-delve/delve/blob/master/Documentation/installation/osx/install.md)
* configuration of go debug in vscode

```json
 {
    // Use IntelliSense to learn about possible attributes.
    // Hover to view descriptions of existing attributes.
    // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Debug file",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${file}"
        },
        {
            "name": "Launch Package",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${fileDirname}"
        }
    ]
}
```
