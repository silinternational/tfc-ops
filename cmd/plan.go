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
	"os"

	migrator "github.com/silinternational/terraform-enterprise-migrator/lib"
	"github.com/spf13/cobra"
)

var legacyOrg string
var newOrg string
var planFile string

// planCmd represents the plan command
var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Generate migration plan file",
	Long: `Generates a plan.csv file with list of environments from legacy organization
for mapping to new organization.`,
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if legacyOrg == "" {
			fmt.Println("Error: The 'legacy' flag is required")
			fmt.Println("")
			os.Exit(1)
		}
		if newOrg == "" {
			fmt.Println("Error: The 'new' flag is required")
			fmt.Println("")
			os.Exit(1)
		}
		fmt.Println("Creating migration plan...")
		envList := migrator.GetAllEnvNamesFromV1API(atlasToken)
		plans := migrator.GetBasePlansFromEnvNames(envList, legacyOrg, newOrg)
		migrator.CreatePlanFile(planFile, plans)
		fmt.Println(envList)
	},
}

func init() {
	rootCmd.AddCommand(planCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// planCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	planCmd.Flags().StringVarP(&legacyOrg, "legacy", "l", "", "required - Name of Terraform Enterprise Legacy Organization")
	planCmd.Flags().StringVarP(&newOrg, "new", "n", "", "required - Name of new Terraform Enterprise Organization")
	planCmd.Flags().StringVarP(&planFile, "file", "f", "plan.csv", "optional - Name of migration plan CSV file")
}
