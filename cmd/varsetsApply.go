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

var varsetsApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply Variable Set to Workspaces",
	Long:  `Apply an existing variable set to workspaces`,
	Args:  cobra.ExactArgs(0),
	Run:   runVarsetsApply,
}

func init() {
	varsetsCmd.AddCommand(varsetsApplyCmd)

	varsetsApplyCmd.Flags().StringVarP(&variableSet, "set", "s", "",
		requiredPrefix+"Terraform variable set to add")
	varsetsApplyCmd.Flags().StringVarP(&workspace, "workspace", "w", "",
		"Name of the Workspace in Terraform Cloud")
	varsetsApplyCmd.Flags().StringVar(&workspaceFilter, "workspace-filter", "",
		"Partial workspace name to search across all workspaces")

	cobra.CheckErr(varsetsApplyCmd.MarkFlagRequired("set"))
	varsetsApplyCmd.MarkFlagsOneRequired("workspace", "workspace-filter")

	// varsetsApplyCmd.Flags().VisitAll(func(f *pflag.Flag) {
	// 	if strings.HasPrefix(f.Usage, requiredPrefix) {
	// 		cobra.CheckErr(varsetsApplyCmd.MarkFlagRequired(f.Name))
	// 	}
	// })
}

func runVarsetsApply(cmd *cobra.Command, args []string) {
	var ws []*tfe.Workspace
	if workspace != "" {
		w, err := client.Workspaces.Read(ctx, organization, workspace)
		cobra.CheckErr(err)

		ws = append(ws, w)
	} else {
		var err error
		list, err := client.Workspaces.List(ctx, organization, &tfe.WorkspaceListOptions{Search: workspaceFilter})
		cobra.CheckErr(err)

		ws = list.Items
	}

	fmt.Printf("Applying variable set '%s' to %s\n", variableSet, workspaceListToString(ws))
	if !readOnlyMode {
		err := client.VariableSets.ApplyToWorkspaces(ctx, variableSet, &tfe.VariableSetApplyToWorkspacesOptions{Workspaces: ws})
		cobra.CheckErr(err)
	}
}
