package main

import (
	"flag"
	"fmt"
	"friday/command"
	"friday/config"
	"friday/logging"
	"friday/sentry"
	"math/rand"
	"os"
	"time"
)

func parseCommand() *command.CommandInfo {
	factory := command.Factory{
		Commands: map[string]command.ICommand{
			"version":  &config.VersionCommand{},
			"confinfo": &config.ConfigurationCommand{},
		},
	}
	factory.Init()
	return factory.ParseCommand()
}

func initConfiguration(configPath string) {
	config.Configuration.Init()
	err := config.Configuration.ReadFrom(configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

func preParseCommand() {
	rand.Seed(time.Now().Unix())
}

func postParseCommand() {
}

func preCommandRun(commandInfo *command.CommandInfo) {
	logging.Init()
	sentry.Init()
}

func postCommandRun(commandInfo *command.CommandInfo) {
}

func main() {
	var (
		configPath string
	)
	flag.StringVar(&configPath, "c", "friday.yaml", "Configuration file")

	preParseCommand()
	commandInfo := parseCommand()
	postParseCommand()

	initConfiguration(configPath)

	preCommandRun(commandInfo)
	commandInfo.Command.Run()
	postCommandRun(commandInfo)
}
