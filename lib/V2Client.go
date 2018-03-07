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

/*
 * @return string - The json needed for the data payload for the api call
 */
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

func CreateAllV2Variables(organization, workspaceName, tfToken string, tfVars []TFVar) {
	for _, nextVar := range tfVars {
		CreateV2Variable(organization, workspaceName, tfToken, nextVar)
	}
}

/*
 * @return string - The json needed for the data payload for the api call
 */
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

func CreateAndPopulateAllV2Workspaces(
	configFile, v1OrgName, v2OrgName string,
	tfToken, vcsTokenID string,
) error {
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
		v1WorkspaceName := line[0]
		v2Workspace := V2Workspace{
			Name:         line[1],
			TFVersion:    line[2],
			VCSRepoID:    line[3],
			VCSBranch:    line[4],
			TFWorkingDir: line[5],
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
