package config

type (
	Configuration struct {
		BaseRoute      string
		QueryPrettyURL bool
		Debug          bool // Debug enables verbose logging of claims / cookies
	}
)

// VoteConfig holds the configuration values from vote-config.yml file
var VoteConfig Configuration
