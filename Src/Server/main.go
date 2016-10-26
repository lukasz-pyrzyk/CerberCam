package main

import (
	"flag"
	"fmt"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("logger")

func main() {
	command := (flag.String("command", "", "a command to run"))
	flag.Parse()

	switch *command {
	case "receive":
		fmt.Println("receive started")
		Receive()
		break
	case "send":
		fmt.Println("sending started")
		break
	default:
		fmt.Println("Invalid operation. Accepting: 'receive'")
	}
}
