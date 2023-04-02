package config

import (
	"os"
	"strconv"

	"github.com/morning-night-guild/platform-app/pkg/log"
)

type APIConfig struct {
	Port                string
	APIKey              string
	JWTSecret           string
	CookieDomain        string
	AppCoreURL          string
	NewRelicAppName     string
	NewRelicLicense     string
	CORSAllowOrigins    string
	CORSDebugEnable     string
	RedisURL            string
	FirebaseSecret      string
	FirebaseAPIEndpoint string
	FirebaseAPIKey      string
}

func NewAPI() APIConfig {
	port := os.Getenv("PORT")

	if _, err := strconv.Atoi(port); err != nil {
		port = "8080"
	}

	conf := APIConfig{
		Port:                port,
		APIKey:              os.Getenv("API_KEY"),
		JWTSecret:           os.Getenv("JWT_SECRET"),
		CookieDomain:        os.Getenv("COOKIE_DOMAIN"),
		AppCoreURL:          os.Getenv("APP_CORE_URL"),
		NewRelicAppName:     os.Getenv("NEWRELIC_APP_NAME"),
		NewRelicLicense:     os.Getenv("NEWRELIC_LICENSE"),
		CORSAllowOrigins:    os.Getenv("CORS_ALLOW_ORIGINS"),
		CORSDebugEnable:     os.Getenv("CORS_DEBUG_ENABLE"),
		RedisURL:            os.Getenv("REDIS_URL"),
		FirebaseSecret:      os.Getenv("FIREBASE_SECRET"),
		FirebaseAPIEndpoint: os.Getenv("FIREBASE_API_ENDPOINT"),
		FirebaseAPIKey:      os.Getenv("FIREBASE_API_KEY"),
	}

	log.Log().Sugar().Infof("config: %+v", conf)

	return conf
}
