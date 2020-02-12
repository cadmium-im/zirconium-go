package core

import (
	"github.com/google/logger"
)

var (
	moduleMgr     ModuleAPI
	router        *Router
	authManager   *AuthManager
	serverDomain  string
	pluginManager *PluginManager
)

func InitializeContext(sDomain string, pluginsDirPath string, enabledPlugins []string) {
	var err error
	moduleMgr, err = NewModuleManager()
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

	for _, v := range BuiltinPlugins {
		go v.Initialize(ModuleAPI(moduleMgr)) // Initialize builtin plugins
	}

	pluginManager = NewPluginManager()
	for _, v := range enabledPlugins {
		pluginManager.StartPlugin(pluginsDirPath, v, moduleMgr.(*ModuleManager))
	}
}
