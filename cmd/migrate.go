// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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
	migrator "github.com/silinternational/terraform-enterprise-migrator/lib"
	"github.com/spf13/cobra"
	"os"
)

var vcsUsername string

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Perform migration plan",
	Long:  `Processes plan.csv to validate migration plan and perform the work`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if vcsUsername == "" {
			fmt.Println("Error: The 'vcs-username' flag is required\n")
			os.Exit(1)
		}
		runMigration(vcsUsername, planFile)
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.Flags().StringVarP(
		&vcsUsername,
		"vcs-username",
		"u",
		"",
		`Name of the VCS User in TF Enterprise (new version) to allow us to get the right VCS Token ID (the GitHub or Bitbucket username used to connect TFE)`,
	)
	migrateCmd.Flags().StringVarP(&planFile, "file", "f", "plan.csv", "optional - Name of migration plan CSV file")
}

func runMigration(vcsUserName, planFile string) {
	fmt.Println("migrate called using config file: " + planFile)
	completed, err := migrator.CreateAndPopulateAllV2Workspaces(planFile, atlasToken, vcsUsername)
	if err != nil {
		fmt.Println(err.Error())
	}
	println("\n  **** Completed Workspaces and their Sensitive Variables ****")
	for workspace, sensitiveVars := range completed {
		println("> " + workspace)
		for _, nextVar := range sensitiveVars {
			println("   - " + nextVar)
		}
		println("")
	}
}

