package config

import (
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App struct {
		Port string `env-required:"true" yaml:"port"`
	} `yaml:"app"`
	CryptoAPI struct {
		URL        string `env-required:"true" yaml:"url"`
		HeaderName string `env-required:"true" yaml:"headerName"`
		APIKey     string `env-required:"true" env:"COINMARKETCAP_API_KEY"`
	} `yaml:"cryptoApi"`
	EmailSending struct {
		SenderAddress string `env-required:"true" yaml:"senderAddress"`
		PublicKey     string `env-required:"true" env:"MAILJET_PUBLIC_KEY"`
		PrivateKey    string `env-required:"true" env:"MAILJET_PRIVATE_KEY"`
	}
	Database struct {
		FilePath string `env-required:"true" yaml:"filePath"`
	} `yaml:"database"`
}

func GetConfig() *Config {
	var cfg = &Config{}
	var once sync.Once

	once.Do(func() {
		err := cleanenv.ReadConfig("config/config.yml", cfg)
		if err != nil {
			return
		}
	})

	return cfg
}
