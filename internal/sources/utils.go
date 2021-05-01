package sources

import "github.com/spf13/viper"

func validateConfigSet(flag string, configs []string) bool {
	if !viper.GetBool(flag) {
		return false
	}
	for _, value := range configs {
		if !viper.IsSet(value) {
			return false
		}
	}
	return true
}
