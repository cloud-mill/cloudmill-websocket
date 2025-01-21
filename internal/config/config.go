package config

import (
	"flag"
	"strings"

	"github.com/cloud-mill/cloudmill-websocket/internal/logger"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var Config CloudmillWebsocketConfig

var InstanceId string

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
	if value := viper.GetString(key); value != "" {
		return strings.Split(value, ",")
	}
	return nil
}

func init() {
	profile := flag.String(
		"profile",
		"local",
		"Environment profile, something similar to spring profiles",
	)
	flag.Parse()

	logger.Logger.Info("profile at", zap.String("profile", *profile))

	viper.SetDefault("profile", *profile)
	viper.AutomaticEnv()

	viper.SetEnvPrefix("cm")

	bindEnv("host")
	bindEnv("port")
	bindEnv("allowed_origins")
	bindEnv("auth_middleware_secret_key")
	bindEnv("jwt_cookie_name")
	bindEnv("csrf_cookie_name")
	bindEnv("csrf_header_name")

	Config = CloudmillWebsocketConfig{
		Host:           viper.GetString("host"),
		Port:           viper.GetInt("port"),
		AllowedOrigins: getEnvAsStringSlice("allowed_origins"),
		Auth: AuthConfig{
			AuthMiddlewareSecretKey: viper.GetString("auth_middleware_secret_key"),
			JwtCookieName:           viper.GetString("jwt_cookie_name"),
			CsrfCookieName:          viper.GetString("csrf_cookie_name"),
			CsrfHeaderName:          viper.GetString("csrf_header_name"),
		},
	}

	logger.Logger.Info("got config", zap.Any("CloudmillWebsocketConfig", Config))

	InstanceId = uuid.New().String()
	logger.Logger.Info("instance id", zap.String("InstanceId", InstanceId))
}
