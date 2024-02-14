package env

import "os"

type Config struct {
	Host string
	Port int
}

func Load() {
	os.Setenv("APP_HOST", "127.0.0.1")
	os.Setenv("APP_PORT", "8080")
}
