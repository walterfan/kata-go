# kata-go

Kata should be called "routine" in Chinese. The secret of practicing martial arts is to master various routines.

First, you must learn from the strengths of hundreds of schools and be familiar with the routines of various schools in the world. Only then can you integrate them and achieve the state of no moves being better than moves.

Dave Thomas - the author of "The Pragmatic Programmer", proposed the idea of ​​Code Kata. Dave also collected a small practice project on his website (http://codekata.com/).

As a professional programmer, I hope to practice some routines that can be often used in work, such as some small routines for file modification, image cutting, and network sending and receiving, so I will organize and collect some routines here.


## Kata

1. [cron-service.go](./kata/cron): demostration of how to implement cron job
1. [encoding_tool](./kata/encoding_tool): demostrate how to implement encoding command line tool by cobra library
1. [ds_web_console](./kata/ds_web_console/): demostrate how to call deep seek api and provide a web console by gin library
1. [ata-auth](./kata/kata-auth/): demostrate how to implement a simple auth system by casbin
1. [list_files.go](./kata/kata-files/list_files.go): demostrate how to list files in a directory
1. [links.go](./kata/kata-http/links.go): demostrate how to use http client to get links from a web page
1. [llm-agent-go](./kata/llm-agent-go/): demostrate how to implement a simple llm agent by gin and vue.js
1. [prompt_service](./kata/prompt_service): demostrate how to implement a simple prompt management service by gin and gorm with sqlite db
2. [prompt_service_v2](./kata/prompt_service_v2): add login and metrics endpoint for prompt management service
3. [service-monitor](./kata/service_monitor): demostrate how to monitor a backend service with prometheus exportor
4. [simple-ai-agent](./kata/simple-ai-agent): demostrate how to use function calling and tools of openai
5. [unix_socket](./kata/unix_socket): demostrate how to use unix socket to communicate with backend service


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
### unit test

```
# Run all tests with verbose output for entire project
go test -v ./...

# Run tests in specific package with verbose output
go test -v ./internal/cmd

# Run specific test function
go test -v ./internal/cmd -run Test_Metrics

# Run tests matching a pattern
go test -v ./internal/cmd -run "Test.*eks.*"

# Run tests with regex pattern
go test -v ./internal/cmd -run "^TestMonitor.*"

# Skip specific tests
go test -v ./internal/cmd -skip "TestFileSync"

# Skip multiple test patterns
go test -v ./internal/cmd -skip "TestFileSync|TestMetrics"

# Run tests but skip long-running ones
go test -v ./internal/cmd -short
```
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

## Reference
* [Go Official Site](https://go.dev/)
* [Go Documentation](https://pkg.go.dev/)
* [Go Playground](https://go.dev/play/)
* [Go Blog](https://blog.golang.org/)
* [Go Wiki](https://github.com/golang/go/wiki)
* [Effective Go](https://golang.org/doc/effective_go)
* [Go Memory Model](https://golang.org/ref/mem)
* [Go Standards](https://google.github.io/styleguide/go/guide)
* [Awesome go projects](https://github.com/avelino/awesome-go)
