package config

type LogConfig struct {
	Level  string `mapstructure:"level"`  // e.g. debug, info, warn, error
	Pretty bool   `mapstructure:"pretty"` // enable colored console?
	Output string `mapstructure:"output"` // stdout, file, etc.
}
