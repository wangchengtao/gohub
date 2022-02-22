package config

import "gohub/pkg/config"

func init() {
	config.Add("app", func() map[string]interface{} {
		return map[string]interface{}{
			"name":     config.Env("APP_NAME", "Gohub"),
			"env":      config.Env("APP_ENV", "production"),
			"debug":    config.Env("APP_DEBUG", false),
			"port":     config.Env("APP_PORT", "3000"),
			"key":      config.Env("APP_KEY", "33446a9dcf9ea060a0a6532b166da32f304af0de"),
			"url":      config.Env("APP_URL", "http://localhost:3000"),
			"timezone": config.Env("TIMEZONE", "Asia/Shanghai"),
		}
	})
}
