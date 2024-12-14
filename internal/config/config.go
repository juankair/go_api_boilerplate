package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const (
	defaultHostFrontend       = "http://localhost:3000/"
	defaultServerPort         = 8080
	defaultJWTExpirationHours = 72
)

type Config struct {
	ServerPort int `yaml:"server_port" env:"SERVER_PORT"`

	HostFrontend string `yaml:"host_frontend" env:"HOST_FRONTEND"`

	DSN string `yaml:"dsn" env:"DSN,secret"`

	JWTSigningKey string `yaml:"jwt_signing_key" env:"JWT_SIGNING_KEY,secret"`

	JWTExpiration int `yaml:"jwt_expiration" env:"JWT_EXPIRATION"`
}

func Load(file string) (*Config, error) {
	c := Config{
		ServerPort:    defaultServerPort,
		JWTExpiration: defaultJWTExpirationHours,
	}

	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(bytes, &c); err != nil {
		return nil, err
	}

	return &c, err
}
