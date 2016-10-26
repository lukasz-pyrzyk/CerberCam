package main

import (
	"flag"
	"time"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("logger")

// CommandType - signature of command functions
type CommandType func()

func main() {
	command := flag.String("command", "", "a command to run")
	flag.Parse()

	switch *command {
	case "receive":
		log.Info("Receive started!")
		mainLoop(HandleReceiveCommand)
		break
	case "send":
		log.Info("Sending started!")
		Send("alertsQueue")
		break
	default:
		log.Errorf("Invalid operation. Accepting: 'receive' or 'send', %s provided", *command)
	}
}

func mainLoop(cmd CommandType) {
	for {
		cmd()
		log.Debug("Thread sleep...")
		time.Sleep(1000)
	}
}
