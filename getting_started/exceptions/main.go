package exceptions

import (
	"fmt"
	. "github.com/joaosoft/go-manager/exceptions"
)

func main() {
	fmt.Println("We started")
	Block{
		Try: func() {
			fmt.Println("I tried")
			Throw("Oh,...sh...")
		},
		Catch: func(e Exception) {
			fmt.Printf("Caught %v\n", e)
		},
		Finally: func() {
			fmt.Println("Finally...")
		},
	}.Do()
	fmt.Println("We went on")
}
