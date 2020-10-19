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

	"github.com/spf13/cobra"

	api "github.com/silinternational/terraform-enterprise-migrator/lib"
)

var keyContains string
var valueContains string

// variablesCmd represents the list command
var variablesCmd = &cobra.Command{
	Use:   "variables",
	Short: "Report on variables",
	Long:  `Show the values of variables with a key or value containing a certain string`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if organization == "" {
			fmt.Println("Error: The 'organization' flag is required")
			fmt.Println("")
			os.Exit(1)
		}
		if len(keyContains) == 0 && len(valueContains) == 0 {
			fmt.Println("Error: Either the 'key_contains' flag or 'value_contains flag must be set")
			fmt.Println("")
			os.Exit(1)
		}

		keyMsg := ""
		valMsg := ""

		if keyContains != "" {
			keyMsg = " key containing " + keyContains
		}

		if valueContains != "" {
			valMsg = " value containing " + valueContains
			if keyContains != "" {
				valMsg = " or value containing " + valueContains
			}
		}

		fmt.Printf("Getting variables from Terraform with%s%s\n",
			keyMsg, valMsg)
		println()
		runVariables()
	},
}

func init() {
	rootCmd.AddCommand(variablesCmd)

	variablesCmd.Flags().StringVarP(&organization, "organization", "o", "",
		"required - Name of Terraform Enterprise Organization")
	variablesCmd.Flags().StringVarP(&keyContains, "key_contains", "k", "",
		"required if value_contains is blank - string contained in the Terraform variable keys to report on")
	variablesCmd.Flags().StringVarP(&valueContains, "value_contains", "v", "",
		"required if key_contains is blank - string contained in the Terraform variable values to report on")
}

func runVariables() {
	allData, err := api.GetV2AllWorkspaceData(organization, atlasToken)
	if err != nil {
		println(err.Error())
		return
	}

	wsVars, err := api.GetAllWorkSpacesVarsFromV2(allData, organization, keyContains, valueContains, atlasToken)
	if err != nil {
		println(err.Error())
		return
	}

	indent := "   "

	for ws, vs := range wsVars {
		if len(vs) == 0 {
			fmt.Printf("%s has no matching variables\n\n", ws)
			continue
		}

		fmt.Printf("%s has %v matching variables ...\n", ws, len(vs))
		for _, v := range vs {
			sens := "not sensitive"
			if v.Sensitive {
				sens = "sensitive"
			}
			fmt.Printf("%s %s = %s (%s)\n", indent, v.Key, v.Value, sens)
		}
		println()
	}
	println()
	return
}
