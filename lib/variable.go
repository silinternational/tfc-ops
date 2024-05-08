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

package lib

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-tfe"
)

type UpdateConfig struct {
	Organization          string
	Workspace             string
	SearchString          string //  must be an exact case-insensitive match (i.e. not a partial match)
	NewValue              string
	AddKeyIfNotFound      bool // If true, then SearchOnVariableValue will be treated as false
	SearchOnVariableValue bool // If false, then will filter on variable key
	SensitiveVariable     bool // Whether to mark the variable as sensitive
}

// ConvertHCLVariable changes a TFVar struct in place by escaping
// the double quotes and line endings in the Value attribute
func ConvertHCLVariable(v *tfe.Variable) {
	if !v.HCL {
		return
	}

	v.Value = strings.ReplaceAll(v.Value, `"`, `\"`)
	v.Value = strings.ReplaceAll(v.Value, "\n", "\\n")
}

// CreateVariable makes a Terraform vars API POST to create a variable
// for a given workspace
func CreateVariable(workspaceID string, v *tfe.Variable) error {
	_, err := client.Variables.Create(ctx, workspaceID, tfe.VariableCreateOptions{
		Key:         &v.Key,
		Value:       &v.Value,
		Description: &v.Description,
		Category:    &v.Category,
		HCL:         &v.HCL,
		Sensitive:   &v.Sensitive,
	})
	return err
}

// CreateAllVariables makes several Terraform vars API POSTs to create
// variables for a given workspace
func CreateAllVariables(workspaceID string, vars []*tfe.Variable) error {
	for _, v := range vars {
		if err := CreateVariable(workspaceID, v); err != nil {
			return err
		}
	}
	return nil
}

// DeleteVariable deletes a variable from a workspace
func DeleteVariable(workspaceID, variableID string) error {
	return client.Variables.Delete(ctx, workspaceID, variableID)
}

// GetWorkspaceVar retrieves the variables from a Workspace and returns the variable that matches the given key
func GetWorkspaceVar(organization, workspaceName, key string) (*tfe.Variable, error) {
	vars, err := GetWorkspaceVars(organization, workspaceName)
	if err != nil {
		return nil, err
	}

	for _, v := range vars {
		if v.Key == key {
			return v, nil
		}
	}
	return nil, nil
}

// GetVarsFromWorkspace returns a list of Terraform variables for a given workspace
func GetWorkspaceVars(organization, workspaceName string) ([]*tfe.Variable, error) {
	workspace, err := GetWorkspaceByName(organization, workspaceName)
	return workspace.Variables, err
}

// SearchVarsInAllWorkspaces returns all the variables that match the search terms 'keyContains' and 'valueContains'
// in all workspaces. The return value is a map of variable lists with the workspace name as the key.
func SearchVarsInAllWorkspaces(organization, keyContains, valueContains string) (map[string][]*tfe.Variable, error) {
	workspaces, err := GetAllWorkspaces(organization)
	if err != nil {
		return nil, err
	}

	return SearchVarsInWorkspaces(workspaces, keyContains, valueContains)
}

// SearchVarsInWorkspaces returns all the variables that match the search terms 'keyContains' and 'valueContains'
// in given workspaces. The return value is a map of variable lists with the workspace name as the key.
func SearchVarsInWorkspaces(workspaces []*tfe.Workspace, keyContains, valueContains string) (map[string][]*tfe.Variable, error) {
	result := map[string][]*tfe.Variable{}
	for _, workspace := range workspaces {
		for _, v := range workspace.Variables {
			if variableContains(v, keyContains, valueContains) {
				result[workspace.Name] = append(result[workspace.Name], v)
			}
		}
	}
	return result, nil
}

// SearchVariables returns a list of variables in the given workspace that match the search terms
// 'keyContains' and 'valueContains'
func SearchVariables(organization, workspaceName, keyContains, valueContains string) ([]*tfe.Variable, error) {
	vars, err := GetWorkspaceVars(organization, workspaceName)
	if err != nil {
		return nil, err
	}

	var wsVars []*tfe.Variable
	for _, v := range vars {
		if variableContains(v, keyContains, valueContains) {
			wsVars = append(wsVars, v)
		}
	}
	return wsVars, nil
}

