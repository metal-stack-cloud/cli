package helpers

import "github.com/spf13/viper"

func IsAnyViperFlagSet(names ...string) bool {
	for _, name := range names {
		if viper.IsSet(name) {
			return true
		}
	}
	return false
}
