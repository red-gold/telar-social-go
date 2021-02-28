package config

type (
	Configuration struct {
		BaseRoute      string
		QueryPrettyURL bool
		Debug          bool // Debug enables verbose logging of claims / cookies
	}
)

// CommentConfig holds the configuration values from comment-config.yml file
var CommentConfig Configuration
