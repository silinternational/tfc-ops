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

	"github.com/spf13/cobra"

	"github.com/silinternational/tfc-ops/lib"
)

var key string

var variablesDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete variables",
	Long:  `Delete variables in matching workspaces having the specified key`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		runVariablesDelete()
	},
}

func init() {
	variablesCmd.AddCommand(variablesDeleteCmd)
	variablesDeleteCmd.Flags().StringVarP(&key, "key", "k", "",
		requiredPrefix+"Terraform variable key to delete, must match exactly")
	if err := variablesDeleteCmd.MarkFlagRequired("key"); err != nil {
		errLog.Fatalln("failed to mark 'key' as a required flag on variablesDeleteCmd: " + err.Error())
	}
}

func runVariablesDelete() {
	if dryRunMode {
		fmt.Println("Dry run mode enabled. No variables will be deleted.")
	}

	if workspace != "" {
		found := deleteWorkspaceVar(organization, workspace, key)
		if !found {
			fmt.Printf("Variable %s not found in workspace %s\n", key, workspace)
		}
		return
	}

	fmt.Printf("Deleting variables with key '%s' from all workspaces...\n", key)
	allWorkspaces, err := lib.GetAllWorkspaces(organization)
	if err != nil {
		println(err.Error())
		return
	}

	for _, w := range allWorkspaces {
		deleteWorkspaceVar(organization, w.Attributes.Name, key)
	}
	return
}

func deleteWorkspaceVar(org, ws, key string) bool {
	v, err := lib.GetWorkspaceVar(org, ws, key)
	if err != nil {
		println(err.Error())
		return false
	}
	if v == nil {
		return false
	}

	fmt.Printf("Deleting variable %s from workspace %s\n", v.Key, ws)
	if !dryRunMode {
		lib.DeleteVariable(v.ID)
	}
	return true
}
