package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cadmium-im/zirconium-go/core"
	"github.com/google/logger"
	"github.com/pelletier/go-toml"
)

func main() {
	var cfg core.Config
	var configPath string
	var generateConfig bool

	defer logger.Init("auth-server", true, false, ioutil.Discard).Close() // TODO Make ability to use file for log output
	flag.StringVar(&configPath, "config", "", "Path to config")
	flag.BoolVar(&generateConfig, "gen-config", false, "Generate the config")
	flag.Parse()
	if generateConfig == true {
		sampleConfig := &core.Config{}
		val, err := toml.Marshal(sampleConfig)
		if err != nil {
			logger.Errorf("Failed to generate config: %s", err.Error())
			os.Exit(1)
		}
		fmt.Println(string(val))
		os.Exit(0)
	}
	if configPath == "" {
		logger.Error("Path to config isn't specified!")
		os.Exit(1)
	}
	cfgData, err := ioutil.ReadFile(configPath)
	if err != nil {
		logger.Error("Failed to read config!")
		os.Exit(1)
	}
	err = toml.Unmarshal(cfgData, &cfg)
	if err != nil {
		logger.Errorf("Failed to read config! (yaml error: %s)", err.Error())
		os.Exit(1)
	}
	err = validateConfig(&cfg)
	if err != nil {
		logger.Errorf("Config validation failed: %s", err.Error())
		os.Exit(1)
	}

	ac := core.NewAppContext(&cfg)
	logger.Info("Zirconium successfully started!")
	logger.Fatal(ac.Run())
}

func validateConfig(config *core.Config) error {
	if config.ServerID == "" {
		return errors.New("server id isn't specified")
	}
	return nil
}
