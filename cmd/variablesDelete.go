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

	"github.com/spf13/cobra"
)

var key string

var variablesDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete variable",
	Long:  `Delete variable in matching workspace having the specified key`,
	Args:  cobra.ExactArgs(0),
	Run:   runVariablesDelete,
}

func init() {
	variablesCmd.AddCommand(variablesDeleteCmd)
	variablesDeleteCmd.Flags().StringVarP(&key, "key", "k", "",
		requiredPrefix+"Terraform variable key to delete, must match exactly")
	cobra.CheckErr(variablesDeleteCmd.MarkFlagRequired("key"))
}

func runVariablesDelete(cmd *cobra.Command, args []string) {
	if workspace == "" {
		cobra.CheckErr("no workspace specified")
	}

	w, err := client.Workspaces.Read(ctx, organization, workspace)
	cobra.CheckErr(err)

	for _, v := range w.Variables {
		if v.Key == key {
			fmt.Printf("Deleting variable %s from workspace %s\n", v.Key, workspace)

			if !readOnlyMode {
				err := client.Variables.Delete(ctx, w.ID, v.ID)
				cobra.CheckErr(err)
			}
			return
		}
	}

	err = fmt.Errorf("variable %s not found in workspace %s", key, workspace)
	cobra.CheckErr(err)
}
