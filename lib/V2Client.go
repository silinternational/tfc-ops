package lib

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	//"io/ioutil"
	"time"
	"encoding/json"
	"text/tabwriter"
)

// V2Workspace holds the information needed to create a v2 terraform workspace
type V2Workspace struct {
	Name         string
	TFVersion    string
	VCSRepoID    string
	VCSBranch    string
	TFWorkingDir string
}

// ConvertHCLVariable changes a TFVar struct in place by escaping
//  the double quotes and line endings in the Value attribute
func ConvertHCLVariable(tfVar *TFVar) {
	if !tfVar.Hcl {
		return
	}

	tfVar.Value = strings.Replace(tfVar.Value, `"`, `\"`, -1)
	tfVar.Value = strings.Replace(tfVar.Value, "\n", "\\n", -1)
}

// GetCreateV2VariablePayload returns the json needed to make a Post to the
// v2 terraform vars api
func GetCreateV2VariablePayload(organization, workspaceName string, tfVar TFVar) string {
	return fmt.Sprintf(`
{
  "data": {
    "type":"vars",
    "attributes": {
      "key":"%s",
      "value":"%s",
      "category":"terraform",
      "hcl":%t,
      "sensitive":false
    }
  },
  "filter": {
    "organization": {
      "name":"%s"
    },
    "workspace": {
      "name":"%s"
    }
  }
}
`, tfVar.Key, tfVar.Value, tfVar.Hcl, organization, workspaceName)
}

// CreateV2Variable makes a v2 terraform vars api post to create a variable
// for a given organization and v2 workspace
func CreateV2Variable(organization, workspaceName, tfToken string, tfVar TFVar) {
	url := "https://app.terraform.io/api/v2/vars"

	ConvertHCLVariable(&tfVar)

	postData := GetCreateV2VariablePayload(organization, workspaceName, tfVar)

	headers := map[string]string{
		"Authorization": "Bearer " + tfToken,
		"Content-Type":  "application/vnd.api+json",
	}
	resp := CallAPI("POST", url, postData, headers)

	defer resp.Body.Close()
	// bodyBytes, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(bodyBytes))
	return
}

// CreateAllV2Variables makes several v2 terraform vars api posts to create
// variables for a given organization and v2 workspace
func CreateAllV2Variables(organization, workspaceName, tfToken string, tfVars []TFVar) {
	for _, nextVar := range tfVars {
		CreateV2Variable(organization, workspaceName, tfToken, nextVar)
	}
}

// GetCreateV2WorkspacePayload returns the json needed to make a Post to the
// v2 terraform workspaces api
func GetCreateV2WorkspacePayload(mp MigrationPlan, vcsTokenID string) string {
	return fmt.Sprintf(`
{
  "data": {
    "attributes": {
      "name": "%s",
      "terraform_version": "%s",
      "working-directory": "%s",
      "vcs-repo": {
        "identifier": "%s",
        "oauth-token-id": "%s",
        "branch": "%s",
        "default-branch": true
      }
    },
    "type": "workspaces"
  }

}
  `, mp.NewName, mp.TerraformVersion, mp.Directory, mp.RepoID, vcsTokenID, mp.Branch)
}

// CreateV2Workspace makes a v2 terraform workspaces api Post to create a
// workspace for a given organization, including setting up its vcs repo integration
func CreateV2Workspace(
	mp MigrationPlan,
	tfToken, vcsTokenID string,
) {
	url := fmt.Sprintf(
		"https://app.terraform.io/api/v2/organizations/%s/workspaces",
		mp.NewOrg,
	)

	postData := GetCreateV2WorkspacePayload(mp, vcsTokenID)

	headers := map[string]string{
		"Authorization": "Bearer " + tfToken,
		"Content-Type":  "application/vnd.api+json",
	}
	resp := CallAPI("POST", url, postData, headers)

	defer resp.Body.Close()
	// bodyBytes, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(bodyBytes))
}

