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

	"github.com/hashicorp/go-tfe"
)

// GetTeamAccessFrom returns the team access data from an existing workspace
func GetTeamAccessFrom(workspaceID string) ([]*tfe.TeamAccess, error) {
	opts := tfe.TeamAccessListOptions{WorkspaceID: workspaceID}
	list, err := client.TeamAccess.List(ctx, &opts)
	if err != nil {
		return nil, fmt.Errorf("error getting team workspace data for %s\n%w", workspaceID, err)
	}
	return list.Items, nil
}

// AssignTeamAccess assigns the requested team access to a workspace on Terraform Cloud
func AssignTeamAccess(workspace *tfe.Workspace, access ...*tfe.TeamAccess) error {
	for _, a := range access {
		_, err := client.TeamAccess.Add(ctx, tfe.TeamAccessAddOptions{
			Access:           &a.Access,
			Runs:             &a.Runs,
			Variables:        &a.Variables,
			StateVersions:    &a.StateVersions,
			SentinelMocks:    &a.SentinelMocks,
			WorkspaceLocking: &a.WorkspaceLocking,
			RunTasks:         &a.RunTasks,
			Team:             a.Team,
			Workspace:        workspace,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
