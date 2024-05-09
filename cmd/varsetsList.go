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

var varsetsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Variable Sets",
	Long:  `List variable sets applied to a workspace`,
	Args:  cobra.ExactArgs(0),
	Run:   runVarsetsList,
}

func init() {
	varsetsCmd.AddCommand(varsetsListCmd)

	varsetsListCmd.Flags().StringVarP(&workspace, "workspace", "w", "",
		"Name of the Workspace in Terraform Cloud")
	varsetsListCmd.Flags().StringVar(&workspaceFilter, "workspace-filter", "",
		"Partial workspace name to search across all workspaces")

	varsetsApplyCmd.MarkFlagsOneRequired("workspace", "workspace-filter")
}

func runVarsetsList(cmd *cobra.Command, args []string) {
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

	for _, w := range ws {
		list, err := client.VariableSets.ListForWorkspace(ctx, w.ID, nil)
		cobra.CheckErr(err)

		fmt.Printf("Workspace %s has the following variable sets:\n", w.Name)
		for _, set := range list.Items {
			fmt.Printf("  %s\n", set.Name)
		}
	}
}
