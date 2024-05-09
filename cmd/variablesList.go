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
	"os"
	"strings"
	"text/tabwriter"

	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
)

var (
	keyContains   string
	valueContains string
	tabularCSV    bool
	sensitive     = map[bool]string{true: "(sensitive)"}
)

var variablesListCmd = &cobra.Command{
	Use:   "list",
	Short: "Report on variables",
	Long:  `Show the values of variables with a key or value containing a certain string`,
	Args:  cobra.ExactArgs(0),
	Run:   runVariablesList,
}

func init() {
	variablesCmd.AddCommand(variablesListCmd)
	variablesListCmd.Flags().StringVarP(&keyContains, "key_contains", "k", "",
		"required if value_contains is blank - string contained in the Terraform variable keys to report on")
	variablesListCmd.Flags().StringVarP(&valueContains, "value_contains", "v", "",
		"required if key_contains is blank - string contained in the Terraform variable values to report on")
	variablesListCmd.Flags().BoolVar(&tabularCSV, "csv", false,
		"output variable list in CSV format")

	variablesListCmd.MarkFlagsOneRequired("key_contains", "value_contains")
}

func runVariablesList(cmd *cobra.Command, args []string) {
	if tabularCSV {
		fmt.Println("workspace,key,value")
	} else {
		msg := ""
		ws := workspace
		if ws == "" {
			ws = "all workspaces"
		}

		if keyContains != "" {
			msg = " key containing " + keyContains
		}

		if valueContains != "" {
			if msg != "" {
				msg += " or"
			}
			msg += " value containing " + valueContains
		}

		fmt.Printf("Getting variables from %s with%s\n", ws, msg)
	}

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
		var vars []*tfe.Variable
		for _, v := range w.Variables {
			if (keyContains != "" && strings.Contains(v.Key, keyContains)) ||
				(valueContains != "" && strings.Contains(v.Value, valueContains)) {
				vars = append(vars, v)
			}
		}

		if tabularCSV {
			printWorkspaceVarsCSV(w.Name, vars)
		} else {
			printWorkspaceVars(w.Name, vars)
		}
	}
}

func printWorkspaceVars(workspaceName string, vars []*tfe.Variable) {
	fmt.Printf("\nWorkspace %s has %v matching variable(s)\n", workspaceName, len(vars))
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "Key \t Value \t Sensitive")
	for _, v := range vars {
		fmt.Fprintf(w, "%s \t %s \t %s\n", v.Key, v.Value, sensitive[v.Sensitive])
	}
	w.Flush()
	fmt.Println()
}

func printWorkspaceVarsCSV(ws string, vs []*tfe.Variable) {
	for _, v := range vs {
		if v.Sensitive {
			v.Value = sensitive[v.Sensitive]
		}
		fmt.Printf(`"%s","%s","%s"`+"\n", escapeString(ws), escapeString(v.Key), escapeString(v.Value))
	}
}

// escapeString escapes characters for CSV encoding, adding a backslash before a double-quote, and converting
// a newline to `\n`
func escapeString(s string) string {
	tmp := strings.ReplaceAll(s, `"`, `\"`)
	return strings.ReplaceAll(tmp, "\n", `\n`)
}
