package config_test

import (
	"friday/config"
	"testing"
)

func TestConf(t *testing.T) {
	err := config.Configuration.Read()
	if err != nil {
		t.Errorf("read error: %v", err)
	}
	if config.Configuration.Debug != true {
		data, _ := config.Configuration.Dumps()
		t.Errorf("testing error:\n%s", data)
	}
}
