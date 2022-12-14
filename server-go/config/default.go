package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBUri                  string        `mapstructure:"MONGODB_URI"`
	RedisUri               string        `mapstructure:"REDIS_URL"`
	Port                   string        `mapstructure:"PORT"`

	ClientOrigin string `mapstructure:"CLIENT_ORIGIN"`

	AccessTokenPrivateKey  string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	AccessTokenPublicKey   string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`
	RefreshTokenPrivateKey string        `mapstructure:"REFRESH_TOKEN_PRIVATE_KEY"`
	RefreshTokenPublicKey  string        `mapstructure:"REFRESH_TOKEN_PUBLIC_KEY"`
	AccessTokenExpiresIn   time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`
	RefreshTokenExpiresIn  time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_IN"`
	AccessTokenMaxAge      int           `mapstructure:"ACCESS_TOKEN_MAXAGE"`
	RefreshTokenMaxAge     int           `mapstructure:"REFRESH_TOKEN_MAXAGE"`

	GoogleClientID         string `mapstructure:"GOOGLE_OAUTH_CLIENT_ID"`
	GoogleClientSecret     string `mapstructure:"GOOGLE_OAUTH_CLIENT_SECRET"`
	GoogleOAuthRedirectUrl string `mapstructure:"GOOGLE_OAUTH_REDIRECT_URL"`

	GitHubClientID         string `mapstructure:"GITHUB_OAUTH_CLIENT_ID"`
	GitHubClientSecret     string `mapstructure:"GITHUB_OAUTH_CLIENT_SECRET"`
	GitHubOAuthRedirectUrl string `mapstructure:"GITHUB_OAUTH_REDIRECT_URL"`

	FacebookClientID        string `mapstructure:"FACEBOOK_OAUTH_CLIENT_ID"`
	FacebookClientSecret    string `mapstructure:"FACEBOOK_OAUTH_CLIENT_SECRET"`
	FacebookRedirectUrl 	string `mapstructure:"FACEBOOK_OAUTH_REDIRECT_URL"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigFile("app.env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}