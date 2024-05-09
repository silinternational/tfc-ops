// Copyright © 2018-2024 SIL International
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	requiredPrefix = "required - "
	atlasToken     = "ATLAS_TOKEN"
)

var (
	cfgFile      string
	token        string
	organization string
	readOnlyMode bool
	debug        bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:              "tfc-ops",
	Short:            "Terraform Cloud operations",
	Long:             `Perform TF Cloud operations, e.g. clone a workspace or manage variables within a workspace`,
	PersistentPreRun: initRoot,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initRoot(cmd *cobra.Command, args []string) {
	// Skip for version command
	if cmd.UseLine() == versionCmd.UseLine() {
		return
	}

	// Get Tokens from env vars
	token = viper.GetString(atlasToken)
	if token == "" {
		err := fmt.Errorf("environment variable for %s is required to execute plan and migration", atlasToken)
		cobra.CheckErr(err)
	}

	debug = viper.GetBool("TFC_OPS_DEBUG")

	if readOnlyMode {
		fmt.Println("###### READ ONLY MODE ENABLED ######")
	}

	cobra.CheckErr(NewClient(""))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".tfc-ops" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".tfc-ops")
	}

	viper.AutomaticEnv() // read in environment variables that match
	viper.SetDefault("TFC_OPS_DEBUG", false)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func addGlobalFlags(command *cobra.Command) {
	command.PersistentFlags().BoolVarP(&readOnlyMode, "read-only-mode", "r", false,
		`read-only mode (e.g. "-r")`,
	)

	command.PersistentFlags().StringVarP(&organization, "organization",
		"o", "", requiredPrefix+"Name of Terraform Cloud Organization")
	if err := command.MarkPersistentFlagRequired("organization"); err != nil {
		panic("MarkPersistentFlagRequired failed with error " + err.Error())
	}
}
