package core

type ServerConfig struct {
	// A list of enabled plugins (or extensions) in server
	EnabledPlugins []string `toml:"enabledPlugins" comment:"A list of enabled plugins (or extensions) in server"`

	// Server domain name (e.g. example.com)
	ServerDomain string `toml:"serverDomain" comment:"Server domain name (e.g. example.com)"`

	// Path to directory with plugin executables
	PluginsDirPath string `toml:"pluginsDirPath" comment:"Path to directory with plugin executables"`
}