// CreateAndPopulateV2Workspace makes several api calls to get the variable values
// for a v1 terraform environment and create a corresponding v2 terraform
// workspace along with those same variables.
// Note that the values for sensitive v1 variables will need to be corrected
// in the v2 workspace.
func CreateAndPopulateV2Workspace(
	mp MigrationPlan,
	tfToken, vcsTokenID string,
) ([]string, error) {

	v1Vars, err := GetTFVarsFromV1Config(mp.LegacyOrg, mp.LegacyName, tfToken)
	if err != nil {
		return []string{}, err
	}

	CreateV2Workspace(mp, tfToken, vcsTokenID)

	CreateAllV2Variables(mp.NewOrg, mp.NewName, tfToken, v1Vars)
	sensitiveVars := []string{}
	sensitiveValue := "TF_ENTERPRISE_SENSITIVE_VAR"

	for _, nextVar := range v1Vars {
		if nextVar.Value == sensitiveValue {
			sensitiveVars = append(sensitiveVars, nextVar.Key)
		}
	}

	return sensitiveVars, nil
}

// RunTFInit ...
//  - removes old terraform.tfstate files
//  - runs terraform init with old versions
//  - runs terraform init with new version
func RunTFInit(mp MigrationPlan, tfToken string) error {
	var tfInit string
	var err error
	var osCmd *exec.Cmd
	var stderr bytes.Buffer

	stateFile := ".terraform"

	// Remove previous state file, if it exists
	_, err = os.Stat(stateFile)
	if err == nil {
		err = os.RemoveAll(stateFile)
		if err != nil {
			return err
		}
	}

	tfInit = fmt.Sprintf(`-backend-config=name=%s/%s`, mp.LegacyOrg, mp.LegacyName)

	osCmd = exec.Command("terraform", "init", tfInit)
	osCmd.Stderr = &stderr

	err = osCmd.Run()
	if err != nil {
		println("Error with Legacy: " + tfInit)
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return err
	}

	// Run tf init with new version
	tfInit = fmt.Sprintf(`-backend-config=name=%s/%s`, mp.NewOrg, mp.NewName)
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

	return nil
}

