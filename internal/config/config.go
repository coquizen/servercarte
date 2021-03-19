package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// Settings contains configuration settings for connecting to the db.

type Security struct {
	Length        int  `yaml:"length" default:"8"`
	MixedCase     bool `yaml:"mixed_case" default:"false"`
	AlphaNum      bool `yaml:"alpha_num" default:"false"`
	SpecialChar   bool `yaml:"special_char" default:"false"`
	CheckPrevious bool `yaml:"check_previous" default:"true"`
}

type Router struct {
	Host                string `yaml:"host" default:"127.0.0.1"`
	Port                string `yaml:"port,omitempty" default:"8080"`
	ReadTimeoutSeconds  int    `yaml:"read_timeout_seconds"`
	WriteTimeoutSeconds int    `yaml:"write_timeout_seconds"`
}

type Database struct {
	Type string `yaml:"type" required:"true"`
	Host string `yaml:"host" required:"true"`
	Port uint64 `yaml:"port,omitempty"`
	User string `required:"true"`
	Pass string `yaml:"pass,omitempty"`
	Name string `yaml:"name,omitempty"`
}

type Authentication struct {
	Algorithm        string `yaml:"algorithm"`
	ExpirationPeriod int    `yaml:"expiration_period"`
	MinKeyLength     int    `yaml:"minimum_key_length"`
	SecretKey        string `yaml:"secret_key"`
}

type config struct {
	Database       Database       `yaml:"database"`
	Server         Router         `yaml:"server"`
	Security       Security       `yaml:"security"`
	Authentication Authentication `yaml:"authentication"`
}

// LoadConfigFile loads the configuration from a local .yml into the struct
func Load(filePath string) (Router, Database, Security, Authentication, error) {
	var cfg config
	f, err := os.Open(filePath)
	if err != nil {
		return cfg.Server, cfg.Database, cfg.Security, cfg.Authentication, fmt.Errorf("error loading config.yml: %v",
			err)
	}

	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return cfg.Server, cfg.Database, cfg.Security, cfg.Authentication, err
	}

	return cfg.Server, cfg.Database, cfg.Security, cfg.Authentication, nil
}
