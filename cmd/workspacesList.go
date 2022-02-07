// Copyright © 2018-2021 SIL International
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

	"github.com/spf13/cobra"

	"github.com/silinternational/tfc-ops/lib"
)

var attributes string

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List Workspaces",
	Long:  `Lists the TF workspaces with (some of) their attributes`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Getting list of workspaces ...")
		runList()
	},
}

func init() {
	workspaceCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&attributes, "attributes", "a", "",
		requiredPrefix+"Workspace attributes to list: "+strings.Join(lib.WorkspaceListAttributes, ", ")+
			" deprecated attributes: "+strings.Join(lib.WorkspaceListAttributesDeprecated, ", "))
	listCmd.MarkFlagRequired("attributes")
}

func runList() {
	allAttrs := strings.Split(attributes, ",")
	for _, attr := range allAttrs {
		if !lib.IsStringInSlice(attr, lib.WorkspaceListAttributes) {
			if lib.IsStringInSlice(attr, lib.WorkspaceListAttributesDeprecated) {
				fmt.Printf("DEPRECATION: attribute '%s' is deprecated and will be removed in a future version of this program\n",
					attr)
			} else {
				fmt.Printf("'%s' is not a valid workspace attribute\n", attr)
			}
		}
	}

	allData, err := lib.GetV2AllWorkspaceData(organization, atlasToken)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, ws := range allData {
		values := make([]string, len(allAttrs))
		for i, a := range allAttrs {
			value, err := ws.AttributeByLabel(strings.Trim(a, " "))
			if err != nil {
				println("\n", err.Error())
				return
			}
			values[i] = value
		}

		println(strings.Join(values, ", "))
	}
}
