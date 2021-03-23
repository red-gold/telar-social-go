package config

type (
	Configuration struct {
		BaseRoute      string
		QueryPrettyURL bool
		Debug          bool // Debug enables verbose logging of claims / cookies
	}
)

// VangConfig holds the configuration values from vang-config.yml file
var VangConfig Configuration
