package config_test

import (
	"friday/config"
	"testing"
)

func TestConf(t *testing.T) {
	if config.Configuration.Debug != true {
		data, _ := config.Configuration.Dumps()
		t.Errorf("testing error:\n%s", data)
	}
}
