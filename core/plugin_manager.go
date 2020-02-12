package core

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/logger"
	"github.com/hashicorp/go-plugin"
)

type PluginManager struct{}

func NewPluginManager() *PluginManager {
	return &PluginManager{}
}

func (p *PluginManager) StartPlugin(pluginsDirPath, pluginFile string, moduleManager *ModuleManager) error {
	pluginsDirectory, _ := filepath.Abs(filepath.Dir(pluginsDirPath))
	pluginFile = filepath.Join(pluginsDirectory, pluginFile)

	logger.Info("Starting plugin: %s", pluginFile)

	client := plugin.NewClient(&plugin.ClientConfig{
		Cmd:        exec.Command(pluginFile),
		Managed:    true,
		SyncStdout: os.Stdout,
		SyncStderr: os.Stderr,

		HandshakeConfig: HandshakeConfig,
		Plugins:         GetPluginMap(nil),
	})

	rpcclient, err := client.Client()

	if err != nil {
		logger.Errorf("Failed to get RPC Client: %s", err)
		client.Kill()
		return err
	}

	// get the interface
	raw, err := rpcclient.Dispense(ModuleInterfaceName)
	if err != nil {
		logger.Errorf("Failed to get interface: %s error: %s", ModuleInterfaceName, err)
		return err
	}

	moduleObj := raw.(Module)

	go func() {
		moduleObj.Initialize(moduleManager)
	}()

	return nil
}
