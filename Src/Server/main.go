package main

import (
	"flag"
	"time"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("logger")
var modeldir *string

// CommandType - signature of command functions
type CommandType func()

func main() {
	modeldir = flag.String("dir", "", "Directory containing the trained model files. The directory will be created and the model downloaded into it if necessary")
	command := flag.String("command", "", "a command to run")
	flag.Parse()

	switch *command {
	case "receive":
		log.Info("Receive started!")
		mainLoop(HandleReceiveCommand)
		break
	case "send":
		log.Info("Sending started!")
		mainLoop(HandleSendCommand)
		//Send("alertsQueue")
		break
	default:
		log.Errorf("Invalid operation. Accepting: 'receive' or 'send', %s provided", *command)
		flag.Usage()
	}

	if *modeldir == "" {
		log.Error("Invalid operation, model directory is not provided!")
		flag.Usage()
	}
}

func mainLoop(cmd CommandType) {
	for {
		cmd()
		log.Debug("Thread sleep...")
		time.Sleep(1000)
	}
}
