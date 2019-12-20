package config

type (
	Configuration struct {
		BaseRoute      string
		QueryPrettyURL bool
		Debug          bool // Debug enables verbose logging of claims / cookies
	}
)

// CircleConfig holds the configuration values from circle-config.yml file
var CircleConfig Configuration
