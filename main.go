package main

import (
	"friday/command"
	"friday/config"
)

func main() {
	factory := command.Factory{
		HelpFlag: "h",
		Commands: map[string]command.ICommand{
			"version": &config.VersionCommand{},
		},
	}
	factory.Init()
	command := factory.ParseCommand()
	command.Run()
}
