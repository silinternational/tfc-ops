// Copyright © 2018-2022 SIL International
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lib

import (
	"strings"

	"github.com/hashicorp/go-tfe"
)

func WorkspaceListToString(workspaces []*tfe.Workspace) string {
	if len(workspaces) == 0 {
		return ""
	}

	names := make([]string, len(workspaces))
	for i, w := range workspaces {
		names[i] = w.Name
	}

	s := ""
	if len(workspaces) > 1 {
		s = "workspaces: " + strings.Join(names, ", ")
	} else {
		s = "workspace '" + names[0] + "'"
	}

	return s
}
