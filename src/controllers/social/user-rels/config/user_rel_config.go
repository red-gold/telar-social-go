package config

type (
	Configuration struct {
		BaseRoute      string
		QueryPrettyURL bool
		Debug          bool // Debug enables verbose logging of claims / cookies
	}
)

// UserRelConfig holds the configuration values from user-rel-config.yml file
var UserRelConfig Configuration
