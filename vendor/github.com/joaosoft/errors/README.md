# errors
[![Build Status](https://travis-ci.org/joaosoft/errors.svg?branch=master)](https://travis-ci.org/joaosoft/errors) | [![codecov](https://codecov.io/gh/joaosoft/errors/branch/master/graph/badge.svg)](https://codecov.io/gh/joaosoft/errors) | [![Go Report Card](https://goreportcard.com/badge/github.com/joaosoft/errors)](https://goreportcard.com/report/github.com/joaosoft/errors) | [![GoDoc](https://godoc.org/github.com/joaosoft/errors?status.svg)](https://godoc.org/github.com/joaosoft/errors)

Error handling with caused-by and stack.

###### If i miss something or you have something interesting, please be part of this project. Let me know! My contact is at the end.

## Dependecy Management
>### Dependency

Project dependencies are managed using Dep. Read more about [Dep](https://github.com/golang/dep).
* Get dependency manager: `go get github.com/joaosoft/dependency`
* Install dependencies: `dependency get`

>### Go
```
go get github.com/joaosoft/errors
```

## Usage 
This examples are available in the project at [examples/main.go](https://github.com/joaosoft/errors/tree/master/example_test.go)
```go
var (
	ErrorOne = errors.New(errors.LevelError, 1, "Error one")
	ErrorTwo = errors.New(errors.LevelError, 2, "Error two")
)

func main() {
	fmt.Println("\nADDING ERRORS!\n")

	errs := errors.Add(ErrorOne).
		Add(ErrorTwo)

	fmt.Println(errs.Cause())

	fmt.Println(errs.Stack)
	fmt.Println(errs.Previous.Stack)

	fmt.Println("\nDONE!")
}
```

## Known issues


## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
