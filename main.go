package main

import (
	"flag"
	"fmt"
	"friday/command"
	"friday/config"
	"math/rand"
	"os"
	"time"
)

func main() {
	var (
		configPath string
		err        error
	)
	factory := command.Factory{
		Commands: map[string]command.ICommand{
			"version":  &config.VersionCommand{},
			"confinfo": &config.ConfigurationCommand{},
		},
	}
	factory.Init()

	rand.Seed(time.Now().Unix())

	flag.StringVar(&configPath, "c", "friday.yaml", "Configuration file")
	command := factory.ParseCommand()

	config.Configuration.Init()
	err = config.Configuration.ReadFrom(configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	command.Run()
}
