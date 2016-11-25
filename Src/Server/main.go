package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"time"

	yaml "gopkg.in/yaml.v2"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("logger")
var GlobalConfig config

func main() {
	command := flag.String("command", "", "a command to run")
	configFile := flag.String("config", "config.yaml", "a configuration file to load")
	flag.Parse()

	loadConfiguration(*configFile)

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
}

func loadConfiguration(cfgFile string) {
	log.Infof("Loading file from %s", cfgFile)
	data, err := ioutil.ReadFile(cfgFile)
	failOnError(err, "Cannot load configuration file")

	yaml.Unmarshal(data, &GlobalConfig)
	log.Info("Loading configuration complete")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Criticalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func mainLoop(cmd commandType) {
	for {
		cmd()
		log.Debug("Thread sleep...")
		time.Sleep(1000)
	}
}
