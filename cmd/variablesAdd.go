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

var variablesAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add variables",
	Long:  `Add variables in matching workspaces having the specified key`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		runVariablesAdd()
	},
}

func init() {
	variablesCmd.AddCommand(variablesAddCmd)
	variablesAddCmd.Flags().StringVarP(&key, "key", "k", "",
		requiredPrefix+"Terraform variable key")
	if err := variablesAddCmd.MarkFlagRequired("key"); err != nil {
		errLog.Fatalln("failed to mark 'key' as a required flag on variablesAddCmd")
	}
	variablesAddCmd.Flags().StringVarP(&value, "value", "v", "",
		requiredPrefix+"Terraform variable value")
	if err := variablesAddCmd.MarkFlagRequired("value"); err != nil {
		errLog.Fatalln("failed to mark 'value' as a required flag on variablesAddCmd")
	}
}

func runVariablesAdd() {
	if readOnlyMode {
		fmt.Println("Read only mode enabled. No variables will be added.")
	}

	if workspace != "" {
		addWorkspaceVar(organization, workspace, key, value)
		return
	}

	fmt.Printf("Adding variables with key '%s' and value '%s' to all workspaces...\n", key, value)
	allWorkspaces, err := lib.GetAllWorkspaces(organization)
	if err != nil {
		println(err.Error())
		return
	}

	for _, w := range allWorkspaces {
		addWorkspaceVar(organization, w.Attributes.Name, key, value)
	}
	return
}

func addWorkspaceVar(org, ws, key, value string) bool {
	if v, err := lib.GetWorkspaceVar(org, ws, key); err != nil {
		println(err)
		return false
	} else if v != nil {
		fmt.Printf("'%s' already exists in '%s'. Use 'variable update' command to change the value.", key, value)
		return false
	}

	fmt.Printf("Adding variable '%s' to workspace '%s' with value '%s'\n", key, ws, value)
	if !readOnlyMode {
		if _, err := lib.AddOrUpdateVariable(lib.UpdateConfig{
			Organization:          organization,
			Workspace:             ws,
			SearchString:          key,
			NewValue:              value,
			AddKeyIfNotFound:      true,
			SearchOnVariableValue: false,
			SensitiveVariable:     false,
		}); err != nil {
			return false
		}
	}
	return true
}
