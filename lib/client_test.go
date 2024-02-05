package lib

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetCreateWorkspacePayload(t *testing.T) {
	assert := require.New(t)
	tests := []struct {
		name       string
		oc         OpsConfig
		vcsTokenID string
		want       string
	}{
		{
			name: "with VCS",
			oc: OpsConfig{
				NewName:          "new-name",
				TerraformVersion: "version",
				RepoID:           "repo-id",
				Branch:           "branch",
				Directory:        "directory",
			},
			vcsTokenID: "token-id",
			want: `{
  "data": {
    "attributes": {
      "name": "new-name",
      "terraform_version": "version",
      "vcs-repo": {
        "branch": "branch",
        "default-branch": true,
        "identifier": "repo-id",
        "oauth-token-id": "token-id"
      },
      "working-directory": "directory"
    },
    "type": "workspaces"
  }
} `,
		},
		{
			name: "without VCS",
			oc: OpsConfig{
				NewName:          "new-name",
				TerraformVersion: "version",
				RepoID:           "repo-id",
				Branch:           "branch",
				Directory:        "directory",
			},
			vcsTokenID: "",
			want: `{
  "data": {
    "attributes": {
      "name": "new-name",
      "terraform_version": "version",
      "working-directory": "directory"
    },
    "type": "workspaces"
  }
} `,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetCreateWorkspacePayload(tt.oc, tt.vcsTokenID)
			assert.Equal(removeWhitespace(tt.want), removeWhitespace(got))
		})
	}
}

func removeWhitespace(s string) string {
	s1 := strings.ReplaceAll(s, " ", "")
	return strings.ReplaceAll(s1, "\n", "")
}
