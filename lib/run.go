package lib

import (
	"github.com/hashicorp/go-tfe"
)

type RunConfig struct {
	Message     string
	WorkspaceID string
}

// CreateRun creates a Run, which starts a Plan, which can later be Applied.
// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/run
func CreateRun(config RunConfig) error {
	ws, err := GetWorkspaceByID(config.WorkspaceID)
	if err != nil {
		return err
	}

	_, err = client.Runs.Create(ctx, tfe.RunCreateOptions{
		Message:   &config.Message,
		Workspace: ws,
	})
	return err
}
