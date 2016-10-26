package main

import (
	"fmt"
	"os"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("logger")

func main() {
	args := os.Args[1:] // get arguments without prog

	switch args[0] {
	case "receive":
		fmt.Println("receive started")
		Receive()
		break
	case "send":
		fmt.Println("sending started")
		break
	}
}
