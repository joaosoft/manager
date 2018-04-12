/*
GoManager helps you to manage you application

With support for
* Processes
* Configurations (with reload and write options)
* NSQ Consumers
* NSQ Producers
* Database Connections
* Web Servers
* Gateways
* Redis Connections
* Work Queues (with FIFO and LIFO modes)

Usage
at https://github.com/joaosoft/go-manager/tree/master/bin/launcher
*/
package gomanager

// Please always use the parenthetical form of import
// even though when only importing one package
// import "fmt"
// will work, it makes it uniform and simpler to add later

import (
	"math"
)

// You can godoc constants

// Some enum examples
const (
	CONSTA = iota
	CONSTB
	CONSTC
	CONSTD
	ANOTHER = 7
)

// You can godoc vars

// This is just a random variable
var Default float64 = 0.7

// var Default = float64(0.7) would've worked as well

// You can godoc types

// Example is a float used for demonstration purposes
type Example float64

// Example2 is also for demonstartion
type Example2 struct {
	X Example
	y int // Private fields do not appear in godoc
}

// You can godoc functions

// NewExample is used to get a ready-to-use Example2
func NewExample(num int) *Example2 {
	return &Example{0.0, num}
}

// You can godoc methods

// Returns the square root of an example
func (e Example) Sqrt() Example {
	return Example(math.Sqrt(float64(e)))
}