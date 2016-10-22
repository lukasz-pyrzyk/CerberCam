package main

import (
	"fmt"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("logger")

func main() {
	Receive()

	fmt.Printf("siema")
}
