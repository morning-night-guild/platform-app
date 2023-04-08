package config_test

import (
	"reflect"
	"testing"

	"github.com/morning-night-guild/platform-app/internal/driver/config"
)

func TestNewCore(t *testing.T) { //nolint:tparallel
	tests := []struct {
		name string
		args config.CoreConfig
		want config.CoreConfig
	}{
		{
			name: "configを作成できる",
			args: config.CoreConfig{
				Port:            "8080",
				DSN:             "dsn",
				APIKey:          "api_key",
				NewRelicAppName: "new_relic_app_name",
				NewRelicLicense: "new_relic_license",
			},
			want: config.CoreConfig{
				Port:            "8080",
				DSN:             "dsn",
				APIKey:          "api_key",
				NewRelicAppName: "new_relic_app_name",
				NewRelicLicense: "new_relic_license",
			},
		},
		{
			name: "PORTの指定がなくてもconfigを作成できる",
			args: config.CoreConfig{
				Port:            "",
				DSN:             "dsn",
				APIKey:          "api_key",
				NewRelicAppName: "new_relic_app_name",
				NewRelicLicense: "new_relic_license",
			},
			want: config.CoreConfig{
				Port:            "8080",
				DSN:             "dsn",
				APIKey:          "api_key",
				NewRelicAppName: "new_relic_app_name",
				NewRelicLicense: "new_relic_license",
			},
		},
		{
			name: "PORTに数値に変換できない文字列が指定されてもconfigを作成できる",
			args: config.CoreConfig{
				Port:            "port",
				DSN:             "dsn",
				APIKey:          "api_key",
				NewRelicAppName: "new_relic_app_name",
				NewRelicLicense: "new_relic_license",
			},
			want: config.CoreConfig{
				Port:            "8080",
				DSN:             "dsn",
				APIKey:          "api_key",
				NewRelicAppName: "new_relic_app_name",
				NewRelicLicense: "new_relic_license",
			},
		},
	}

	// 環境変数を操作するため直列でテスト
	for _, tt := range tests { //nolint:paralleltest
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("PORT", tt.args.Port)
			t.Setenv("DATABASE_URL", tt.args.DSN)
			t.Setenv("API_KEY", tt.args.APIKey)
			t.Setenv("NEWRELIC_APP_NAME", tt.args.NewRelicAppName)
			t.Setenv("NEWRELIC_LICENSE", tt.args.NewRelicLicense)
			if got := config.NewCore(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCore() = %v, want %v", got, tt.want)
			}
		})
	}
}
