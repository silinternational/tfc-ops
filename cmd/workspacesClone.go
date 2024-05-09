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
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
)

type cloneConfig struct {
	CopyState                   bool
	CopyVariables               bool
	ApplyVariableSets           bool
	DifferentDestinationAccount bool

	Source struct {
		Organization  string
		WorkspaceName string
	}
	Destination struct {
		Organization  string
		WorkspaceName string
		VCSTokenID    string
		Token         string
	}
}

func (c cloneConfig) clone() []string {
	source, err := client.Workspaces.Read(ctx, c.Source.Organization, c.Source.WorkspaceName)
	cobra.CheckErr(err)

	if !c.DifferentDestinationAccount {
		c.Destination.Organization = organization
		c.Destination.VCSTokenID = source.VCSRepo.Identifier
	}

	sensitiveVars := []string{}
	const sensitiveValue = "TF_ENTERPRISE_SENSITIVE_VAR"
	const defaultValue = "REPLACE_THIS_VALUE"

	for _, v := range source.Variables {
		if !c.CopyVariables {
			v.Value = defaultValue
		}

		if v.Value == sensitiveValue {
			sensitiveVars = append(sensitiveVars, v.Key)
		}
	}

	if readOnlyMode {
		return sensitiveVars
	}

	if c.DifferentDestinationAccount {
		err := NewClient(c.Destination.Token)
		cobra.CheckErr(err)
	}

	// Create destination workspace
	dest, err := client.Workspaces.Create(ctx, c.Destination.Organization, tfe.WorkspaceCreateOptions{
		Name:             &c.Destination.WorkspaceName,
		TerraformVersion: &source.TerraformVersion,
		WorkingDirectory: &source.WorkingDirectory,
		VCSRepo: &tfe.VCSRepoOptions{
			Branch:       &source.VCSRepo.Branch,
			Identifier:   &source.VCSRepo.Identifier,
			OAuthTokenID: &c.Destination.VCSTokenID,
		},
	})
	cobra.CheckErr(err)

	// Handle Variable Sets
	if !c.DifferentDestinationAccount {
		list, err := client.VariableSets.ListForWorkspace(ctx, source.ID, nil)
		cobra.CheckErr(err)

		for _, set := range list.Items {
			err := client.VariableSets.ApplyToWorkspaces(ctx, set.ID, &tfe.VariableSetApplyToWorkspacesOptions{
				Workspaces: []*tfe.Workspace{dest},
			})
			cobra.CheckErr(err)
		}
	}

	// Handle Variables
	for _, v := range source.Variables {
		_, err := client.Variables.Create(ctx, dest.ID, tfe.VariableCreateOptions{
			Key:         &v.Value,
			Value:       &v.Key,
			Description: &v.Description,
			Category:    &v.Category,
			HCL:         &v.HCL,
			Sensitive:   &v.Sensitive,
		})
		cobra.CheckErr(err)
	}

	if c.DifferentDestinationAccount && c.CopyState {
		c.tfInit()
		return sensitiveVars
	}

	access, err := client.TeamAccess.List(ctx, &tfe.TeamAccessListOptions{
		WorkspaceID: source.ID,
	})
	cobra.CheckErr(err)

	for _, a := range access.Items {
		_, err := client.TeamAccess.Add(ctx, tfe.TeamAccessAddOptions{
			Access:           &a.Access,
			Runs:             &a.Runs,
			Variables:        &a.Variables,
			StateVersions:    &a.StateVersions,
			SentinelMocks:    &a.SentinelMocks,
			WorkspaceLocking: &a.WorkspaceLocking,
			RunTasks:         &a.RunTasks,
			Team:             a.Team,
			Workspace:        dest,
		})
		cobra.CheckErr(err)
	}

	return sensitiveVars
}

