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
	"context"

	"github.com/hashicorp/go-tfe"
)

var (
	client *tfe.Client
	ctx    = context.Background()
)

func NewClient(t string) error {
	cfg := tfe.DefaultConfig()
	if t == "" {
		t = token
	}
	cfg.Token = t

	var err error
	client, err = tfe.NewClient(cfg)
	return err
}
