package config

import (
	"github.com/jinzhu/configor"
)

// Структура конфигурации приложения
type Config struct {
	ConfigFile string
	LogLevel   uint32 `default:"4"`
	Api        struct {
		HttpPort int  `default:"8080"`
		Logging  bool `default:"false"`
	}
	Store struct {
		Host     string `required:"true"`
		Port     int    `required:"true"`
		User     string `required:"true"`
		Password string `required:"true"`
		Dbname   string `required:"true"`
	}
}

// Создать конфигурацию из файла configFile
func NewConfig(configFile string) (*Config, error) {
	config := &Config{ConfigFile: configFile}
	if err := configor.Load(config, configFile); err != nil {
		return nil, err
	}
	return config, nil
}
