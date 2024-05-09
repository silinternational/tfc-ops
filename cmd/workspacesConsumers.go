// Copyright © 2018-2024 SIL International
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

package cmd

import (
	"github.com/spf13/cobra"
)

// workspaceConsumersCmd represents the top level command for workspace consumers
var workspaceConsumersCmd = &cobra.Command{
	Use:   "consumers",
	Short: "Manage workspace remote state consumers",
	Long:  `Add, list, update, or delete workspace remote state consumers`,
	Args:  cobra.MinimumNArgs(1),
}

func init() {
	workspaceCmd.AddCommand(workspaceConsumersCmd)
}
