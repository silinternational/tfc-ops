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
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Perform migration plan",
	Long:  `Processes plan.csv to validate migration plan and peform the work`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		runMigration(planFile)
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// planCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	migrateCmd.Flags().StringVarP(&planFile, "file", "f", "plan.csv", "optional - Name of migration plan CSV file")
}

func runMigration(planFile string) {
	fmt.Println("migrate called using config file: " + planFile)
	err := migrator.CreateAndPopulateAllV2Workspaces(planFile, atlasToken, vcsToken)
	if err != nil {
		fmt.Println(err.Error())
	}

}
