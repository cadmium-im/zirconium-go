package internal

import (
	"github.com/ChronosX88/zirconium/shared"
	"github.com/google/logger"
)

var (
	ModuleMgr    *ModuleManager
	router       *Router
	authManager  *AuthManager
	serverDomain string
)

func InitializeContext(sDomain string) {
	var err error
	ModuleMgr, err = NewModuleManager()
	if err != nil {
		logger.Fatalf("Unable to initialize module manager: %s", err.Error())
	}

	router, err = NewRouter()
	if err != nil {
		logger.Fatalf("Unable to initialize router: %s", err.Error())
	}

	authManager, err = NewAuthManager()
	if err != nil {
		logger.Fatalf("Unable to initialize authentication manager: %s", err.Error())
	}
	serverDomain = sDomain

	for _, v := range shared.Plugins {
		go v.Initialize() // Initialize provided plugins
	}
}
