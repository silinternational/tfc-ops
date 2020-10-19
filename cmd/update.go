package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	updater "github.com/silinternational/terraform-enterprise-migrator/lib"
)

var workspace string
var variableSearchString string
var newVariableValue string
var searchOnVariableValue bool
var addKeyIfNotFound bool
var dryRunMode bool

// cloneCmd represents the clone command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update/add a variable in a V2 Workspace",
	Long:  `Update or add a variable in a TF Enterprise Version 2 Workspace based on a complete case-insensitive match`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if organization == "" {
			fmt.Println("Error: The 'organization' flag is required\n")
			os.Exit(1)
		}
		if workspace == "" {
			fmt.Println("Error: The 'workspace' flag is required\n")
			os.Exit(1)
		}

		if variableSearchString == "" {
			fmt.Println("Error: The 'variable-search-string' flag is required\n")
			os.Exit(1)
		}
		if newVariableValue == "" {
			fmt.Println("Error: The 'new-variable-value' flag is required\n")
			os.Exit(1)
		}
		if addKeyIfNotFound && searchOnVariableValue {
			fmt.Println("Error: The 'add-key-if-not-found' flag may not be used with the 'search-on-variable-value' flag\n")
			os.Exit(1)
		}
		config := updater.V2UpdateConfig{
			Organization:          organization,
			NewOrganization:       newOrganization,
			Workspace:             workspace,
			SearchString:          variableSearchString,
			NewValue:              newVariableValue,
			SearchOnVariableValue: searchOnVariableValue,
			AddKeyIfNotFound:      addKeyIfNotFound,
			DryRunMode:            dryRunMode,
		}

		runUpdateVariable(config)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringVarP(
		&organization,
		"organization",
		"o",
		"",
		`Name of the Organization in TF Enterprise (version 2)`,
	)
	updateCmd.Flags().StringVarP(
		&workspace,
		"workspace",
		"w",
		"",
		`Name of the Workspace in TF Enterprise (version 2)`,
	)
	updateCmd.Flags().StringVarP(
		&variableSearchString,
		"variable-search-string",
		"s",
		"",
		`The string to match in the current variables (either in the Key or Value - see other flags)`,
	)
	updateCmd.Flags().StringVarP(
		&newVariableValue,
		"new-variable-value",
		"n",
		"",
		`The desired new value of the variable`,
	)
	updateCmd.Flags().BoolVarP(
		&addKeyIfNotFound,
		"add-key-if-not-found",
		"a",
		false,
		`optional (e.g. "-a=true") whether to add a new variable if a matching key is not found.`,
	)
	updateCmd.Flags().BoolVarP(
		&searchOnVariableValue,
		"search-on-variable-value",
		"v",
		false,
		`optional (e.g. "-v=true") whether to do the search based on the value of the variables. (Must be false if add-key-if-not-found is true`,
	)
	updateCmd.Flags().BoolVarP(
		&dryRunMode,
		"dry-run-mode",
		"d",
		false,
		`optional (e.g. "-d=true") dry run mode only.`,
	)
}

func runUpdateVariable(cfg updater.V2UpdateConfig) {
	if cfg.AddKeyIfNotFound {
		if cfg.SearchOnVariableValue {
			println("update variable aborted. Because addKeyIfNotFound was true, searchOnVariableValue must be set to false")
			return
		}
		cfg.SearchOnVariableValue = false
	}

	if cfg.DryRunMode {
		println("\n ****  DRY RUN MODE  ****")
	}
	fmt.Printf("update variable called using %s, %s, search string: %s, new value: %s, add-key-if-not-found: %t, search-on-variable-value: %t\n",
		cfg.Organization, cfg.Workspace, cfg.SearchString, cfg.NewValue, cfg.AddKeyIfNotFound, cfg.SearchOnVariableValue)
	cfg.AtlasToken = atlasToken

	message, err := updater.AddOrUpdateV2Variable(cfg)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	println("\n  **** Completed Updating ****")
	println(message)

}
