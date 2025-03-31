package config

type AppConfig struct {
	HostPort string `mapstructure:"hostport"`
	Debug    bool   `mapstructure:"debug"`
	Env      string `mapstructure:"env"`
}
