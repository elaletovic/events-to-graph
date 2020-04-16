package config

import "os"

// Config --
type Config struct {
	DBAddress string
}

// Load returns Config obj
func Load() Config {
	return Config{
		DBAddress: env("DB_ADDRESS", "http://localhost:8529"),
	}
}

func env(key, def string) string {
	val := os.Getenv(key)

	if len(val) == 0 {
		val = def
	}
	return val
}
