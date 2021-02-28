package config

type (
	Configuration struct {
		BaseRoute      string
		QueryPrettyURL bool
		Debug          bool // Debug enables verbose logging of claims / cookies
	}
)

// MediaConfig holds the configuration values from media-config.yml file
var MediaConfig Configuration
