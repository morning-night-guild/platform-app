package config_test

import (
	"reflect"
	"testing"

	"github.com/morning-night-guild/platform-app/internal/driver/config"
)

func TestNewAPI(t *testing.T) { //nolint:tparallel
	t.Parallel()

	tests := []struct {
		name string
		args config.APIConfig
		want config.APIConfig
	}{
		{
			name: "configを作成できる",
			args: config.APIConfig{
				Port:                "8080",
				APIKey:              "api_key",
				TokenSecret:         "",
				CookieDomain:        "",
				AppCoreURL:          "http://localhost:8080",
				NewRelicAppName:     "new_relic_app_name",
				NewRelicLicense:     "new_relic_license",
				CORSAllowOrigins:    "http://localhost:8080",
				CORSDebugEnable:     "true",
				RedisURL:            "",
				FirebaseSecret:      "",
				FirebaseAPIEndpoint: "",
				FirebaseAPIKey:      "",
			},
			want: config.APIConfig{
				Port:                "8080",
				APIKey:              "api_key",
				TokenSecret:         "",
				CookieDomain:        "",
				AppCoreURL:          "http://localhost:8080",
				NewRelicAppName:     "new_relic_app_name",
				NewRelicLicense:     "new_relic_license",
				CORSAllowOrigins:    "http://localhost:8080",
				CORSDebugEnable:     "true",
				RedisURL:            "",
				FirebaseSecret:      "",
				FirebaseAPIEndpoint: "",
				FirebaseAPIKey:      "",
			},
		},
		{
			name: "PORTの指定がなくてもconfigを作成できる",
			args: config.APIConfig{
				Port:                "",
				APIKey:              "api_key",
				TokenSecret:         "",
				CookieDomain:        "",
				AppCoreURL:          "http://localhost:8080",
				NewRelicAppName:     "new_relic_app_name",
				NewRelicLicense:     "new_relic_license",
				CORSAllowOrigins:    "http://localhost:8080",
				CORSDebugEnable:     "true",
				RedisURL:            "",
				FirebaseSecret:      "",
				FirebaseAPIEndpoint: "",
				FirebaseAPIKey:      "",
			},
			want: config.APIConfig{
				Port:                "8080",
				APIKey:              "api_key",
				TokenSecret:         "",
				CookieDomain:        "",
				AppCoreURL:          "http://localhost:8080",
				NewRelicAppName:     "new_relic_app_name",
				NewRelicLicense:     "new_relic_license",
				CORSAllowOrigins:    "http://localhost:8080",
				CORSDebugEnable:     "true",
				RedisURL:            "",
				FirebaseSecret:      "",
				FirebaseAPIEndpoint: "",
				FirebaseAPIKey:      "",
			},
		},
		{
			name: "PORTに数値に変換できない文字列が指定されてもconfigを作成できる",
			args: config.APIConfig{
				Port:                "8080",
				APIKey:              "api_key",
				TokenSecret:         "",
				CookieDomain:        "",
				AppCoreURL:          "http://localhost:8080",
				NewRelicAppName:     "new_relic_app_name",
				NewRelicLicense:     "new_relic_license",
				CORSAllowOrigins:    "http://localhost:8080",
				CORSDebugEnable:     "true",
				RedisURL:            "",
				FirebaseSecret:      "",
				FirebaseAPIEndpoint: "",
				FirebaseAPIKey:      "",
			},
			want: config.APIConfig{
				Port:                "8080",
				APIKey:              "api_key",
				TokenSecret:         "",
				CookieDomain:        "",
				AppCoreURL:          "http://localhost:8080",
				NewRelicAppName:     "new_relic_app_name",
				NewRelicLicense:     "new_relic_license",
				CORSAllowOrigins:    "http://localhost:8080",
				CORSDebugEnable:     "true",
				RedisURL:            "",
				FirebaseSecret:      "",
				FirebaseAPIEndpoint: "",
				FirebaseAPIKey:      "",
			},
		},
	}

	// 環境変数を操作するため直列でテスト
	for _, tt := range tests { //nolint:paralleltest
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("PORT", tt.args.Port)
			t.Setenv("API_KEY", tt.args.APIKey)
			t.Setenv("APP_CORE_URL", tt.args.AppCoreURL)
			t.Setenv("NEWRELIC_APP_NAME", tt.args.NewRelicAppName)
			t.Setenv("NEWRELIC_LICENSE", tt.args.NewRelicLicense)
			t.Setenv("CORS_ALLOW_ORIGINS", tt.args.CORSAllowOrigins)
			t.Setenv("CORS_DEBUG_ENABLE", tt.args.CORSDebugEnable)
			if got := config.NewAPI(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAPI() = %v, want %v", got, tt.want)
			}
		})
	}
}
