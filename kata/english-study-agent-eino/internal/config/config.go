package config

import (
	"github.com/spf13/viper"
)

// FeedConfig represents a single RSS feed configuration
type FeedConfig struct {
	Title    string `mapstructure:"title" json:"title"`
	URL      string `mapstructure:"url" json:"url"`
	Category string `mapstructure:"category" json:"category,omitempty"`
}

type Config struct {
	AI struct {
		Provider string `mapstructure:"provider"`
		APIKey   string `mapstructure:"api_key"`
		Model    string `mapstructure:"model"`
		BaseURL  string `mapstructure:"base_url"`
		SkipTLS  bool   `mapstructure:"skip_tls"`
	} `mapstructure:"ai"`
	Server struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`
	Feeds []FeedConfig `mapstructure:"feeds"`
}

var GlobalConfig Config

func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.english-agent")

	viper.SetDefault("server.port", "8080")
	viper.SetDefault("ai.provider", "openai")
	viper.SetDefault("ai.model", "deepseek-chat")
	viper.SetDefault("ai.base_url", "https://api.deepseek.com")
	viper.SetDefault("ai.skip_tls", false)
	
	viper.AutomaticEnv() // Read from env vars

	// Map env vars like LLM_API_KEY to ai.api_key
	viper.BindEnv("ai.api_key", "LLM_API_KEY")
	viper.BindEnv("ai.base_url", "LLM_BASE_URL")
	viper.BindEnv("ai.model", "LLM_MODEL")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			// Config file was found but another error was produced
			panic(err)
		}
	}

	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		panic(err)
	}
}

func Get() *Config {
	return &GlobalConfig
}

