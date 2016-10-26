package main

import (
	"flag"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("logger")

func main() {
	command := flag.String("command", "", "a command to run")
	flag.Parse()

	switch *command {
	case "receive":
		log.Info("Receive started")
		Receive("picturesQueue")
		break
	case "send":
		log.Info("sending started")
		Send("alertsQueue")
		break
	default:
		log.Error("Invalid operation. Accepting: 'receive'")
	}
}
