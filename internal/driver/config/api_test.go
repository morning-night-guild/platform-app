package config_test

import (
	"os"
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
				Port:             "8080",
				APIKey:           "api_key",
				AppCoreURL:       "http://localhost:8080",
				NewRelicAppName:  "new_relic_app_name",
				NewRelicLicense:  "new_relic_license",
				CORSAllowOrigins: "http://localhost:8080",
				CORSDebugEnable:  "true",
			},
			want: config.APIConfig{
				Port:             "8080",
				APIKey:           "api_key",
				AppCoreURL:       "http://localhost:8080",
				NewRelicAppName:  "new_relic_app_name",
				NewRelicLicense:  "new_relic_license",
				CORSAllowOrigins: "http://localhost:8080",
				CORSDebugEnable:  "true",
			},
		},
		{
			name: "PORTの指定がなくてもconfigを作成できる",
			args: config.APIConfig{
				Port:             "",
				APIKey:           "api_key",
				AppCoreURL:       "http://localhost:8080",
				NewRelicAppName:  "new_relic_app_name",
				NewRelicLicense:  "new_relic_license",
				CORSAllowOrigins: "http://localhost:8080",
				CORSDebugEnable:  "true",
			},
			want: config.APIConfig{
				Port:             "8080",
				APIKey:           "api_key",
				AppCoreURL:       "http://localhost:8080",
				NewRelicAppName:  "new_relic_app_name",
				NewRelicLicense:  "new_relic_license",
				CORSAllowOrigins: "http://localhost:8080",
				CORSDebugEnable:  "true",
			},
		},
		{
			name: "PORTに数値に変換できない文字列が指定されてもconfigを作成できる",
			args: config.APIConfig{
				Port:             "8080",
				APIKey:           "api_key",
				AppCoreURL:       "http://localhost:8080",
				NewRelicAppName:  "new_relic_app_name",
				NewRelicLicense:  "new_relic_license",
				CORSAllowOrigins: "http://localhost:8080",
				CORSDebugEnable:  "true",
			},
			want: config.APIConfig{
				Port:             "8080",
				APIKey:           "api_key",
				AppCoreURL:       "http://localhost:8080",
				NewRelicAppName:  "new_relic_app_name",
				NewRelicLicense:  "new_relic_license",
				CORSAllowOrigins: "http://localhost:8080",
				CORSDebugEnable:  "true",
			},
		},
	}

	// 環境変数を操作するため直列でテスト
	for _, tt := range tests { //nolint:paralleltest
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("PORT", tt.args.Port)
			os.Setenv("API_KEY", tt.args.APIKey)
			os.Setenv("APP_CORE_URL", tt.args.AppCoreURL)
			os.Setenv("NEWRELIC_APP_NAME", tt.args.NewRelicAppName)
			os.Setenv("NEWRELIC_LICENSE", tt.args.NewRelicLicense)
			os.Setenv("CORS_ALLOW_ORIGINS", tt.args.CORSAllowOrigins)
			os.Setenv("CORS_DEBUG_ENABLE", tt.args.CORSDebugEnable)
			if got := config.NewAPI(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAPI() = %v, want %v", got, tt.want)
			}
		})
	}
}
