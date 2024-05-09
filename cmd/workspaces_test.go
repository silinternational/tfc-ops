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
	"testing"

	"github.com/hashicorp/go-tfe"
)

func Test_workspaceListToString(t *testing.T) {
	tests := []struct {
		name       string
		workspaces []*tfe.Workspace
		want       string
	}{
		{
			name:       "nil",
			workspaces: nil,
			want:       "",
		},
		{
			name:       "empty",
			workspaces: []*tfe.Workspace{},
			want:       "",
		},
		{
			name: "one",
			workspaces: []*tfe.Workspace{
				{Name: "one"},
			},
			want: "workspace 'one'",
		},
		{
			name: "two",
			workspaces: []*tfe.Workspace{
				{Name: "one"},
				{Name: "two"},
			},
			want: "workspaces: one, two",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := workspaceListToString(tt.workspaces); got != tt.want {
				t.Errorf("workspaceListToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
