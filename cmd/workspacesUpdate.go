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
	"reflect"

	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
)

var (
	attribute       string
	value           string
	workspaceFilter string
)

// workspaceUpdateCmd represents the workspace update command
var workspaceUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Workspaces",
	Long:  `Updates an attribute of Terraform workspaces`,
	Args:  cobra.ExactArgs(0),
	Run:   runWorkspaceUpdate,
}

func init() {
	workspaceCmd.AddCommand(workspaceUpdateCmd)

	workspaceUpdateCmd.Flags().StringVarP(&attribute, "attribute", "a", "",
		requiredPrefix+"Workspace attribute to update, use Terraform Cloud API workspace attribute names")
	workspaceUpdateCmd.Flags().StringVarP(&value, "value", "v", "",
		requiredPrefix+"Value")
	workspaceUpdateCmd.Flags().StringVarP(&workspaceFilter, "workspace", "w", "",
		requiredPrefix+"Partial workspace name to search across all workspaces")

	cobra.CheckErr(workspaceUpdateCmd.MarkFlagRequired("attribute"))
	cobra.CheckErr(workspaceUpdateCmd.MarkFlagRequired("value"))
	cobra.CheckErr(workspaceUpdateCmd.MarkFlagRequired("workspace"))
}

func runWorkspaceUpdate(cmd *cobra.Command, args []string) {
	fmt.Println("Updating workspaces ...")

	list, err := client.Workspaces.List(ctx, organization, &tfe.WorkspaceListOptions{Search: workspaceFilter})
	cobra.CheckErr(err)

	var opts tfe.WorkspaceUpdateOptions
	r := reflect.ValueOf(opts)
	f := reflect.Indirect(r).FieldByName(attribute)
	f.Set(reflect.ValueOf(value))

	for _, w := range list.Items {
		fmt.Printf("set '%s' to '%s' on workspace %s\n", attribute, value, w.Name)
		if !readOnlyMode {
			_, err := client.Workspaces.UpdateByID(ctx, w.ID, opts)
			cobra.CheckErr(err)
		}
	}
}
