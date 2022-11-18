package config

import "github.com/spf13/viper"

type Config struct {
	Port       string `mapstructure:"PORT"`
	DBUrl      string `mapstructure:"DB_URL"`
	DbAuthdb   string `mapstructure:"DB_AUTHDB"`
	DbTable    string `mapstructure:"DB_TABLE"`
	DbUser     string `mapstructure:"DB_USER"`
	DbPassword string `mapstructure:"DB_PASSWORD"`
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath("")
	viper.SetConfigFile("dev.env")
	viper.SetConfigName("dev")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	return
}
