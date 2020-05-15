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
	"strings"

	"github.com/spf13/cobra"
	api "github.com/silinternational/terraform-enterprise-migrator/lib"
)

var attributes string

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List Workspaces",
	Long: `Lists the TF workspaces with (some of) their attributes`,
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if organization == "" {
			fmt.Println("Error: The 'organization' flag is required")
			fmt.Println("")
			os.Exit(1)
		}
		if len(attributes) == 0 {
			fmt.Println("Error: The 'attributes' flag is required")
			fmt.Println("")
			os.Exit(1)
		}
		fmt.Println("Getting list of workspaces ...")
		runList()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&organization, "organization", "o", "",
		"required - Name of Terraform Enterprise Organization")
	listCmd.Flags().StringVarP(&attributes, "attributes", "a", "",
		"required - Workspace attributes to list: id,name,createdat,environment,workingdirectory,terraformversion,vcsrepo")

}

func runList() {
	allData, err := api.GetV2AllWorkspaceData(organization, atlasToken)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	allAttrs := strings.Split(attributes, ",")

	for _, ws := range allData {
		for _, a := range allAttrs {
			value, err := ws.AttributeByLabel(strings.Trim(a, " "))
			if err != nil {
				println("\n", err.Error())
				return
			}
			print(value, ", ")
		}

		println()
	}
}