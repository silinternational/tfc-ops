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

var workspaceConsumersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List workspace remote state consumers",
	Args:  cobra.ExactArgs(0),
	Run:   runWorkspaceConsumersList,
}

func init() {
	workspaceConsumersCmd.AddCommand(workspaceConsumersListCmd)

	workspaceConsumersListCmd.Flags().StringVarP(&workspace, "workspace", "w", "",
		requiredPrefix+"Partial workspace name to search across all workspaces")

	cobra.CheckErr(workspaceConsumersListCmd.MarkFlagRequired("workspace"))
}

func runWorkspaceConsumersList(cmd *cobra.Command, args []string) {
	ws, err := client.Workspaces.Read(ctx, organization, workspace)
	cobra.CheckErr(err)

	list, err := client.Workspaces.ListRemoteStateConsumers(ctx, ws.ID, nil)
	cobra.CheckErr(err)

	fmt.Printf("Workspace %s has consumer %s", ws.Name, workspaceListToString(list.Items))
}
