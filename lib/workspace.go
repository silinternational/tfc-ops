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
	"reflect"

	"github.com/hashicorp/go-tfe"
)

type CloneConfig struct {
	Organization                string
	NewOrganization             string
	SourceWorkspace             string
	NewWorkspace                string
	NewVCSTokenID               string
	AtlasToken                  string
	AtlasTokenDestination       string
	CopyState                   bool
	CopyVariables               bool
	ApplyVariableSets           bool
	DifferentDestinationAccount bool
}

type WorkspaceUpdateParams struct {
	Organization    string
	WorkspaceFilter string
	Attribute       string
	Value           string
}

// GetAllWorkspaces retrieves all workspaces from Terraform Cloud and returns a list of Workspace objects
func GetAllWorkspaces(organization string) ([]*tfe.Workspace, error) {
	list, err := client.Workspaces.List(ctx, organization, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting workspace data for %s: %w", organization, err)
	}

	return list.Items, nil
}

func GetWorkspaceByName(organization, workspaceName string) (*tfe.Workspace, error) {
	return client.Workspaces.Read(ctx, organization, workspaceName)
}

func GetWorkspaceByID(workspaceID string) (*tfe.Workspace, error) {
	return client.Workspaces.ReadByID(ctx, workspaceID)
}

// CreateWorkspace makes a Terraform workspaces API call to create a workspace for a given organization, including
// setting up its VCS repo integration. Returns the properties of the new workspace.
func CreateWorkspace(oc OpsConfig, vcsTokenID string) (*tfe.Workspace, error) {
	opts := tfe.WorkspaceCreateOptions{
		Name:             &oc.NewName,
		TerraformVersion: &oc.TerraformVersion,
		WorkingDirectory: &oc.Directory,
	}
	if vcsTokenID != "" {
		opts.VCSRepo = &tfe.VCSRepoOptions{
			Branch:       &oc.Branch,
			Identifier:   &oc.RepoID,
			OAuthTokenID: &vcsTokenID,
		}
	}
	return client.Workspaces.Create(ctx, oc.NewOrg, opts)
}

// CloneWorkspace gets the data, variables and team access data for an existing Terraform Cloud workspace
// and then creates a clone of it with the same data.
//
// If the copyVariables param is set to true, then all the non-sensitive variable values will be added to the new
// workspace.  Otherwise, they will be set to "REPLACE_THIS_VALUE"
func CloneWorkspace(cfg CloneConfig) (*tfe.Workspace, []string, error) {
	source, err := GetWorkspaceByName(cfg.Organization, cfg.SourceWorkspace)
	if err != nil {
		return nil, nil, err
	}

	vars, err := GetWorkspaceVars(cfg.Organization, cfg.SourceWorkspace)
	if err != nil {
		return nil, nil, err
	}

	if !cfg.DifferentDestinationAccount {
		cfg.NewOrganization = cfg.Organization
		cfg.NewVCSTokenID = source.VCSRepo.Identifier
	}

	oc := OpsConfig{
		SourceOrg:        cfg.Organization,
		SourceName:       source.Name,
		NewOrg:           cfg.NewOrganization,
		NewName:          cfg.NewWorkspace,
		TerraformVersion: source.TerraformVersion,
		RepoID:           source.VCSRepo.Identifier,
		Branch:           source.VCSRepo.Branch,
		Directory:        source.WorkingDirectory,
	}

	sensitiveVars := []string{}
	const sensitiveValue = "TF_ENTERPRISE_SENSITIVE_VAR"
	const defaultValue = "REPLACE_THIS_VALUE"

	for _, v := range vars {
		if !cfg.CopyVariables {
			v.Value = defaultValue
		}

		if v.Value == sensitiveValue {
			sensitiveVars = append(sensitiveVars, v.Key)
		}
	}

	if config.readOnly {
		return nil, sensitiveVars, nil
	}

	if cfg.DifferentDestinationAccount {
		if err := NewClient(cfg.AtlasTokenDestination); err != nil {
			return nil, nil, err
		}

		workspace, err := CreateWorkspace(oc, cfg.NewVCSTokenID)
		if err != nil {
			return nil, nil, err
		}

		if err := CreateAllVariables(workspace.ID, vars); err != nil {
			return workspace, nil, err
		}

		if cfg.CopyState {
			if err := RunTFInit(oc, cfg.AtlasToken, cfg.AtlasTokenDestination); err != nil {
				return workspace, sensitiveVars, err
			}
		}

		return workspace, sensitiveVars, nil
	}

	dest, err := CreateWorkspace(oc, source.VCSRepo.OAuthTokenID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create new workspace: %w", err)
	}

	if err := copyVariableSets(source, dest); err != nil {
		return dest, nil, fmt.Errorf("failed to clone variable sets: %w", err)
	}

	if err := CreateAllVariables(dest.ID, vars); err != nil {
		return dest, nil, err
	}

	// Get Team Access Data for source Workspace
	access, err := GetTeamAccessFrom(source.ID)
	if err != nil {
		return dest, sensitiveVars, err
	}

	return dest, sensitiveVars, AssignTeamAccess(dest, access...)
}

