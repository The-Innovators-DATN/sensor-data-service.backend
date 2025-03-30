package config

type S3Config struct {
	Region    string `mapstructure:"region"`
	Bucket    string `mapstructure:"bucket"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"scret_key"`
	Endpoint  string `mapstructure:"endpoint"`
}
