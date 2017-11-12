package main

import (
	"flag"
	"fmt"
	"friday/command"
	"friday/config"
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
	var (
		configPath string
	)
	flag.StringVar(&configPath, "c", "friday.yaml", "Configuration file")

	preParseCommand()
	commandInfo := parseCommand()
	postParseCommand()

	initConfiguration(configPath)

	preCommandRun(commandInfo)
	err := commandInfo.Command.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
	postCommandRun(commandInfo)
}
