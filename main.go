package main

import (
	"flag"
	"fmt"
	"friday/command"
	"friday/config"
	"friday/firework"
	"friday/logging"
	"friday/storage/migration"
	"friday/vm"
	"math/rand"
	"os"
	"time"
)

func parseCommand() *command.CommandInfo {
	factory := command.Factory{
		Commands: map[string]command.ICommand{
			"version":  &config.VersionCommand{},
			"confinfo": &config.ConfigurationCommand{},
			"vm":       &vm.Command{},
			"migrate":  &migration.Command{},
			"run":      &firework.Command{},
		},
	}
	factory.Init()
	return factory.ParseCommand()
}

func initConfiguration() {
	err := config.Configuration.Read()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	err = config.Configuration.Validate()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(-1)
	}
}

func preParseCommand() {
	rand.Seed(time.Now().Unix())
}

func postParseCommand() {
}

func preCommandRun(commandInfo *command.CommandInfo) {
	logging.Init()
}

func postCommandRun(commandInfo *command.CommandInfo) {
}

func main() {
	flag.StringVar(
		&(config.Configuration.ConfigurationPath),
		"c", config.Configuration.ConfigurationPath,
		"Configuration file",
	)

	preParseCommand()
	commandInfo := parseCommand()
	postParseCommand()

	initConfiguration()

	preCommandRun(commandInfo)
	err := commandInfo.Command.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
	postCommandRun(commandInfo)
}
