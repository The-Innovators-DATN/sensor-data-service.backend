package config

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	App        AppConfig
	DB         DBConfig
	Redis      RedisConfig
	Log        LogConfig
	Clickhouse ClickhouseConfig
}

func LoadAllConfigs(dir string) (*Config, error) {
	_ = godotenv.Load(filepath.Join(dir, ".env"))

	cfg := &Config{}
	load := []struct {
		Name   string
		Target interface{}
		EnvFn  func(*viper.Viper)
	}{
		{"app", &cfg.App, nil},
		{"database", &cfg.DB, func(v *viper.Viper) { overrideEnv(v, "database") }},
		{"redis", &cfg.Redis, func(v *viper.Viper) { overrideEnv(v, "redis") }},
		{"log", &cfg.Log, func(v *viper.Viper) { overrideEnv(v, "log") }},
		{"clickhouse", &cfg.Clickhouse, func(v *viper.Viper) { overrideEnv(v, "clickhouse") }},
	}
	log.Print(load)
	for _, item := range load {

		v := viper.New()
		v.SetConfigName(item.Name)
		v.SetConfigType("yaml")
		v.AddConfigPath(dir)
		v.AutomaticEnv()
		v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("error loading %s.yaml: %w", item.Name, err)
		}
		if item.EnvFn != nil {
			item.EnvFn(v)
			fmt.Printf("Final %s config:\n%+v\n", item.Name, item.Target)
		}
		if err := v.Unmarshal(item.Target); err != nil {
			return nil, fmt.Errorf("error decoding %s config: %w", item.Name, err)
		}
	}

	return cfg, nil
}
