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

var consumers []string

var workspaceConsumersAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add workspace remote state consumers",
	Args:  cobra.ExactArgs(0),
	Run:   runWorkspaceConsumersAdd,
}

func init() {
	workspaceConsumersCmd.AddCommand(workspaceConsumersAddCmd)

	workspaceConsumersAddCmd.Flags().StringVarP(&workspace, "workspace", "w", "",
		requiredPrefix+"Partial workspace name to search across all workspaces")
	workspaceConsumersAddCmd.Flags().StringSliceVar(&consumers, "consumers", nil,
		requiredPrefix+"List of remote state consumer workspaces, comma-separated")

	cobra.CheckErr(workspaceConsumersAddCmd.MarkFlagRequired("workspace"))
	cobra.CheckErr(workspaceConsumersAddCmd.MarkFlagRequired("consumers"))
}

func runWorkspaceConsumersAdd(cmd *cobra.Command, args []string) {
	ws, err := client.Workspaces.Read(ctx, organization, workspace)
	cobra.CheckErr(err)

	cs := make([]*tfe.Workspace, len(consumers))
	for i, consumer := range consumers {
		cs[i], err = client.Workspaces.Read(ctx, organization, consumer)
		cobra.CheckErr(err)
	}

	fmt.Printf("Adding %s consumer %s", workspace, workspaceListToString(cs))
	if !readOnlyMode {
		err := client.Workspaces.AddRemoteStateConsumers(ctx, ws.ID, tfe.WorkspaceAddRemoteStateConsumersOptions{Workspaces: cs})
		cobra.CheckErr(err)
	}
}
