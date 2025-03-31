package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	App        AppConfig
	DB         DBConfig
	Clickhouse ClickhouseConfig
	Redis      RedisConfig
}

func LoadAllConfigs(dir string) (*Config, error) {

	_ = godotenv.Load(dir + "/.env")

	cfg := &Config{}
	files := map[string]interface{}{
		"app":        &cfg.App,
		"database":   &cfg.DB,
		"clickhouse": &cfg.Clickhouse,
		"redis":      &cfg.Redis,
	}

	for name, target := range files {
		viper := viper.New()
		viper.SetConfigName(name)
		viper.SetConfigType("yaml")
		viper.AddConfigPath(dir)
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		if err := viper.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		if err := viper.Unmarshal(target); err != nil {
			return nil, fmt.Errorf("failed to unmarshal config: %w", err)
		}
		if name == "database" {
			overrideDBEnv(&cfg.DB, viper)
		}
		if name == "clickhouse" {
			overrideClickhouseEnv(&cfg.Clickhouse, viper)
		}
		if name == "redis" {
			overrideRedisEnv(&cfg.Redis, viper)
		}
	}
	return cfg, nil
}
func overrideDBEnv(dbCfg *DBConfig, v *viper.Viper) {
	if val := os.Getenv("DATABASE_USER"); val != "" {
		dbCfg.User = val
	}
	if val := os.Getenv("DATABASE_PASSWORD"); val != "" {
		dbCfg.Password = val
	}
	if val := os.Getenv("DATABASE_HOST"); val != "" {
		dbCfg.Host = val
	}
	if val := os.Getenv("DATABASE_PORT"); val != "" {
		dbCfg.Port = v.GetInt("DATABASE_PORT")
	}
	if val := os.Getenv("DATABASE_NAME"); val != "" {
		dbCfg.Name = val
	}
}

func overrideClickhouseEnv(dbCfg *ClickhouseConfig, v *viper.Viper) {
	if val := os.Getenv("CLICKHOUSE_USER"); val != "" {
		dbCfg.User = val
	}
	if val := os.Getenv("CLICKHOUSE_PASSWORD"); val != "" {
		dbCfg.Password = val
	}
	if val := os.Getenv("CLICKHOUSE_ADDR"); val != "" {
		dbCfg.Addr = val
	}
	if val := os.Getenv("CLICKHOUSE_NAME"); val != "" {
		dbCfg.Name = val
	}
}

func overrideRedisEnv(dbCfg *RedisConfig, v *viper.Viper) {
	if val := os.Getenv("REDIS_ADDR"); val != "" {
		dbCfg.Addr = val
	}
	if val := os.Getenv("REDIS_PASSWORD"); val != "" {
		dbCfg.Password = val
	}
	if val := os.Getenv("REDIS_DB"); val != "" {
		dbCfg.DB = v.GetInt("REDIS_DB")
	}
	if val := os.Getenv("REDIS_PROTOCOL"); val != "" {
		dbCfg.Protocol = v.GetInt("REDIS_PROTOCOL")
	}
}
