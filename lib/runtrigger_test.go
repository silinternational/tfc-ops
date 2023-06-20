package lib

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_buildRunTriggerPayload(t *testing.T) {
	got := buildRunTriggerPayload("ws_id")
	if got != `{"data":{"relationships":{"sourceable":{"data":{"id":"ws_id","type":"workspaces"}}}}}` {
		t.Fatalf("did not get expected result, got %s", got)
	}
}

func Test_parseRunTriggerListResponse(t *testing.T) {
	r := bytes.NewReader([]byte(listTriggerSampleBody))
	triggers, err := parseRunTriggerListResponse(r)
	require.NoError(t, err)
	require.Equal(t, triggers[0].WorkspaceID, "ws-abcdefghijklmnop")
	require.Equal(t, triggers[0].SourceID, "ws-qrstuvwxyzABCDEF")
	require.Equal(t, triggers[0].WorkspaceName, "a-workspace-name")
	require.Equal(t, triggers[0].SourceName, "source-ws-1")
	require.Equal(t, triggers[0].CreatedAt, time.Date(2023, 6, 20, 8, 56, 50, 996e6, time.UTC))
}

const listTriggerSampleBody = `{
	"data": [
		{
			"id": "rt-abcdefghijklmnop",
			"type": "run-triggers",
			"attributes": {
				"workspace-name": "a-workspace-name",
				"sourceable-name": "source-ws-1",
				"created-at": "2023-06-20T08:56:50.996Z"
			},
			"relationships": {
				"workspace": {
					"data": {
						"id": "ws-abcdefghijklmnop",
						"type": "workspaces"
					}
				},
				"sourceable": {
					"data": {
						"id": "ws-qrstuvwxyzABCDEF",
						"type": "workspaces"
					}
				}
			}
		}
	]
}`
