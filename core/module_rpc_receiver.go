package core

import (
	"errors"

	"github.com/hashicorp/go-plugin"
)

type greeterServer struct {
	Broker *plugin.MuxBroker
	Module Module
}

// Server implmentation of go-plugin.plugin.Plugin.Server
func (p *ModuleRef) Server(b *plugin.MuxBroker) (interface{}, error) {
	if p.F == nil {
		return nil, errors.New("Greeter interface not implemeted")
	}
	return &greeterServer{Broker: b, Module: p.F()}, nil
}

// Name calls the plugin implementation to get the name of the plugin
func (p *greeterServer) Name(nothing interface{}, result *string) error {
	*result = p.Module.Name()
	return nil
}

// Version calls the plugin implementation to get the version of the plugin
func (p *greeterServer) Version(nothing interface{}, result *string) error {
	*result = p.Module.Version()
	return nil
}

// StartTime calls the plugin implementation to initialize plugin
func (p *greeterServer) Initialize(moduleAPI ModuleAPI, emptyResult interface{}) error {
	p.Module.Initialize(moduleAPI)
	return nil
}
