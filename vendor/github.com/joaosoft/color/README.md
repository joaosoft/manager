# color
[![Build Status](https://travis-ci.org/joaosoft/color.svg?branch=master)](https://travis-ci.org/joaosoft/color) | [![codecov](https://codecov.io/gh/joaosoft/color/branch/master/graph/badge.svg)](https://codecov.io/gh/joaosoft/color) | [![Go Report Card](https://goreportcard.com/badge/github.com/joaosoft/color)](https://goreportcard.com/report/github.com/joaosoft/color) | [![GoDoc](https://godoc.org/github.com/joaosoft/color?status.svg)](https://godoc.org/github.com/joaosoft/color)

A color formatter that allows you to add color on your output.
The easy way to use the color:
``` Go
import log github.com/joaosoft/color

fmt.Fprintf(os.Stdout, fmt.Sprintf("%s joao", color.WithColor("hello", color.FormatBold, color.ForegroundRed, color.BackgroundCyan)))
```

###### If i miss something or you have something interesting, please be part of this project. Let me know! My contact is at the end.

## With support for
* text format
* foreground color
* background color

## Dependecy Management 
>### Dep

Project dependencies are managed using Dep. Read more about [Dep](https://github.com/golang/dep).
* Install dependencies: `dep ensure`
* Update dependencies: `dep ensure -update`

>### Go
```
go get github.com/joaosoft/color
```

## Usage 
This examples are available in the project at [color/examples](https://github.com/joaosoft/color/tree/master/examples)

```go
func main() {
	fmt.Fprintf(os.Stdout, fmt.Sprintf("%s joao", color.WithColor("hello", color.FormatBold, color.ForegroundRed, color.BackgroundCyan)))
}
```

## Known issues

## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
