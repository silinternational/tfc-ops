// Copyright © 2018-2022 SIL International
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

const requiredPrefix = "required - "

var (
	atlasToken   string
	cfgFile      string
	debug        bool
	organization string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tfc-ops",
	Short: "Terraform Cloud operations",
	Long:  `Perform TF Cloud operations, e.g. clone a workspace or manage variables within a workspace`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
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

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	foundError := false

	// Get Tokens from env vars
	atlasToken = os.Getenv("ATLAS_TOKEN")
	if atlasToken == "" {
		fmt.Println("Error: Environment variable for ATLAS_TOKEN is required to execute plan and migration")
		fmt.Println("")
		foundError = true
	}

	debugStr := os.Getenv("TFC_OPS_DEBUG")
	if debugStr == "TRUE" || debugStr == "true" {
		debug = true
	}

	if foundError {
		os.Exit(1)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".tfc-ops" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".tfc-ops")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func makeOrgFlagRequired(command *cobra.Command) {
	command.PersistentFlags().StringVarP(&organization, "organization",
		"o", "", requiredPrefix+"Name of Terraform Enterprise Organization")
	if err := command.MarkPersistentFlagRequired("organization"); err != nil {
		panic("MarkPersistentFlagRequired failed with error " + err.Error())
	}
}
