package config

import "github.com/kelseyhightower/envconfig"

// Specification for basic configurations
type Specification struct {
	Port     string `envconfig:"PORT" default:"8080"`
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`
	Database Database
}

type Database struct {
	ReadDSN  string `envconfig:"DATABASE_READ_DSN" required:"true"`
	WriteDSN string `envconfig:"DATABASE_WRITE_DSN" required:"true"`
}

//LoadEnv loads environment variables
func LoadEnv() (*Specification, error) {
	var c Specification
	err := envconfig.Process("", &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
