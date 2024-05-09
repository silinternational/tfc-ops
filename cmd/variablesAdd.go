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

	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
)

var variablesAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add new variable (if not already present)",
	Long:  `Add variable in matching workspace. Will not update existing variable.`,
	Args:  cobra.ExactArgs(0),
	Run:   runVariablesAdd,
}

func init() {
	variablesCmd.AddCommand(variablesAddCmd)
	variablesAddCmd.Flags().StringVarP(&key, "key", "k", "",
		requiredPrefix+"Terraform variable key")
	variablesAddCmd.Flags().StringVarP(&value, "value", "v", "",
		requiredPrefix+"Terraform variable value")

	cobra.CheckErr(variablesAddCmd.MarkFlagRequired("key"))
	cobra.CheckErr(variablesAddCmd.MarkFlagRequired("value"))
}

func runVariablesAdd(cmd *cobra.Command, args []string) {
	var ws []*tfe.Workspace
	if workspace != "" {
		w, err := client.Workspaces.Read(ctx, organization, workspace)
		cobra.CheckErr(err)

		ws = append(ws, w)
	} else {
		var err error
		list, err := client.Workspaces.List(ctx, organization, nil)
		cobra.CheckErr(err)

		workspace = "all workspaces"
		ws = list.Items
	}
	fmt.Printf("Adding variables with key '%s' and value '%s' to %s...\n", key, value, workspace)

	for _, w := range ws {
		for _, v := range w.Variables {
			if v.Key == key {
				err := fmt.Errorf("'%s' already exists in '%s'. Use 'variable update' command to change the value", key, w.Name)
				cobra.CheckErr(err)
			}
		}

		fmt.Printf("Workspace %s: Adding variable %s = %s\n", w.Name, key, value)
		if !readOnlyMode {
			_, err := client.Variables.Create(ctx, w.ID, tfe.VariableCreateOptions{
				Key:   &key,
				Value: &value,
			})
			cobra.CheckErr(err)
		}
	}
}
