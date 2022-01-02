// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{Use: "eve"}

func init() {
	viper.SetEnvPrefix("eve")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	rootCmd.AddCommand(pxeCmd)
	pxeInit()
}

// Execute primary function for cobra
func Execute() {
	rootCmd.Execute()
}
