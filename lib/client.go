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
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/hashicorp/go-tfe"
)

var ctx = context.Background()
var client *tfe.Client
var config struct {
	token    string
	debug    bool
	readOnly bool
}

func EnableDebug() {
	config.debug = true
}

func EnableReadOnlyMode() {
	config.readOnly = true
}

func SetToken(t string) {
	config.token = t
}

func NewClient(token string) error {
	cfg := tfe.DefaultConfig()
	if token != "" {
		token = config.token
	}
	cfg.Token = token

	var err error
	client, err = tfe.NewClient(cfg)
	return err
}

// RunTFInit ...
//   - removes old terraform.tfstate files
//   - runs terraform init with old versions
//   - runs terraform init with new version
//
// NOTE: This procedure can be used to copy/migrate a workspace's state to a new one.
// (see the -backend-config mention below and the backend.tf file in this repo)
func RunTFInit(oc OpsConfig, tfToken, tfTokenDestination string) error {
	var tfInit string
	var err error
	var osCmd *exec.Cmd
	var stderr bytes.Buffer

	tokenEnv := "ATLAS_TOKEN"

	stateFile := ".terraform"

	// Remove previous state file, if it exists
	_, err = os.Stat(stateFile)
	if err == nil {
		err = os.RemoveAll(stateFile)
		if err != nil {
			return err
		}
	}

	if err := os.Setenv(tokenEnv, tfToken); err != nil {
		return fmt.Errorf("error setting %s environment variable to source value: %s", tokenEnv, err)
	}

	tfInit = fmt.Sprintf(`-backend-config=name=%s/%s`, oc.SourceOrg, oc.SourceName)

	osCmd = exec.Command("terraform", "init", tfInit)
	osCmd.Stderr = &stderr

	err = osCmd.Run()
	if err != nil {
		println("Error with Legacy: " + tfInit)
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return err
	}

	if err := os.Setenv(tokenEnv, tfTokenDestination); err != nil {
		return fmt.Errorf("error setting %s environment variable to destination value: %s", tokenEnv, err)
	}

	// Run tf init with new version
	tfInit = fmt.Sprintf(`-backend-config=name=%s/%s`, oc.NewOrg, oc.NewName)
	osCmd = exec.Command("terraform", "init", tfInit)
	osCmd.Stderr = &stderr

	// Needed to run the command interactively, in order to allow for an automated reply
	cmdStdin, err := osCmd.StdinPipe()
	if err != nil {
		println("Error with StdinPipe: " + tfInit)
		return err
	}

	err = osCmd.Start()
	if err != nil {
		return err
	}

	defer cmdStdin.Close()
	io.Copy(cmdStdin, bytes.NewBufferString("yes\n"))

	//  Answer "yes" to the question about creating the new state
	err = osCmd.Wait()
	if err != nil {
		println("Error waiting for new tf init: " + tfInit)
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return err
	}

	if err := os.Setenv(tokenEnv, tfToken); err != nil {
		return fmt.Errorf("error resetting %s environment variable back to source value: %s", tokenEnv, err)
	}

	return nil
}
