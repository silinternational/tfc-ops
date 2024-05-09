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
	"strings"

	"github.com/hashicorp/go-tfe"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	workspace             string
	variableSearchString  string
	newVariableValue      string
	searchOnVariableValue bool
	addKeyIfNotFound      bool
	sensitiveVariable     bool
)

// cloneCmd represents the clone command
var variablesUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update/add a variable in a Workspace",
	Long:  `Update or add a variable in a Terraform Cloud Workspace based on a complete case-insensitive match`,
	Args:  cobra.ExactArgs(0),
	Run:   runVariablesUpdate,
}

func init() {
	variablesCmd.AddCommand(variablesUpdateCmd)
	variablesUpdateCmd.Flags().StringVarP(&variableSearchString, "variable-search-string", "s", "",
		requiredPrefix+`The string to match in the current variables (either in the Key or Value - see other flags)`)
	variablesUpdateCmd.Flags().StringVarP(&newVariableValue, "new-variable-value", "n", "",
		requiredPrefix+`The desired new value of the variable`)
	variablesUpdateCmd.Flags().BoolVarP(&addKeyIfNotFound, "add-key-if-not-found", "a", false,
		`optional (e.g. "-a=true") whether to add a new variable if a matching key is not found.`)
	variablesUpdateCmd.Flags().BoolVarP(&searchOnVariableValue, "search-on-variable-value", "v", false,
		`optional (e.g. "-v=true") whether to do the search based on the value of the variables. (Must be false if add-key-if-not-found is true`)
	variablesUpdateCmd.Flags().BoolVarP(&sensitiveVariable, "sensitive-variable", "x", false,
		`optional (e.g. "-x=true") make the variable sensitive.`)

	cobra.CheckErr(variablesUpdateCmd.MarkFlagRequired("variable-search-string"))
	cobra.CheckErr(variablesUpdateCmd.MarkFlagRequired("new-variable-value"))
	variablesUpdateCmd.MarkFlagsMutuallyExclusive("add-key-if-not-found", "search-on-variable-value")
}

func runVariablesUpdate(cmd *cobra.Command, args []string) {
	fmt.Printf("update variable called using %s, %s, search string: %s, new value: %s, add-key-if-not-found: %t, search-on-variable-value: %t\n",
		organization, workspace, variableSearchString, newVariableValue, addKeyIfNotFound, searchOnVariableValue)

	var ws []*tfe.Workspace
	if workspace != "" {
		w, err := client.Workspaces.Read(ctx, organization, workspace)
		cobra.CheckErr(err)

		ws = append(ws, w)
	} else {
		var err error
		list, err := client.Workspaces.List(ctx, organization, nil)
		cobra.CheckErr(err)

		ws = list.Items
	}

	for _, w := range ws {
		if len(ws) > 1 {
			fmt.Printf("Do you want to update the variable %s across the workspace: %s\n\n", variableSearchString, w.Name)
			if !awaitUserResponse() {
				continue
			}
		}

		addOrUpdateVariable(w, variableSearchString, searchOnVariableValue, addKeyIfNotFound, newVariableValue, sensitiveVariable)
	}
}

// addOrUpdateVariable adds or updates an existing Terraform Cloud workspace variable
// If the copyVariables param is set to true, then all the non-sensitive variable values will be added to the new
// workspace.  Otherwise, they will be set to "REPLACE_THIS_VALUE"
func addOrUpdateVariable(w *tfe.Workspace, search string, searchValue, addKeyIfNotFound bool, value string, sensitive bool) {
	for _, v := range w.Variables {
		oldValue := v.Value

		if (searchValue && !strings.EqualFold(v.Value, search)) ||
			(!searchValue && !strings.EqualFold(v.Key, search)) {
			continue
		}

		if !readOnlyMode {
			_, err := client.Variables.Update(ctx, w.ID, v.ID, tfe.VariableUpdateOptions{
				Value:     &newVariableValue,
				Sensitive: &sensitive,
			})
			cobra.CheckErr(err)
		}
		fmt.Printf("Replaced the value of %s from %s to %s\n", v.Key, oldValue, value)
		return
	}

	if addKeyIfNotFound {
		if !readOnlyMode {
			cat := tfe.CategoryTerraform
			_, err := client.Variables.Create(ctx, w.ID, tfe.VariableCreateOptions{
				Key:       &search,
				Value:     &value,
				Category:  &cat,
				Sensitive: &sensitive,
			})
			cobra.CheckErr(err)
		}
		fmt.Printf("Added variable %s = %s\n", search, value)
		return
	}

	fmt.Println("No match found and no variable added")
}

func awaitUserResponse() bool {
	prompt := promptui.Select{
		Label: "Select[Yes/No]",
		Items: []string{"No", "Yes"},
	}
	_, result, err := prompt.Run()
	cobra.CheckErr(err)

	return result == "Yes"
}
