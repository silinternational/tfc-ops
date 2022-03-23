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

	"github.com/silinternational/tfc-ops/lib"
)

var key string

var variablesRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove variables",
	Long:  `Remove variables in matching workspaces having the specified key`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if len(key) == 0 {
			fmt.Println("Error: The 'key' flag must be set")
			fmt.Println("")
			os.Exit(1)
		}

		runVariablesRemove()
	},
}

func init() {
	variablesCmd.AddCommand(variablesRemoveCmd)
	variablesRemoveCmd.Flags().StringVarP(&key, "key", "k", "",
		requiredPrefix+"Terraform variable key to remove, must match exactly")
}

func runVariablesRemove() {
	if dryRunMode {
		fmt.Println("Dry run mode enabled. No variables will be removed.")
	}

	if workspace != "" {
		found := removeWorkspaceVar(organization, workspace, key)
		if !found {
			fmt.Printf("Variable %s not found in workspace %s\n", key, workspace)
		}
		return
	}

	fmt.Printf("Removing variables with key '%s' from all workspaces...\n", key)
	allWorkspaces, err := lib.GetAllWorkspaces(organization, atlasToken)
	if err != nil {
		println(err.Error())
		return
	}

	for _, w := range allWorkspaces {
		removeWorkspaceVar(organization, w.Attributes.Name, key)
	}
	return
}

func removeWorkspaceVar(org, ws, key string) bool {
	v, err := lib.GetWorkspaceVar(org, ws, atlasToken, key)
	if err != nil {
		println(err.Error())
		return false
	}
	if v == nil {
		return false
	}

	fmt.Printf("Removing variable %s from workspace %s\n", v.Key, ws)
	if !dryRunMode {
		lib.DeleteVariable(v.ID, atlasToken)
	}
	return true
}
