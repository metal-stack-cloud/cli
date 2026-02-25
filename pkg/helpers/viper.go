package helpers

import "slices"

import "github.com/spf13/viper"

func IsAnyViperFlagSet(names ...string) bool {
	return slices.ContainsFunc(names, viper.IsSet)
}
