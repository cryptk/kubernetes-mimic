package sources

import (
	"fmt"

	"github.com/spf13/viper"
)

func validateConfigSet(flag string, configs []string) (bool, error) {
	if !viper.GetBool(flag) {
		return false, nil
	}

	for _, value := range configs {
		if !viper.IsSet(value) || viper.GetString(value) == "" {
			return false, fmt.Errorf("invalid configuration detected, %s is set but %s is not", flag, value)
		}
	}

	return true, nil
}
