// Copyright © 2023 SIL International
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

	"github.com/silinternational/tfc-ops/v3/lib"
)

var varsetsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Variable Sets",
	Long:  `List variable sets applied to a workspace`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		runVarsetsList()
	},
}

func init() {
	varsetsCmd.AddCommand(varsetsListCmd)

	varsetsListCmd.Flags().StringVarP(&workspace, "workspace", "w", "",
		"Name of the Workspace in Terraform Cloud")

	varsetsListCmd.Flags().StringVar(&workspaceFilter, "workspace-filter", "",
		"Partial workspace name to search across all workspaces")
}

func runVarsetsList() {
	if workspace == "" && workspaceFilter == "" {
		errLog.Fatalln("Either --workspace or --workspace-filter must be specified.")
	}

	var err error
	var workspaces []*tfe.Workspace
	if workspace != "" {
		w, err := lib.GetWorkspaceByName(organization, workspace)
		if err != nil {
			errLog.Fatalf("error getting workspace %q from Terraform: %s", workspace, err)
		}
		workspaces = append(workspaces, w)
	} else {
		workspaces, err = lib.FindWorkspaces(organization, workspaceFilter)
		if err != nil {
			errLog.Fatalf(err.Error())
		}
		if len(workspaces) == 0 {
			errLog.Fatalf("no workspaces match the filter '%s'", workspaceFilter)
		}
	}

	for _, w := range workspaces {
		sets, err := lib.GetWorkspaceVariableSets(w.ID)
		if err != nil {
			errLog.Fatalf(err.Error())
		}
		fmt.Printf("Workspace %s has the following variable sets:\n", w.Name)
		for _, set := range sets {
			fmt.Printf("  %s\n", set.Name)
		}
	}
}
