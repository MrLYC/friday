package command

import (
	"flag"
	"fmt"
	"os"
)

// Factory : factory to get command
type Factory struct {
	Name        string
	Command     ICommand
	CommandName string
	HelpFlag    string
	Help        bool
	Commands    map[string]ICommand
}

// SetFlags : set parsing flags
func (f *Factory) SetFlags() {
}

// GetDescription : get command description
func (f *Factory) GetDescription() string {
	return "Print usage"
}

// Run : run command
func (f *Factory) Run() error {
	if f.Command == nil || f.Command == f {
		fmt.Printf("Supported commands:\n %v", f.Name)
		for k := range f.Commands {
			fmt.Printf(" %v", k)
		}
		fmt.Printf("\n")
	} else {
		fmt.Printf("Command %v: %v\n", f.CommandName, f.Command.GetDescription())
	}
	flag.Usage()
	os.Exit(1)
	return nil
}

// Init : init factory
func (f *Factory) Init() {
	f.Name = "usage"
	flag.BoolVar(&f.Help, f.HelpFlag, false, f.GetDescription())
}

// ParseCommand : parse from command line arguments and get command
func (f *Factory) ParseCommand() ICommand {
	var (
		argv        = len(os.Args)
		commandArgs []string
		command     ICommand
		ok          bool
	)

	if argv > 1 {
		f.CommandName = os.Args[1]
	}
	if argv > 2 {
		commandArgs = os.Args[2:]
	}

	command, ok = f.Commands[f.CommandName]
	if ok {
		f.Command = command
	} else {
		f.CommandName = f.Name
		command = f
	}

	command.SetFlags()
	flag.CommandLine.Parse(commandArgs)

	if f.Help {
		f.Run()
	}
	return command
}
