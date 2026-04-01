package config

import "os"

type Config struct {
	AppPort string
	DBURL   string
}

func Load() Config {
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "3000"
	}

	return Config{
		AppPort: appPort,
		DBURL:   os.Getenv("DATABASE_URL"),
	}
}
