package config

import (
	"fmt"
	"path/filepath"
	"runtime"
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
	} `yaml:"emailSending"`
	Database struct {
		FilePath string `env-required:"true" yaml:"filePath"`
	} `yaml:"database"`
}

var cfg = &Config{}
var once sync.Once

func GetConfig() *Config {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	once.Do(func() {
		err := cleanenv.ReadConfig(filepath.Join(basepath, "config.yml"), cfg)
		if err != nil {
			fmt.Println(err)
			return
		}
	})

	return cfg
}
