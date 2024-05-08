package lib

import (
	"fmt"

	"github.com/hashicorp/go-tfe"
)

type RunTriggerConfig struct {
	WorkspaceID       string
	SourceWorkspaceID string
}

func CreateRunTrigger(config RunTriggerConfig) error {
	ws, err := GetWorkspaceByID(config.SourceWorkspaceID)
	if err != nil {
		return err
	}

	_, err = client.RunTriggers.Create(ctx, config.WorkspaceID, tfe.RunTriggerCreateOptions{
		Sourceable: ws,
	})
	return err
}

type FindRunTriggerConfig struct {
	SourceWorkspaceID string
	WorkspaceID       string
}

// FindRunTrigger searches all the run triggers inbound to the given WorkspaceID. If a run trigger is configured for
// the given SourceWorkspaceID, that trigger is returned. Otherwise, nil is returned.
func FindRunTrigger(config FindRunTriggerConfig) (*tfe.RunTrigger, error) {
	triggers, err := ListRunTriggers(ListRunTriggerConfig{
		WorkspaceID: config.WorkspaceID,
		Type:        "inbound",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list run triggers: %w", err)
	}
	for _, t := range triggers {
		if t.Sourceable.ID == config.SourceWorkspaceID {
			return t, nil
		}
	}
	return nil, nil
}

type ListRunTriggerConfig struct {
	WorkspaceID string
	Type        string // must be either "inbound" or "outbound"
}

// ListRunTriggers returns a list of run triggers configured for the given workspace
// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/run-triggers#list-run-triggers
func ListRunTriggers(config ListRunTriggerConfig) ([]*tfe.RunTrigger, error) {
	opts := tfe.RunTriggerListOptions{RunTriggerType: tfe.RunTriggerFilterOp(config.Type)}
	list, err := client.RunTriggers.List(ctx, config.WorkspaceID, &opts)
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}
