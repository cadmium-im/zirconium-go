package core

type Config struct {
	ServerDomains []string `comment:"Server domain names (e.g. example.com)"`
	ServerID      string
	Websocket     struct {
		Host     string
		Port     int
		Endpoint string
	}
	Mongo struct {
		Host     string
		Port     int
		User     string
		Password string
		Database string
	}
}
