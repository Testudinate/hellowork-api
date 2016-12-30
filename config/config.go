package config

import "github.com/kelseyhightower/envconfig"

// Specification for basic configurations
type Specification struct {
	Port        string `envconfig:"PORT" default:"8080"`
	LogLevel    string `envconfig:"LOG_LEVEL" default:"info"`
	DatabaseDSN string `envconfig:"DATABASE_DSN" required:"true"`
}

//LoadEnv loads environment variables
func LoadEnv() (*Specification, error) {
	var config Specification
	err := envconfig.Process("", &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
