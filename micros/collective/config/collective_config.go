package config

type (
	Configuration struct {
		BaseRoute      string
		QueryPrettyURL bool
		Debug          bool // Debug enables verbose logging of claims / cookies
	}
)

// CollectiveConfig holds the configuration values from collectives_config.yml file
var CollectiveConfig Configuration
