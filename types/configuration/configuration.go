package configuration

type Config struct {
	ListenAddress string
}

var currentConfig Config

// TODO: Validate and return error on invalid configurations
func SetConfig(config *Config) {
	currentConfig = *config
}

func GetConfig() *Config {
	return &currentConfig
}
