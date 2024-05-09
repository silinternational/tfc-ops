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
	"strings"

	"github.com/spf13/cobra"
)

var attributes []string

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List Workspaces",
	Long:  `Lists the TF workspaces with (some of) their attributes`,
	Args:  cobra.ExactArgs(0),
	Run:   runWorkspacesList,
}

func init() {
	workspaceCmd.AddCommand(listCmd)
	listCmd.Flags().StringSliceVarP(&attributes, "attributes", "a", nil,
		requiredPrefix+"Workspace attributes to list, use Terraform Cloud API workspace attribute names")

	cobra.CheckErr(listCmd.MarkFlagRequired("attributes"))
}

func runWorkspacesList(cmd *cobra.Command, args []string) {
	fmt.Println("Getting list of workspaces ...")
	list, err := client.Workspaces.List(ctx, organization, nil)
	cobra.CheckErr(err)

	fmt.Println(strings.Join(attributes, ", "))
	for _, w := range list.Items {
		r := reflect.ValueOf(w)
		data := make([]string, len(attributes))
		for i, a := range attributes {
			data[i] = reflect.Indirect(r).FieldByName(a).String()
		}
		fmt.Println(strings.Join(data, ", "))
	}
}