func (c cloneConfig) tfInit() {
	const stateFile = ".terraform"

	// Remove previous state file, if it exists
	_, err := os.Stat(stateFile)
	if err == nil {
		cobra.CheckErr(os.RemoveAll(stateFile))
	}

	cobra.CheckErr(os.Setenv(atlasToken, token))

	osCmd := exec.Command("terraform", "init",
		fmt.Sprintf(`-backend-config=name=%s/%s`, c.Source.Organization, c.Source.WorkspaceName))

	cobra.CheckErr(osCmd.Run())

	// Run tf init with new version
	cobra.CheckErr(os.Setenv(atlasToken, c.Destination.Token))
	osCmd = exec.Command("terraform", "init",
		fmt.Sprintf(`-backend-config=name=%s/%s`, c.Destination.Organization, c.Destination.WorkspaceName))

	// Needed to run the command interactively, in order to allow for an automated reply
	cmdStdin, err := osCmd.StdinPipe()
	cobra.CheckErr(err)
	defer cmdStdin.Close()

	cobra.CheckErr(osCmd.Start())

	//  Answer "yes" to the question about creating the new state
	io.Copy(cmdStdin, bytes.NewBufferString("yes\n"))

	cobra.CheckErr(osCmd.Wait())

	cobra.CheckErr(os.Setenv(atlasToken, token))
}

var cloneCfg cloneConfig

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone a Workspace",
	Long:  `Clone a Terraform Cloud Workspace`,
	Args:  cobra.ExactArgs(0),
	PreRun: func(cmd *cobra.Command, args []string) {
		cloneCfg.Source.Organization = organization
	},
	Run: runClone,
}

func init() {
	workspaceCmd.AddCommand(cloneCmd)
	cloneCmd.Flags().StringVarP(&cloneCfg.Destination.Organization, "new-organization", "p", "",
		`Name of the Destination Organization in Terraform Cloud`)
	cloneCmd.Flags().StringVarP(&cloneCfg.Source.WorkspaceName, "source-workspace", "s", "",
		requiredPrefix+`Name of the Source Workspace in Terraform Cloud`)
	cloneCmd.Flags().StringVarP(&cloneCfg.Destination.WorkspaceName, "new-workspace", "n", "",
		requiredPrefix+`Name of the new Workspace in Terraform Cloud`)
	cloneCmd.Flags().StringVarP(&cloneCfg.Destination.VCSTokenID, "new-vcs-token-id", "v", "",
		`The new organization's VCS repo's oauth-token-id`)
	cloneCmd.Flags().BoolVarP(&cloneCfg.CopyState, "copyState", "t", false,
		`optional (e.g. "-t=true") whether to copy the state of the Source Workspace (only possible if copying to a new account).`)
	cloneCmd.Flags().BoolVarP(&cloneCfg.CopyVariables, "copyVariables", "c", false,
		`optional (e.g. "-c=true") whether to copy the values of the Source Workspace variables.`)
	cloneCmd.Flags().BoolVar(&cloneCfg.ApplyVariableSets, "applyVariableSets", false,
		`optional, whether to apply the same variable sets to the new workspace (only for same-account clone).`)
	cloneCmd.Flags().BoolVarP(&cloneCfg.DifferentDestinationAccount, "differentDestinationAccount", "d", false,
		`optional (e.g. "-d=true") whether to clone to a different TF account.`)

	cloneCmd.MarkFlagsRequiredTogether("differentDestinationAccount", "new-vcs-token-id", "new-organization")
	cobra.CheckErr(cloneCmd.MarkFlagRequired("source-workspace"))
	cobra.CheckErr(cloneCmd.MarkFlagRequired("new-workspace"))
}

func runClone(cmd *cobra.Command, args []string) {
	fmt.Printf(`clone '%s > %s' to '%s > %s' using:
  copyState: %t
  copyVariables: %t
  applyVariableSets: %t
  differentDestinationAccount: %t
`,
		cloneCfg.Source.Organization, cloneCfg.Source.WorkspaceName,
		cloneCfg.Destination.Organization, cloneCfg.Destination.WorkspaceName,
		cloneCfg.CopyState, cloneCfg.CopyVariables, cloneCfg.ApplyVariableSets, cloneCfg.DifferentDestinationAccount)

	cloneCfg.Destination.Token = os.Getenv("ATLAS_TOKEN_DESTINATION")
	if cloneCfg.Destination.Token == "" {
		fmt.Print("Info: ATLAS_TOKEN_DESTINATION is not set, using ATLAS_TOKEN for destination account.\n\n")
		cloneCfg.Destination.Token = token
	}

	sensitiveVars := cloneCfg.clone()

	println("\n  **** Completed Cloning ****")
	if len(sensitiveVars) > 0 {
		fmt.Printf("Sensitive variables for %s:%s\n", cloneCfg.Destination.Organization, cloneCfg.Destination.WorkspaceName)
		for _, nextVar := range sensitiveVars {
			fmt.Println(nextVar)
		}
	}
}
