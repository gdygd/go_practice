package main

import "github.com/spf13/viper"

type Config struct {
	OpenApiKey string `mapstructure:"OPENAI_API_KEY"`
}

func LoadConfig(path string) (Config, error) {
	var config Config
	var err error = nil
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	return config, nil
}
