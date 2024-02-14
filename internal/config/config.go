package config

import (
	"github.com/juankair/go_api_boilerplate/pkg/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const (
	defaultServerPort         = 8080
	defaultJWTExpirationHours = 72
)

type Config struct {
	ServerPort int `yaml:"server_port" env:"SERVER_PORT"`

	DSN string `yaml:"dsn" env:"DSN,secret"`

	JWTSigningKey string `yaml:"jwt_signing_key" env:"JWT_SIGNING_KEY,secret"`

	JWTExpiration int `yaml:"jwt_expiration" env:"JWT_EXPIRATION"`
}

func Load(file string, logger log.Logger) (*Config, error) {
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

	//if err = env.New("APP_", logger.Infof).Load(&c); err != nil {
	//	return nil, err
	//}

	return &c, err
}
