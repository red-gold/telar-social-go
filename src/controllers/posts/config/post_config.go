package config

type (
	Configuration struct {
		BaseRoute      string
		QueryPrettyURL bool
		Debug          bool // Debug enables verbose logging of claims / cookies
	}
)

// PostConfig holds the configuration values from post-config.yml file
var PostConfig Configuration
