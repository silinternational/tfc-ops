package cmd

import (
	"reflect"
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
			got, got1 := stringMapToSlice(tt.input)
			if !reflect.DeepEqual(got, tt.wantKeys) {
				t.Errorf("stringMapToSlice() got = %v, want %v", got, tt.wantKeys)
			}
			if !reflect.DeepEqual(got1, tt.wantValues) {
				t.Errorf("stringMapToSlice() got1 = %v, want %v", got1, tt.wantValues)
			}
		})
	}
}

func Test_workspaceListToString(t *testing.T) {
	tests := []struct {
		name           string
		workspaceNames []string
		want           string
	}{
		{
			name:           "nil",
			workspaceNames: nil,
			want:           "",
		},
		{
			name:           "empty",
			workspaceNames: []string{},
			want:           "",
		},
		{
			name:           "one",
			workspaceNames: []string{"one"},
			want:           "workspace 'one'",
		},
		{
			name:           "two",
			workspaceNames: []string{"one", "two"},
			want:           "workspaces: one, two",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := workspaceListToString(tt.workspaceNames); got != tt.want {
				t.Errorf("workspaceListToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
