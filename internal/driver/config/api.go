package config

import (
	"os"
	"strconv"

	"github.com/morning-night-guild/platform-app/pkg/log"
)

type APIConfig struct {
	Port             string
	APIKey           string
	AppCoreURL       string
	NewRelicAppName  string
	NewRelicLicense  string
	CORSAllowOrigins string
	CORSDebugEnable  string
}

func NewAPI() APIConfig {
	port := os.Getenv("PORT")

	if _, err := strconv.Atoi(port); err != nil {
		port = "8080"
	}

	c := APIConfig{
		Port:             port,
		APIKey:           os.Getenv("API_KEY"),
		AppCoreURL:       os.Getenv("APP_CORE_URL"),
		NewRelicAppName:  os.Getenv("NEWRELIC_APP_NAME"),
		NewRelicLicense:  os.Getenv("NEWRELIC_LICENSE"),
		CORSAllowOrigins: os.Getenv("CORS_ALLOW_ORIGINS"),
		CORSDebugEnable:  os.Getenv("CORS_DEBUG_ENABLE"),
	}

	log.Log().Sugar().Infof("config: %+v", c)

	return c
}
