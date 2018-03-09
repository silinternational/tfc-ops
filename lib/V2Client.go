package lib

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type V2Workspace struct {
	Name         string
	TFVersion    string
	VCSRepoID    string
	VCSBranch    string
	TFWorkingDir string
}

/*
 * @param tfVar - Changes the struct in place by escaping
 *  the double quotes and line endings in the Value attribute
 */
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
	resp := CallApi("POST", url, postData, headers)

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
func GetCreateV2WorkspacePayload(
	name, tfVersion, workingDir string,
	vcsID, vcsTokenID, vcsBranch string,
) string {
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
  `, name, tfVersion, workingDir, vcsID, vcsTokenID, vcsBranch)
}

// CreateV2Workspace makes a v2 terraform workspaces api Post to create a
// workspace for a given organization, including setting up its vcs repo integration
func CreateV2Workspace(
	organization, name, tfToken, tfVersion, workingDir string,
	vcsID, vcsTokenID, vcsBranch string,
) {
	url := fmt.Sprintf(
		"https://app.terraform.io/api/v2/organizations/%s/workspaces",
		organization,
	)

	postData := GetCreateV2WorkspacePayload(name, tfVersion, workingDir, vcsID, vcsTokenID, vcsBranch)

	headers := map[string]string{
		"Authorization": "Bearer " + tfToken,
		"Content-Type":  "application/vnd.api+json",
	}
	resp := CallApi("POST", url, postData, headers)

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
	v2Workspace V2Workspace,
	v1WorkspaceName, v1OrgName, v2OrgName, tfToken, vcsTokenID string,
) error {

	v1Vars, err := GetTFVarsFromV1Config(v1OrgName, v1WorkspaceName, tfToken)
	if err != nil {
		return err
	}

	CreateV2Workspace(
		v2OrgName,
		v2Workspace.Name,
		tfToken,
		v2Workspace.TFVersion,
		v2Workspace.TFWorkingDir,
		v2Workspace.VCSRepoID,
		vcsTokenID,
		v2Workspace.VCSBranch,
	)

	CreateAllV2Variables(v2OrgName, v2Workspace.Name, tfToken, v1Vars)

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
func CreateAndPopulateAllV2Workspaces(configFile, tfToken, vcsTokenID string) error {
	// Get config contents
	csvFile, err := os.Open(configFile)
	if err != nil {
		return err
	}

	reader := csv.NewReader(bufio.NewReader(csvFile))

	// Throw away the first line with its column headers
	_, err = reader.Read()
	if err != nil {
		return err
	}

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		v1OrgName := line[0]
		v1WorkspaceName := line[1]
		v2OrgName := line[2]

		v2Workspace := V2Workspace{
			Name:         line[3],
			TFVersion:    line[4],
			VCSRepoID:    line[5],
			VCSBranch:    line[6],
			TFWorkingDir: line[7],
		}

		err = CreateAndPopulateV2Workspace(
			v2Workspace,
			v1WorkspaceName,
			v1OrgName,
			v2OrgName,
			tfToken,
			vcsTokenID,
		)
		if err != nil {
			return err
		}

	}
	return nil
}
