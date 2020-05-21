// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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
	cloner "github.com/silinternational/terraform-enterprise-migrator/lib"
	"github.com/spf13/cobra"
)

var copyState bool
var copyVariables bool
var differentDestinationAccount bool
var organization string
var newOrganization string
var sourceWorkspace string
var newWorkspace string
var newVCSTokenID string

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone a V2 Workspace",
	Long: `Clone a TF Enterprise Version 2 Workspace`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if organization == "" {
			fmt.Println("Error: The 'organization' flag is required\n")
			os.Exit(1)
		}
		if sourceWorkspace == "" {
			fmt.Println("Error: The 'source-workspace' flag is required\n")
			os.Exit(1)
		}
		if newWorkspace == "" {
			fmt.Println("Error: The 'new-workspace' flag is required\n")
			os.Exit(1)
		}
		if differentDestinationAccount {

		    if newOrganization == "" {
			    fmt.Println("Error: The 'new-organization' '-p' flag is required for a different destination account.\n")
			    os.Exit(1)
		    }
		    if newVCSTokenID == "" {
			    fmt.Println("Error: The 'new-vcs-token-id' '-v' flag is required for a different destination account.\n")
			    os.Exit(1)
		    }
		}

		config := cloner.V2CloneConfig{
			Organization:                organization,
			NewOrganization:             newOrganization,
			SourceWorkspace:             sourceWorkspace,
			NewWorkspace:                newWorkspace,
			NewVCSTokenID:               newVCSTokenID,
			CopyState:                   copyState,
			CopyVariables:               copyVariables,
			DifferentDestinationAccount: differentDestinationAccount,
		}

		runClone(config)
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)
	cloneCmd.Flags().StringVarP(
		&organization,
		"organization",
		"o",
		"",
		`Name of the Organization in TF Enterprise (version 2)`,
	)
	cloneCmd.Flags().StringVarP(
		&newOrganization,
		"new-organization",
		"p",
		"",
		`Name of the Destination Organization in TF Enterprise (version 2)`,
	)
	cloneCmd.Flags().StringVarP(
		&sourceWorkspace,
			"source-workspace",
			"s",
			"",
			`Name of the Source Workspace in TF Enterprise (version 2)`,
	)
	cloneCmd.Flags().StringVarP(
		&newWorkspace,
		"new-workspace",
		"n",
		"",
		`Name of the new Workspace in TF Enterprise (version 2)`,
	)
	cloneCmd.Flags().StringVarP(
		&newVCSTokenID,
		"new-vcs-token-id",
		"v",
		"",
		`The new organization's VCS repo's oauth-token-id`,
	)
	cloneCmd.Flags().BoolVarP(
		&copyState,
		"copyState",
		"t",
		false,
		`optional (e.g. "-t=true") whether to copy the state of the Source Workspace (only possible if copying to a new account).`,
	)
	cloneCmd.Flags().BoolVarP(
		&copyVariables,
		"copyVariables",
		"c",
		false,
		`optional (e.g. "-c=true") whether to copy the values of the Source Workspace variables.`,
	)
	cloneCmd.Flags().BoolVarP(
		&differentDestinationAccount,
		"differentDestinationAccount",
		"d",
		false,
		`optional (e.g. "-d=true") whether to clone to a different TF account.`,
	)
}

func runClone(cfg cloner.V2CloneConfig) {
	fmt.Printf("clone called using %s, %s, %s, copyState: %t, copyVariables: %t, differentDestinationAccount: %t\n",
		cfg.Organization, cfg.SourceWorkspace, cfg.NewWorkspace, cfg.CopyState, cfg.CopyVariables, cfg.DifferentDestinationAccount)
	cfg.AtlasToken = atlasToken
	cfg.AtlasTokenDestination = atlasTokenDestination

	sensitiveVars, err := cloner.CloneV2Workspace(cfg)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	println("\n  **** Completed Cloning ****")
	if len(sensitiveVars) > 0 {
		fmt.Printf("Sensitive variables for %s:%s\n", cfg.Organization, cfg.NewWorkspace)
		for _, nextVar := range sensitiveVars {
			println(nextVar)
		}
	}

}
