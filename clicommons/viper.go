package clicommons

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// MustBindPFlag calls viper.BindPFlag and panics if there is an error. This is a useful wrapper for properly handling
// errors in situations where we almost certainly know the binding will succeed.
func MustBindPFlag(cfgName string, flag *pflag.Flag) {
	err := viper.BindPFlag(cfgName, flag)
	if err != nil {
		panic(err)
	}
}
