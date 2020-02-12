package core

import (
	"net/rpc"

	"github.com/google/logger"
	"github.com/hashicorp/go-plugin"
)

type moduleClient struct {
	Broker *plugin.MuxBroker
	Client *rpc.Client
}

// Client implmentation of go-plugin.plugin.Plugin.Client
func (p *ModuleRef) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &moduleClient{Broker: b, Client: c}, nil
}

// Name initiates an RPC call to the plugin name
func (p *moduleClient) Name() string {
	var resp string
	err := p.Client.Call("Plugin.Name", new(interface{}), &resp)
	if err != nil {
		logger.Fatal(err) // FIXME
	}
	return resp
}

// Version initiates an RPC call to the plugin version
func (p *moduleClient) Version() string {
	var resp string
	err := p.Client.Call("Plugin.Version", new(interface{}), &resp)
	if err != nil {
		logger.Fatal(err) // FIXME
	}
	return resp
}

// StartTime initiates an RPC call to the plugin for initializing
func (p *moduleClient) Initialize(moduleAPI ModuleAPI) {
	var resp interface{}
	err := p.Client.Call("Plugin.Initialize", map[string]interface{}{
		"moduleAPI": moduleAPI,
	}, &resp)
	if err != nil {
		logger.Fatal(err) // FIXME
	}
}
