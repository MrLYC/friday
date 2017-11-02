package vm

import (
	"flag"
	"io/ioutil"
)

// Command :
type Command struct {
	Path string
}

// GetDescription : get command description
func (c *Command) GetDescription() string {
	return "Lua shell"
}

// SetFlags : set parsing flags
func (c *Command) SetFlags() {
	flag.StringVar(&(c.Path), "path", "", "Path to lua file")
}

// Run : run command
func (c *Command) Run() error {
	data, err := ioutil.ReadFile(c.Path)
	if err != nil {
		return err
	}

	vm := &VM{}
	vm.Init()
	vm.Execute(string(data))

	return nil
}
