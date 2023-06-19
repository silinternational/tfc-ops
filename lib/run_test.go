package lib

import (
	"testing"
)

func Test_buildRunPayload(t *testing.T) {
	got := buildRunPayload("my message", "ws_id")
	if got != `{"data":{"attributes":{"message":"my message"},"relationships":{"workspace":{"data":{"id":"ws_id"}}}}}` {
		t.Fatalf("did not get expected result, got %q", got)
	}
}
