package lib

import (
	"testing"

	"github.com/hashicorp/go-tfe"
)

func Test_WorkspaceListToString(t *testing.T) {
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
			if got := WorkspaceListToString(tt.workspaces); got != tt.want {
				t.Errorf("workspaceListToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
