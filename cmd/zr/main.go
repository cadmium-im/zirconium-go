package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	core "github.com/ChronosX88/zirconium/core"
	"github.com/google/logger"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/pelletier/go-toml"
)

var connectionHandler = core.NewConnectionHandler()
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {

	var cfg core.ServerConfig
	var configPath string
	var generateConfig bool
	flag.StringVar(&configPath, "config", "", "Path to config")
	flag.BoolVar(&generateConfig, "gen_config", false, "Generate the config")
	flag.Parse()
	if generateConfig == true {
		sampleConfig := &core.ServerConfig{}
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

	core.InitializeContext(cfg.ServerDomain, cfg.PluginsDirPath, cfg.EnabledPlugins)
	router := mux.NewRouter()
	router.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		response.Write([]byte("Zirconium server is up and running!"))
	}).Methods("GET")
	router.HandleFunc("/ws", wsHandler)

	logger.Info("Zirconium successfully started!")
	logger.Fatal(http.ListenAndServe(":8844", router))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	connectionHandler.HandleNewConnection(ws)
}

func validateConfig(config *core.ServerConfig) error {
	if config.ServerDomain == "" {
		return errors.New("server domain isn't specified")
	} else if config.PluginsDirPath == "" {
		return errors.New("plugin directory path isn't specified")
	}
	return nil
}
