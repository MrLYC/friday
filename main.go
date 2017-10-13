package main

import (
	"friday/command"
)

func main() {
	factory := command.Factory{
		HelpFlag: "h",
		Commands: map[string]command.ICommand{},
	}
	factory.Init()
	command := factory.ParseCommand()
	command.Run()
}
