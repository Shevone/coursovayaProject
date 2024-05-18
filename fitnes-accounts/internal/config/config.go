package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	StoragePath string        `yaml:"storage_path" env-required:"true"` // 2 - если будет пуст, то приложение не запустится
	GRPC        GRPCConfig    `yaml:"grpc" env-required:"true"`
	TokenTTL    time.Duration `yaml:"tokenTTL" env-default:"1h"`
}
type GRPCConfig struct {
	Port    int           `yaml:"port" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env-default:"5s"`
}

// MustLoad Must - функция не будет возвращать ошибку, а будет паниковать
func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath)
}
func MustLoadPath(configPath string) *Config {
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {

		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

// Метод, который будет искать путь до файла конфигурации
// 1. из флагов при запуске программы через консоль
// 2. из переменных окружения
// upd берет чисто переменные окружения
func fetchConfigPath() string {
	return os.Getenv("CONFIG_PATH")
}
