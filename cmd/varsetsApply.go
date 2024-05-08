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

	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"

	"github.com/silinternational/tfc-ops/v3/lib"
)

var varsetsApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply Variable Set to Workspaces",
	Long:  `Apply an existing variable set to workspaces`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		runVarsetsApply(variableSet)
	},
}

func init() {
	varsetsCmd.AddCommand(varsetsApplyCmd)

	varsetsApplyCmd.Flags().StringVarP(&variableSet, "set", "s", "",
		requiredPrefix+"Terraform variable set to add")
	if err := varsetsApplyCmd.MarkFlagRequired("set"); err != nil {
		errLog.Fatalln("failed to mark 'set' as a required flag on varsetsApplyCmd")
	}

	varsetsApplyCmd.Flags().StringVarP(&workspace, "workspace", "w", "",
		"Name of the Workspace in Terraform Cloud")

	varsetsApplyCmd.Flags().StringVar(&workspaceFilter, "workspace-filter", "",
		"Partial workspace name to search across all workspaces")
}

func runVarsetsApply(name string) {
	if readOnlyMode {
		fmt.Println("Read only mode enabled. No variable set will be applied.")
	}

	if workspace == "" && workspaceFilter == "" {
		errLog.Fatalln("Either --workspace or --workspace-filter must be specified.")
	}

	var err error
	var workspaces []*tfe.Workspace
	if workspace != "" {
		w, err := lib.GetWorkspaceByName(organization, workspace)
		if err != nil {
			errLog.Fatalf("error getting workspace from Terraform: %s", err)
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

	_ = applyVariableSet(organization, name, workspaces)
}

func applyVariableSet(org, vsName string, workspaces []*tfe.Workspace) bool {
	vs, err := lib.GetVariableSet(org, vsName)
	if err != nil {
		errLog.Fatalf("Error retrieving variable set: %s", err)
	}
	if vs == nil {
		errLog.Fatalf("No variable set matches the name given (%s)", vsName)
	}

	fmt.Printf("Applying variable set '%s' to %s\n", vs.Name, lib.WorkspaceListToString(workspaces))
	if err = lib.ApplyVariableSet(vs.ID, workspaces); err != nil {
		errLog.Fatalf("Error while applying variable set: %s", err)
	}
	return true
}
