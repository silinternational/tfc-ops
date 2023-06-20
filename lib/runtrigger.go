package lib

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Jeffail/gabs/v2"
)

type RunTriggerConfig struct {
	WorkspaceID       string
	SourceWorkspaceID string
}

func CreateRunTrigger(config RunTriggerConfig) error {
	u := NewTfcUrl("/workspaces/" + config.WorkspaceID + "/run-triggers")
	payload := buildRunTriggerPayload(config.SourceWorkspaceID)
	_ = callAPI(http.MethodPost, u.String(), payload, nil)
	return nil
}

func buildRunTriggerPayload(sourceWorkspaceID string) string {
	data := gabs.New()
	_, err := data.Object("data")
	if err != nil {
		return "unable to create run trigger payload:" + err.Error()
	}

	workspaceObject := gabs.Wrap(map[string]any{
		"type": "workspaces",
		"id":   sourceWorkspaceID,
	})
	if _, err = data.SetP(workspaceObject, "data.relationships.sourceable.data"); err != nil {
		return "unable to complete run trigger payload:" + err.Error()
	}

	return data.String()
}

type FindRunTriggerConfig struct {
	SourceID    string
	WorkspaceID string
}

// FindRunTrigger searches all the run triggers inbound to the given WorkspaceID. If a run trigger is configured for
// the given SourceID, that trigger is returned. Otherwise, nil is returned.
func FindRunTrigger(config FindRunTriggerConfig) (*RunTrigger, error) {
	triggers, err := ListRunTriggers(ListRunTriggerConfig{
		WorkspaceID: config.WorkspaceID,
		Type:        "inbound",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list run triggers: %w", err)
	}
	for _, t := range triggers {
		if t.SourceID == config.SourceID {
			return &t, nil
		}
	}
	return nil, nil
}

type RunTrigger struct {
	CreatedAt     time.Time
	SourceName    string
	SourceID      string
	WorkspaceName string
	WorkspaceID   string
}

type ListRunTriggerConfig struct {
	WorkspaceID string
	Type        string // must be either "inbound" or "outbound"
}

// ListRunTriggers returns a list of run triggers configured for the given workspace
// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/run-triggers#list-run-triggers
func ListRunTriggers(config ListRunTriggerConfig) ([]RunTrigger, error) {
	u := NewTfcUrl("/workspaces/" + config.WorkspaceID + "/run-triggers?filter%5Brun-trigger%5D%5Btype%5D=" + config.Type)
	resp := callAPI(http.MethodGet, u.String(), "", nil)
	triggers, err := parseRunTriggerListResponse(resp.Body)
	if err != nil {
		return nil, err
	}
	return triggers, nil
}

func parseRunTriggerListResponse(r io.Reader) ([]RunTrigger, error) {
	parsed, err := gabs.ParseJSONBuffer(r)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response data: %w", err)
	}

	attributes := parsed.Search("data", "*").Children()
	triggers := make([]RunTrigger, len(attributes))
	for i, attr := range attributes {
		trigger := RunTrigger{
			SourceID:      attr.Path("relationships.sourceable.data.id").Data().(string),
			SourceName:    attr.Path("attributes.sourceable-name").Data().(string),
			WorkspaceID:   attr.Path("relationships.workspace.data.id").Data().(string),
			WorkspaceName: attr.Path("attributes.workspace-name").Data().(string),
		}
		createdAt, _ := time.Parse(time.RFC3339, attr.Path("attributes.created-at").Data().(string))
		trigger.CreatedAt = createdAt
		triggers[i] = trigger
	}
	return triggers, nil
}
