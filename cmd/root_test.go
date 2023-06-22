package cmd

import (
	"testing"
)

func Test_stringMapToSlice(t *testing.T) {
	tests := []struct {
		name       string
		input      map[string]string
		wantKeys   []string
		wantValues []string
	}{
		{
			name:       "nil",
			input:      nil,
			wantKeys:   []string{},
			wantValues: []string{},
		},
		{
			name:       "one",
			input:      map[string]string{"1": "one"},
			wantKeys:   []string{"1"},
			wantValues: []string{"one"},
		},
		{
			name:       "two",
			input:      map[string]string{"1": "one", "2": "two"},
			wantKeys:   []string{"1", "2"},
			wantValues: []string{"one", "two"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keys, values := stringMapToSlice(tt.input)
			for i := range keys {
				if tt.input[keys[i]] != values[i] {
					t.Errorf("stringMapToSlice() got = %v, want %v", values[i], tt.input[keys[i]])
				}
			}
		})
	}
}

func Test_workspaceListToString(t *testing.T) {
	tests := []struct {
		name           string
		workspaceNames []string
		suppressHeader bool
		want           string
	}{
		{
			name:           "nil",
			workspaceNames: nil,
			suppressHeader: false,
			want:           "",
		},
		{
			name:           "empty",
			workspaceNames: []string{},
			suppressHeader: false,
			want:           "",
		},
		{
			name:           "one (header)",
			workspaceNames: []string{"one"},
			suppressHeader: false,
			want:           "workspace 'one'",
		},
		{
			name:           "one (no header)",
			workspaceNames: []string{"one"},
			suppressHeader: true,
			want:           "one",
		},
		{
			name:           "two (header)",
			workspaceNames: []string{"one", "two"},
			suppressHeader: false,
			want:           "workspaces: one, two",
		},
		{
			name:           "two (no header)",
			workspaceNames: []string{"one", "two"},
			suppressHeader: true,
			want:           "one, two",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suppressCSVHeader = tt.suppressHeader
			if got := workspaceListToString(tt.workspaceNames); got != tt.want {
				t.Errorf("workspaceListToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