// UpdateWorkspace updates one attribute of one or more Terraform Cloud workspaces.
func UpdateWorkspace(params WorkspaceUpdateParams) error {
	if err := validateUpdateWorkspaceParams(params); err != nil {
		return err
	}
	var opts tfe.WorkspaceUpdateOptions
	r := reflect.ValueOf(opts)
	f := reflect.Indirect(r).FieldByName(params.Attribute)
	f.Set(reflect.ValueOf(params.Value))

	foundWs, err := FindWorkspaces(params.Organization, params.WorkspaceFilter)
	if err != nil {
		return err
	}
	if len(foundWs) == 0 {
		return fmt.Errorf("no workspaces found matching the filter '%s'", params.WorkspaceFilter)
	}

	for _, w := range foundWs {
		fmt.Printf("set '%s' to '%s' on workspace %s\n", params.Attribute, params.Value, w.Name)
		if _, err := client.Workspaces.Update(ctx, params.Organization, w.Name, opts); err != nil {
			return err
		}
	}

	fmt.Printf("Updated %d workspace(s)\n", len(foundWs))
	return nil
}

func validateUpdateWorkspaceParams(params WorkspaceUpdateParams) error {
	if config.debug {
		fmt.Printf("params:\n    %#v\n", params)
	}

	if len(params.WorkspaceFilter) < 3 {
		return fmt.Errorf("workspace filter must be at least 3 characters, given: '%s'", params.WorkspaceFilter)
	}

	return nil
}

// FindWorkspaces uses the `search[name]` parameter to retrieve a list of workspaces in Terraform Cloud that
// match the workspaceFilter by the workspace name. The list is returned as a map with the ID in the key
// and the name in the value.
func FindWorkspaces(organization, workspaceFilter string) ([]*tfe.Workspace, error) {
	list, err := client.Workspaces.List(ctx, organization, &tfe.WorkspaceListOptions{Search: workspaceFilter})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

// GetWorkspaceAttributes returns a list of all workspaces in `organization` and the values of the attributes requested
// in the `attributes` list. The value of unrecognized attribute names will be returned as `null`.
func GetWorkspaceAttributes(organization string, attributes []string) ([][]string, error) {
	workspaces, err := GetAllWorkspaces(organization)
	if err != nil {
		return nil, err
	}

	attributeData := make([][]string, len(workspaces))
	for i, w := range workspaces {
		r := reflect.ValueOf(w)
		attributeData[i] = make([]string, len(attributes))
		for j, a := range attributes {
			attributeData[i][j] = reflect.Indirect(r).FieldByName(a).String()
		}
	}
	return attributeData, nil
}

func AddRemoteStateConsumers(workspaceID string, consumerIDs []string) error {
	var err error
	opts := tfe.WorkspaceAddRemoteStateConsumersOptions{
		Workspaces: make([]*tfe.Workspace, len(consumerIDs)),
	}
	for i, id := range consumerIDs {
		opts.Workspaces[i], err = client.Workspaces.ReadByID(ctx, id)
		if err != nil {
			return err
		}
	}

	return client.Workspaces.AddRemoteStateConsumers(ctx, workspaceID, opts)
}
