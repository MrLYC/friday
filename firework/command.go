package firework

import (
	"errors"
	"friday/logging"
	"os"
	"os/signal"
)

// Command :
type Command struct {
	Emitter *Emitter
}

// GetDescription : get command description
func (c *Command) GetDescription() string {
	return "Run firework command"
}

// SetFlags : set parsing flags
func (c *Command) SetFlags() {

}

// Run : run command
func (c *Command) Run() error {
	var err error

	c.Emitter = &Emitter{}
	c.Emitter.Init()

	c.Emitter.Ready()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	go func(ch chan os.Signal) {
		for s := range ch {
			if c.Emitter.Status == StatusControllerTerminating {
				logging.Infof("Killing emitter by: %v", s)
				c.Emitter.Kill()
				err = errors.New("Command killed")
			} else {
				logging.Infof("Terminating emitter by: %v", s)
				c.Emitter.Terminate()
			}
		}
	}(signalCh)

	c.Emitter.Run()

	return err
}
