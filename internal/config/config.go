package config

import (
	"flag"
	"strings"

	"github.com/cloud-mill/cloudmill-websocket/internal/logger"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	Config     CloudmillWebsocketConfig
	InstanceId string
)

type AuthConfig struct {
	AuthMiddlewareSecretKey string `json:"auth_middleware_secret_key"`
	JwtCookieName           string `json:"jwt_cookie_name"`
	CsrfCookieName          string `json:"csrf_cookie_name"`
	CsrfHeaderName          string `json:"csrf_header_name"`
}

type CloudmillWebsocketConfig struct {
	Host           string     `json:"host"`
	Port           int        `json:"port"`
	AllowedOrigins []string   `json:"allowed_origins"`
	Auth           AuthConfig `json:"auth"`
}

const (
	EnvHost                    = "host"
	EnvPort                    = "port"
	EnvAllowedOrigins          = "allowed_origins"
	EnvAuthMiddlewareSecretKey = "auth_middleware_secret_key"
	EnvJwtCookieName           = "jwt_cookie_name"
	EnvCsrfCookieName          = "csrf_cookie_name"
	EnvCsrfHeaderName          = "csrf_header_name"
)

func bindEnv(key string) {
	if err := viper.BindEnv(key); err != nil {
		logger.Logger.Fatal(
			"error binding environment variable",
			zap.String("key", key),
			zap.Error(err),
		)
	}
}

func getEnvAsStringSlice(key string) []string {
	value := viper.GetString(key)
	if value != "" {
		return strings.Split(value, ",")
	}
	return nil
}

func validateConfig() {
	requiredKeys := []string{
		EnvHost,
		EnvPort,
		EnvAuthMiddlewareSecretKey,
		EnvJwtCookieName,
		EnvCsrfCookieName,
		EnvCsrfHeaderName,
	}

	for _, key := range requiredKeys {
		if !viper.IsSet(key) {
			logger.Logger.Warn("missing required environment variable", zap.String("key", key))
		}
	}
}

func init() {
	profile := flag.String(
		"profile",
		"local",
		"environment profile (e.g., local, dev, prod)",
	)
	flag.Parse()

	logger.Logger.Info("environment profile", zap.String("profile", *profile))

	viper.SetDefault("profile", *profile)
	viper.SetEnvPrefix("cm")
	viper.AutomaticEnv()

	bindEnv(EnvHost)
	bindEnv(EnvPort)
	bindEnv(EnvAllowedOrigins)
	bindEnv(EnvAuthMiddlewareSecretKey)
	bindEnv(EnvJwtCookieName)
	bindEnv(EnvCsrfCookieName)
	bindEnv(EnvCsrfHeaderName)

	validateConfig()

	Config = CloudmillWebsocketConfig{
		Host:           viper.GetString(EnvHost),
		Port:           viper.GetInt(EnvPort),
		AllowedOrigins: getEnvAsStringSlice(EnvAllowedOrigins),
		Auth: AuthConfig{
			AuthMiddlewareSecretKey: viper.GetString(EnvAuthMiddlewareSecretKey),
			JwtCookieName:           viper.GetString(EnvJwtCookieName),
			CsrfCookieName:          viper.GetString(EnvCsrfCookieName),
			CsrfHeaderName:          viper.GetString(EnvCsrfHeaderName),
		},
	}

	logger.Logger.Info("loaded config", zap.Any("Config", Config))

	InstanceId = uuid.New().String()
	logger.Logger.Info("initialised instance", zap.String("instance_id", InstanceId))
}