// CreateAndPopulateAllV2Workspaces makes several api calls to retrieve variables
// from v1 environments and create and populate corresponding v2 workspaces.
//   It relies on a csv file with the following columns (the first row will be ignored).
//   Note: these columns are defined by the MigrationPlan struct
// - the name of the organization in Legacy
// - the name of the legacy environment
// - the of the organization in the new Enterprise
// - the name of the new workspace
// - the new terraform version (e.g. ""0.11.3")
// - the id of the version control repo (e.g. "myorg/myproject")
// - the version control branch that the new workspace should be linked to
// - the directory that holds the terraform configuration files in the vcs repo
//
// It also runs `terraform init` with the legacy information and then again with
// the new version.
func CreateAndPopulateAllV2Workspaces(configFile, tfToken, vcsUsername string) (map[string][]string, error) {
	var sensitiveVars []string
	var allPlans []MigrationPlan

	completed := map[string][]string{}

	// Get config contents
	csvFile, err := os.Open(configFile)
	if err != nil {
		return completed, err
	}

	reader := csv.NewReader(bufio.NewReader(csvFile))

	// Throw away the first line with its column headers
	_, err = reader.Read()
	if err != nil {
		newErr := fmt.Errorf("Error reading first row of the plan ...\n %s", err.Error())
		return completed, newErr
	}

	vcsTokenIDs := map[string]string{}

	for rowNum := 2; ; rowNum++ {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			newErr := fmt.Errorf("Error reading row %d of the plan ...\n %s", rowNum, err.Error())
			return completed, newErr
		}

		migrationPlan, err := NewMigrationPlan(line)
		if err != nil {
			fmt.Printf("Skipping row %d because of an error ...\n %s\n", rowNum, err.Error())
			continue
		}
		allPlans = append(allPlans, migrationPlan)
	}

	println("\n\n *** The following environments will be migrated\n")

	w := tabwriter.NewWriter(os.Stdout, 0, 3, 3, ' ', 0)
	tabbedNames := "Legacy Org/Legacy Environment\tNew Org/New Workspace"
	fmt.Fprintln(w, tabbedNames)
	fmt.Fprintln(w, "-----------------------------\t---------------------")

	for _, migrationPlan := range allPlans {
		nextMigration := fmt.Sprintf(
			"%s/%s\t%s/%s",
			migrationPlan.LegacyOrg,
			migrationPlan.LegacyName,
			migrationPlan.NewOrg,
			migrationPlan.NewName,
		)
		fmt.Fprintln(w, nextMigration)
	}
	w.Flush()
	fmt.Print("\nContinue? [y/N]  ")
	var userResponse string
	fmt.Scanln(&userResponse)

	if userResponse != "y" && userResponse != "Y" {
		err = fmt.Errorf(" NOTICE: User canceled migration")
		return completed, err
	}

	println()

	for _, migrationPlan := range allPlans {

		// If we have a vcsTokenID for this org use it. Otherwise, get one from the api
		var vcsTokenID string
		oldID, haveOneAlready := vcsTokenIDs[migrationPlan.NewOrg]
		if haveOneAlready {
			vcsTokenID = oldID
		} else {
			tokenID, err := getVCSToken(vcsUsername, migrationPlan.NewOrg, tfToken)
			if err != nil {
				fmt.Printf(
					"Skipping workspace %s because of an error getting the VCS Token ID for %s with %s ...\n %s\n",
					migrationPlan.NewName,
					vcsUsername,
					migrationPlan.NewOrg,
					err.Error(),
				)
				continue
			}
			vcsTokenIDs[migrationPlan.NewOrg] = tokenID
			vcsTokenID = tokenID
		}

		if vcsTokenID == "" {
			fmt.Printf(
				"Skipping workspace %s because no VCS Token ID was available for %s with %s",
				migrationPlan.NewName,
				vcsUsername,
				migrationPlan.NewOrg,
			)
			continue
		}

		fmt.Printf("  >>> Migrating %s/%s ... ", migrationPlan.NewOrg, migrationPlan.NewName)
		sensitiveVars, err = CreateAndPopulateV2Workspace(migrationPlan, tfToken, vcsTokenID)
		if err != nil {
			return completed, err
		}

		completed[migrationPlan.NewName] = sensitiveVars

		err = RunTFInit(migrationPlan, tfToken)
		if err != nil {
			return completed, err
		}
		println("Done")
	}

	return completed, nil
}



type OAuthTokens struct {
	Data []struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			CreatedAt           time.Time `json:"created-at"`
			ServiceProviderUser string    `json:"service-provider-user"`
			HasSSHKey           bool      `json:"has-ssh-key"`
		} `json:"attributes"`
		Relationships struct {
			OauthClient struct {
				Data struct {
					ID   string `json:"id"`
					Type string `json:"type"`
				} `json:"data"`
				Links struct {
					Related string `json:"related"`
				} `json:"links"`
			} `json:"oauth-client"`
		} `json:"relationships"`
		Links struct {
			Self string `json:"self"`
		} `json:"links"`
	} `json:"data"`
}


func getVCSToken(vcsUsername, orgName, tfToken string) (string, error) {
	url := fmt.Sprintf("https://app.terraform.io/api/v2/organizations/%s/oauth-tokens", orgName)
	headers := map[string]string{"Authorization": "Bearer " + tfToken}
	resp := CallAPI("GET", url, "", headers)

	defer resp.Body.Close()
	//bodyBytes, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(bodyBytes))

	var oauthTokens OAuthTokens

	if err := json.NewDecoder(resp.Body).Decode(&oauthTokens); err != nil {
		return "", err
	}

	vcsTokenID := ""

	for _, nextToken := range oauthTokens.Data {
		if nextToken.Attributes.ServiceProviderUser == vcsUsername {
			vcsTokenID = nextToken.ID
			break
		}
	}

	return vcsTokenID, nil
}