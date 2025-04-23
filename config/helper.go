package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func overrideEnv(v *viper.Viper, prefix string) {
	for _, key := range v.AllKeys() {
		envKey := strings.ToUpper(fmt.Sprintf("%s_%s", prefix, strings.ReplaceAll(key, ".", "_")))
		if val, ok := os.LookupEnv(envKey); ok {
			v.Set(key, val)
		}
	}
}