// UpdateVariable makes a Terraform vars API call to update a variable
// for a given workspace
func UpdateVariable(workspaceID string, v *tfe.Variable) error {
	_, err := client.Variables.Update(ctx, workspaceID, v.ID, tfe.VariableUpdateOptions{
		Key:         &v.Key,
		Value:       &v.Value,
		Description: &v.Description,
		Category:    &v.Category,
		HCL:         &v.HCL,
		Sensitive:   &v.Sensitive,
	})
	return err
}

// AddOrUpdateVariable adds or updates an existing Terraform Cloud workspace variable
// If the copyVariables param is set to true, then all the non-sensitive variable values will be added to the new
// workspace.  Otherwise, they will be set to "REPLACE_THIS_VALUE"
func AddOrUpdateVariable(cfg UpdateConfig) (string, error) {
	workspace, err := GetWorkspaceByName(cfg.Organization, cfg.Workspace)
	if err != nil {
		return "", err
	}

	for _, v := range workspace.Variables {
		oldValue := v.Value
		v.Value = cfg.NewValue
		v.Sensitive = cfg.SensitiveVariable
		v.HCL = false

		if cfg.SearchOnVariableValue {
			if !strings.EqualFold(v.Value, cfg.SearchString) {
				continue
			}

			// Found a match
			if !config.readOnly {
				err = UpdateVariable(workspace.ID, v)
			}
			return fmt.Sprintf("Replaced the value of %s from %s to %s", v.Key, oldValue, cfg.NewValue), err
		}

		// Search on variable key, since search on value is not true
		if !strings.EqualFold(v.Key, cfg.SearchString) {
			continue
		}

		// Found a match
		// Only add if there isn't a match
		if cfg.AddKeyIfNotFound {
			return "", fmt.Errorf("addKeyIfNotFound was set to true but a variable already exists with key %s", v.Key)
		}

		if !config.readOnly {
			err = UpdateVariable(workspace.ID, v)
		}
		return fmt.Sprintf("Replaced the value of %s from %s to %s", v.Key, oldValue, cfg.NewValue), err
	}

	// At this point, we haven't found a match
	if cfg.AddKeyIfNotFound {
		if !config.readOnly {
			err = CreateVariable(workspace.ID, &tfe.Variable{
				Key:       cfg.SearchString,
				Value:     cfg.NewValue,
				HCL:       false,
				Sensitive: cfg.SensitiveVariable,
			})
		}
		return fmt.Sprintf("Added variable %s = %s", cfg.SearchString, cfg.NewValue), err
	}

	return "No match found and no variable added", nil
}

func GetVariableSet(org, vsName string) (*tfe.VariableSet, error) {
	sets, err := GetAllVariableSets(org)
	if err != nil {
		return nil, fmt.Errorf("error getting list of variable sets in org: %w", err)
	}
	for _, l := range sets {
		if l.Name == vsName {
			return l, nil
		}
	}
	return nil, nil
}

func GetAllVariableSets(organization string) ([]*tfe.VariableSet, error) {
	list, err := client.VariableSets.List(ctx, organization, nil)
	if err != nil {
		return nil, err
	}
	return list.Items, err
}

func GetWorkspaceVariableSets(workspaceID string) ([]*tfe.VariableSet, error) {
	list, err := client.VariableSets.ListForWorkspace(ctx, workspaceID, nil)
	if err != nil {
		return nil, err
	}
	return list.Items, err
}

func ApplyVariableSet(varsetID string, workspaces []*tfe.Workspace) error {
	return client.VariableSets.ApplyToWorkspaces(ctx, varsetID, &tfe.VariableSetApplyToWorkspacesOptions{Workspaces: workspaces})
}

func ApplyVariableSets(sets []*tfe.VariableSet, workspaces ...*tfe.Workspace) error {
	for _, varset := range sets {
		return ApplyVariableSet(varset.ID, workspaces)
	}
	return nil
}

func copyVariableSets(source, dest *tfe.Workspace) error {
	sets, err := GetWorkspaceVariableSets(source.ID)
	if err != nil {
		return fmt.Errorf("copy variable sets: %w", err)
	}

	if err := ApplyVariableSets(sets, dest); err != nil {
		return fmt.Errorf("copy variable sets: %w", err)
	}
	return nil
}

func variableContains(v *tfe.Variable, key, value string) bool {
	return (key != "" && strings.Contains(v.Key, key)) ||
		(value != "" && strings.Contains(v.Value, value))
}
