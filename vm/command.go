package vm

import (
	"flag"
	"io/ioutil"

	"friday/command"
)

// Command :
type Command struct {
	Path command.TFlagStringArr
}

// GetDescription : get command description
func (c *Command) GetDescription() string {
	return "Lua shell"
}

// SetFlags : set parsing flags
func (c *Command) SetFlags() {
	flag.Var(&c.Path, "path", "Path to lua file")
}

// Run : run command
func (c *Command) Run() error {
	vm := &VM{}
	vm.Init()

	for _, i := range c.Path {
		path := string(i)
		if path == "" {
			continue
		}
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		err = vm.Execute(string(data))
		if err != nil {
			return err
		}
	}

	return nil
}
