package lib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewTfcUrl(t *testing.T) {
	orgs := NewTfcUrl("/organizations")
	require.Equal(t, baseURL+"/organizations", orgs.String())

	withQuery := NewTfcUrl("/organizations?q=foo")
	require.Equal(t, baseURL+"/organizations", withQuery.String())
}
