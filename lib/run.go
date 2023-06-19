package lib

import (
	"net/http"

	"github.com/Jeffail/gabs/v2"
)

type RunConfig struct {
	Message     string
	WorkspaceID string
}

// CreateRun creates a Run, which starts a Plan, which can later be Applied.
// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/run
func CreateRun(config RunConfig) error {
	u := NewTfcUrl("/runs")
	payload := buildRunPayload(config.Message, config.WorkspaceID)
	_ = callAPI(http.MethodPost, u.String(), payload, nil)
	return nil
}

func buildRunPayload(message, workspaceID string) string {
	data := gabs.New()

	_, err := data.Object("data")
	if err != nil {
		return "error"
	}

	if _, err = data.SetP(message, "data.attributes.message"); err != nil {
		return "unable to process attribute for update:" + err.Error()
	}

	if _, err = data.SetP(workspaceID, "data.relationships.workspace.data.id"); err != nil {
		return "unable to process attribute for update:" + err.Error()
	}

	return data.String()
}
