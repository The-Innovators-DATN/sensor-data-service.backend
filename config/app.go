package config

type AppConfig struct {
	Port  int    `mapstructure:"port"`
	Host  string `mapstructure:"host"`
	Debug bool   `mapstructure:"debug"`
	Env   string `mapstructure:"env"`
}
