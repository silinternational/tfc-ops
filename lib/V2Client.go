package lib

import (
	"bytes"
	//"io/ioutil"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

type V2UpdateConfig struct {
	Organization          string
	NewOrganization       string
	Workspace             string
	AtlasToken            string
	SearchString          string //  must be an exact case-insensitive match (i.e. not a partial match)
	NewValue              string
	AddKeyIfNotFound      bool // If true, then SearchOnVariableValue will be treated as false
	SearchOnVariableValue bool // If false, then will filter on variable key
	DryRunMode            bool
	SensitiveVariable     bool // Whether to mark the variable as sensitive
}

type V2CloneConfig struct {
	Organization                string
	NewOrganization             string
	SourceWorkspace             string
	NewWorkspace                string
	NewVCSTokenID               string
	AtlasToken                  string
	AtlasTokenDestination       string
	CopyState                   bool
	CopyVariables               bool
	DifferentDestinationAccount bool
}

// V2Workspace holds the information needed to create a v2 terraform workspace
type V2Workspace struct {
	Name         string
	TFVersion    string
	VCSRepoID    string
	VCSBranch    string
	TFWorkingDir string
}

// V2Var is what is returned by the api for one variable
type V2Var struct {
	ID        string `json:"-"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	Sensitive bool   `json:"sensitive"`
	Category  string `json:"category"`
	Hcl       bool   `json:"hcl"`
}

// V2VarsResponse is what is returned by the api when requesting the variables of a workspace
type V2VarsResponse struct {
	Data []struct {
		ID            string `json:"id"`
		Type          string `json:"type"`
		Variable      V2Var  `json:"attributes"`
		Relationships struct {
			Configurable struct {
				Data struct {
					ID   string `json:"id"`
					Type string `json:"type"`
				} `json:"data"`
				Links struct {
					Related string `json:"related"`
				} `json:"links"`
			} `json:"configurable"`
		} `json:"relationships"`
		Links struct {
			Self string `json:"self"`
		} `json:"links"`
	} `json:"data"`
}

// V2WorkspaceData is what is returned by the api for each workspace
type V2WorkspaceData struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		Name             string    `json:"name"`
		Environment      string    `json:"environment"`
		AutoApply        bool      `json:"auto-apply"`
		Locked           bool      `json:"locked"`
		CreatedAt        time.Time `json:"created-at"`
		WorkingDirectory string    `json:"working-directory"`
		VCSRepo          struct {
			Branch  string `json:"branch"`
			ID      string `json:"identifier"`
			TokenID string `json:"oauth-token-id"`
		} `json:"vcs-repo"`
		TerraformVersion string `json:"terraform-version"`
		Permissions      struct {
			CanUpdate         bool `json:"can-update"`
			CanDestroy        bool `json:"can-destroy"`
			CanQueueDestroy   bool `json:"can-queue-destroy"`
			CanQueueRun       bool `json:"can-queue-run"`
			CanUpdateVariable bool `json:"can-update-variable"`
			CanLock           bool `json:"can-lock"`
			CanReadSettings   bool `json:"can-read-settings"`
		} `json:"permissions"`
		Actions struct {
			IsDestroyable bool `json:"is-destroyable"`
		} `json:"actions"`
	} `json:"attributes"`
	Relationships struct {
		Organization struct {
			Data struct {
				ID   string `json:"id"`
				Type string `json:"type"`
			} `json:"data"`
		} `json:"organization"`
		LatestRun struct {
			Data interface{} `json:"data"`
		} `json:"latest-run"`
		CurrentRun struct {
			Data interface{} `json:"data"`
		} `json:"current-run"`
	} `json:"relationships"`
	Links struct {
		Self string `json:"self"`
	} `json:"links"`
}

func (v *V2WorkspaceData) AttributeByLabel(label string) (string, error) {
	switch strings.ToLower(label) {
	case "id":
		return v.ID, nil
	case "name":
		return v.Attributes.Name, nil
	case "createdat":
		return v.Attributes.CreatedAt.String(), nil
	case "environment":
		return v.Attributes.Environment, nil
	case "workingdirectory":
		return v.Attributes.WorkingDirectory, nil
	case "terraformversion":
		return v.Attributes.TerraformVersion, nil
	case "vcsrepo":
		return v.Attributes.VCSRepo.ID, nil
	}

	return "", fmt.Errorf("Attribute label not valid: %s", label)
}

// V2WorkspaceJSON is what is returned by the api when requesting the data for a workspace
type V2WorkspaceJSON struct {
	Data V2WorkspaceData `json:"data"`
}

// AllV2WorkspacesJSON is what is returned by the api when requesting the data for all workspaces
type AllV2WorkspacesJSON struct {
	Data []V2WorkspaceData `json:"data"`
}

// TeamWorkspaceData is what is returned by the api for one team access object for a workspace
type TeamWorkspaceData struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		Access string `json:"access"`
	} `json:"attributes"`
	Relationships struct {
		Team struct {
			Data struct {
				ID   string `json:"id"`
				Type string `json:"type"`
			} `json:"data"`
			Links struct {
				Related string `json:"related"`
			} `json:"links"`
		} `json:"team"`
		Workspace struct {
			Data struct {
				ID   string `json:"id"`
				Type string `json:"type"`
			} `json:"data"`
			Links struct {
				Related string `json:"related"`
			} `json:"links"`
		} `json:"workspace"`
	} `json:"relationships"`
	Links struct {
		Self string `json:"self"`
	} `json:"links"`
}

// AllTeamWorkspaceData is what is returned by the api when requesting the team access data for a workspace
type AllTeamWorkspaceData struct {
	Data []TeamWorkspaceData `json:"data"`
}

// TFVar matches the attributes of a terraform environment/workspace's variable
type TFVar struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	Hcl       bool   `json:"hcl"`
	Sensitive bool   `json:"sensitive"`
}

// TFConfig matches the json return value of the v1 terraform configurations api
type TFConfig struct {
	Version struct {
		Version  int `json:"version"`
		Metadata struct {
			Foo string `json:"foo"`
		} `json:"metadata"`
		TfVars    []TFVar           `json:"tf_vars"`
		Variables map[string]string `json:"variables"`
	} `json:"version"`
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
      "sensitive":%t
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
`, tfVar.Key, tfVar.Value, tfVar.Hcl, tfVar.Sensitive, organization, workspaceName)
}

// GetUpdateV2VariablePayload returns the json needed to make a Post to the
// v2 terraform vars api
func GetUpdateV2VariablePayload(organization, workspaceName, variableID string, tfVar TFVar) string {
	return fmt.Sprintf(`
{
  "data": {
    "id":"%s",
    "type":"vars",
    "attributes": {
      "key":"%s",
      "value":"%s",
      "category":"terraform",
      "description":"",
      "hcl":%t,
      "sensitive":%t
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
`, variableID, tfVar.Key, tfVar.Value, tfVar.Hcl, tfVar.Sensitive, organization, workspaceName)
}

func GetV2AllWorkspaceData(organization, tfToken string) ([]V2WorkspaceData, error) {

	baseURL := fmt.Sprintf("https://app.terraform.io/api/v2/organizations/%s/workspaces?page%%5Bnumber%%5D=", organization)

	headers := map[string]string{
		"Authorization": "Bearer " + tfToken,
		"Content-Type":  "application/vnd.api+json",
	}

	allWsData := []V2WorkspaceData{}

	for page := 1; ; page++ {
		url := fmt.Sprintf("%s%d", baseURL, page)
		resp := CallAPI("GET", url, "", headers)

		defer resp.Body.Close()
		//bodyBytes, _ := ioutil.ReadAll(resp.Body)
		//fmt.Println(string(bodyBytes))

		var nextWsData AllV2WorkspacesJSON

		if err := json.NewDecoder(resp.Body).Decode(&nextWsData); err != nil {
			return []V2WorkspaceData{}, fmt.Errorf("Error getting all workspaces' data for %s:%s\n%s", organization, err.Error())
		}
		allWsData = append(allWsData, nextWsData.Data...)

		// If there isn't a whole page of contents, then we're on the last one.
		if len(nextWsData.Data) < 20 {
			break
		}
	}
	return allWsData, nil
}

func GetV2WorkspaceData(organization, workspaceName, tfToken string) (V2WorkspaceJSON, error) {

	url := fmt.Sprintf(
		"https://app.terraform.io/api/v2/organizations/%s/workspaces/%s",
		organization,
		workspaceName,
	)

	headers := map[string]string{
		"Authorization": "Bearer " + tfToken,
		"Content-Type":  "application/vnd.api+json",
	}
	resp := CallAPI("GET", url, "", headers)

	defer resp.Body.Close()
	//bodyBytes, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(bodyBytes))

	var v2WsData V2WorkspaceJSON

	if err := json.NewDecoder(resp.Body).Decode(&v2WsData); err != nil {
		return V2WorkspaceJSON{}, fmt.Errorf("Error getting workspace data for %s:%s\n%s", organization, workspaceName, err.Error())
	}

	return v2WsData, nil
}

//  GetVarsFromV2 returns a list of Terraform variables for a given workspace
func GetVarsFromV2(organization, workspaceName, tfToken string) ([]V2Var, error) {

	url := fmt.Sprintf(
		"https://app.terraform.io/api/v2/vars?filter%%5Borganization%%5D%%5Bname%%5D=%s&filter%%5Bworkspace%%5D%%5Bname%%5D=%s",
		organization,
		workspaceName,
	)

	headers := map[string]string{
		"Authorization": "Bearer " + tfToken,
		"Content-Type":  "application/vnd.api+json",
	}
	resp := CallAPI("GET", url, "", headers)

	defer resp.Body.Close()
	// bodyBytes, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(bodyBytes))

	var v2Resp V2VarsResponse

	if err := json.NewDecoder(resp.Body).Decode(&v2Resp); err != nil {
		return []V2Var{}, fmt.Errorf("Error getting variables for %s:%s ...\n%s", organization, workspaceName, err.Error())
	}

	variables := []V2Var{}
	for _, data := range v2Resp.Data {
		data.Variable.ID = data.ID // push the ID down into the Variable for future reference
		variables = append(variables, data.Variable)
	}

	return variables, nil
}

func GetAllWorkSpacesVarsFromV2(wsData []V2WorkspaceData, organization, keyContains, valueContains, tfToken string) (map[string][]V2Var, error) {
	allVars := map[string][]V2Var{}

	for _, ws := range wsData {
		wsName := ws.Attributes.Name

		wsVars, err := GetMatchingVarsFromV2(organization, wsName, tfToken, keyContains, valueContains)
		if err != nil {
			return nil, err
		}
		allVars[wsName] = wsVars
	}

	return allVars, nil
}

func GetMatchingVarsFromV2(organization string, wsName string, tfToken string, keyContains string, valueContains string) ([]V2Var, error) {

	vars, err := GetVarsFromV2(organization, wsName, tfToken)
	if err != nil {
		err := fmt.Errorf("Error getting variables for %s:%s\n%s", organization, wsName, err.Error())
		return []V2Var{}, err
	}

	var wsVars []V2Var

	for _, v := range vars {
		if keyContains != "" && strings.Contains(v.Key, keyContains) {
			wsVars = append(wsVars, v)
			continue
		}
		if valueContains != "" && strings.Contains(v.Value, valueContains) {
			wsVars = append(wsVars, v)
		}
	}

	return wsVars, nil
}

// GetTeamAccessFromV2 returns the team access data from an existing workspace
func GetTeamAccessFromV2(workspaceID, tfToken string) (AllTeamWorkspaceData, error) {
	url := fmt.Sprintf(
		"https://app.terraform.io/api/v2/team-workspaces?filter%%5Bworkspace%%5D%%5Bid%%5D=%s",
		workspaceID,
	)

	headers := map[string]string{
		"Authorization": "Bearer " + tfToken,
		"Content-Type":  "application/vnd.api+json",
	}

	resp := CallAPI("GET", url, "", headers)

	defer resp.Body.Close()

	var allTeamData AllTeamWorkspaceData

	if err := json.NewDecoder(resp.Body).Decode(&allTeamData); err != nil {
		return AllTeamWorkspaceData{}, fmt.Errorf("Error getting team workspace data for %s\n%s", workspaceID, err.Error())
	}

	return allTeamData, nil
}

func getAssignTeamAccessPayload(accessLevel, workspaceID, teamID string) string {
	return fmt.Sprintf(`
{
  "data": {
    "attributes": {
      "access":"%s"
    },
    "relationships": {
      "workspace": {
        "data": {
          "type":"workspaces",
          "id":"%s"
        }
      },
      "team": {
        "data": {
          "type":"teams",
          "id":"%s"
        }
      }
    },
    "type":"team-workspaces"
  }
}
`, accessLevel, workspaceID, teamID)
}

// AssignTeamAccessOnV2 assigns the requested team access to a workspace on Terraform Enterprise V.2
func AssignTeamAccessOnV2(workspaceID, tfToken string, allTeamData AllTeamWorkspaceData) {
	url := fmt.Sprintf("https://app.terraform.io/api/v2/team-workspaces")

	headers := map[string]string{
		"Authorization": "Bearer " + tfToken,
		"Content-Type":  "application/vnd.api+json",
	}

	for _, teamData := range allTeamData.Data {
		postData := getAssignTeamAccessPayload(
			teamData.Attributes.Access,
			workspaceID,
			teamData.Relationships.Team.Data.ID,
		)

		resp := CallAPI("POST", url, postData, headers)
		defer resp.Body.Close()
	}
	return
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
func GetCreateV2WorkspacePayload(oc OpsConfig, vcsTokenID string) string {
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
  `, oc.NewName, oc.TerraformVersion, oc.Directory, oc.RepoID, vcsTokenID, oc.Branch)
}

// UpdateV2Variable makes a v2 terraform vars api post to update a variable
// for a given organization and v2 workspace
func UpdateV2Variable(organization, workspaceName, variableID, tfToken string, tfVar TFVar) {
	url := fmt.Sprintf("https://app.terraform.io/api/v2/vars/%s", variableID)

	ConvertHCLVariable(&tfVar)

	patchData := GetUpdateV2VariablePayload(organization, workspaceName, variableID, tfVar)

	headers := map[string]string{
		"Authorization": "Bearer " + tfToken,
		"Content-Type":  "application/vnd.api+json",
	}
	resp := CallAPI("PATCH", url, patchData, headers)

	defer resp.Body.Close()
	// bodyBytes, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(bodyBytes))
	return
}

// CreateV2Workspace makes a v2 terraform workspaces api Post to create a
// workspace for a given organization, including setting up its vcs repo integration
func CreateV2Workspace(
	oc OpsConfig,
	tfToken, vcsTokenID string,
) (string, error) {
	url := fmt.Sprintf(
		"https://app.terraform.io/api/v2/organizations/%s/workspaces",
		oc.NewOrg,
	)

	postData := GetCreateV2WorkspacePayload(oc, vcsTokenID)

	headers := map[string]string{
		"Authorization": "Bearer " + tfToken,
		"Content-Type":  "application/vnd.api+json",
	}
	resp := CallAPI("POST", url, postData, headers)

	defer resp.Body.Close()
	// bodyBytes, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(bodyBytes))

	var v2WsData V2WorkspaceJSON

	if err := json.NewDecoder(resp.Body).Decode(&v2WsData); err != nil {
		return "", fmt.Errorf("error getting created workspace data: %s\n", err)
	}
	return v2WsData.Data.ID, nil
}


// RunTFInit ...
//  - removes old terraform.tfstate files
//  - runs terraform init with old versions
//  - runs terraform init with new version
// NOTE: This procedure can be used to copy/migrate a workspace's state to a new one.
//  (see the -backend-config mention below and the backend.tf file in this repo)
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
		return fmt.Errorf("Error setting %s environment variable to source value: %s", tokenEnv, err)
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
		return fmt.Errorf("Error setting %s environment variable to destination value: %s", tokenEnv, err)
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
		return fmt.Errorf("Error resetting %s environment variable back to source value: %s", tokenEnv, err)
	}

	return nil
}


// CloneV2Workspace gets the data, variables and team access data for an existing Terraform Enterprise workspace
//  and then creates a clone of it with the same data.
// If the copyVariables param is set to true, then all the non-sensitive variable values will be added to the new
//   workspace.  Otherwise, they will be set to "REPLACE_THIS_VALUE"
func CloneV2Workspace(cfg V2CloneConfig) ([]string, error) {

	v2WsData, err := GetV2WorkspaceData(cfg.Organization, cfg.SourceWorkspace, cfg.AtlasToken)
	if err != nil {
		return []string{}, err
	}

	variables, err := GetVarsFromV2(cfg.Organization, cfg.SourceWorkspace, cfg.AtlasToken)
	if err != nil {
		return []string{}, err
	}

	if !cfg.DifferentDestinationAccount {
		cfg.NewOrganization = cfg.Organization
		cfg.NewVCSTokenID = v2WsData.Data.Attributes.VCSRepo.ID
	}

	oc := OpsConfig{
		SourceOrg:        cfg.Organization,
		SourceName:       v2WsData.Data.Attributes.Name,
		NewOrg:           cfg.NewOrganization,
		NewName:          cfg.NewWorkspace,
		TerraformVersion: v2WsData.Data.Attributes.TerraformVersion,
		RepoID:           v2WsData.Data.Attributes.VCSRepo.ID,
		Branch:           v2WsData.Data.Attributes.VCSRepo.Branch,
		Directory:        v2WsData.Data.Attributes.WorkingDirectory,
	}

	sensitiveVars := []string{}
	sensitiveValue := "TF_ENTERPRISE_SENSITIVE_VAR"
	defaultValue := "REPLACE_THIS_VALUE"

	tfVars := []TFVar{}
	var tfVar TFVar

	for _, nextVar := range variables {
		if cfg.CopyVariables {
			tfVar = TFVar{
				Key:   nextVar.Key,
				Value: nextVar.Value,
				Hcl:   nextVar.Hcl,
			}
		} else {
			tfVar = TFVar{
				Key:   nextVar.Key,
				Value: defaultValue,
				Hcl:   nextVar.Hcl,
			}
		}
		if nextVar.Value == sensitiveValue {
			sensitiveVars = append(sensitiveVars, nextVar.Key)
		}
		tfVars = append(tfVars, tfVar)
	}

	if cfg.DifferentDestinationAccount {
		CreateV2Workspace(oc, cfg.AtlasTokenDestination, cfg.NewVCSTokenID)
		CreateAllV2Variables(oc.NewOrg, oc.NewName, cfg.AtlasTokenDestination, tfVars)

		if cfg.CopyState {
			if err := RunTFInit(oc, cfg.AtlasToken, cfg.AtlasTokenDestination); err != nil {
				return sensitiveVars, err
			}
		}

		return sensitiveVars, nil
	}

	CreateV2Workspace(oc, cfg.AtlasToken, v2WsData.Data.Attributes.VCSRepo.TokenID)
	CreateAllV2Variables(oc.NewOrg, oc.NewName, cfg.AtlasToken, tfVars)

	// Get Team Access Data for source Workspace
	allTeamData, err := GetTeamAccessFromV2(v2WsData.Data.ID, cfg.AtlasToken)
	if err != nil {
		return sensitiveVars, err
	}

	// Get new Workspace data for its ID
	newV2WsData, err := GetV2WorkspaceData(cfg.Organization, cfg.NewWorkspace, cfg.AtlasToken)
	if err != nil {
		return sensitiveVars, err
	}

	AssignTeamAccessOnV2(newV2WsData.Data.ID, cfg.AtlasToken, allTeamData)

	return sensitiveVars, nil
}

// AddOrUpdateV2Variable adds or updates an existing Terraform Enterprise workspace variable
// If the copyVariables param is set to true, then all the non-sensitive variable values will be added to the new
//   workspace.  Otherwise, they will be set to "REPLACE_THIS_VALUE"
func AddOrUpdateV2Variable(cfg V2UpdateConfig) (string, error) {
	variables, err := GetVarsFromV2(cfg.Organization, cfg.Workspace, cfg.AtlasToken)
	if err != nil {
		return "", err
	}

	loweredSearchString := strings.ToLower(cfg.SearchString)

	for _, nextVar := range variables {
		oldValue := nextVar.Value
		if cfg.SearchOnVariableValue {
			if strings.ToLower(nextVar.Value) != loweredSearchString {
				continue
			}
			// Found a match
			tfVar := TFVar{Key: nextVar.Key, Value: cfg.NewValue, Hcl: false, Sensitive: cfg.SensitiveVariable}
			if !cfg.DryRunMode {
				UpdateV2Variable(cfg.Organization, cfg.Workspace, nextVar.ID, cfg.AtlasToken, tfVar)
			}
			return fmt.Sprintf("Replaced the value of %s from %s to %s", nextVar.Key, oldValue, cfg.NewValue), nil
		}

		// Search on variable key, since search on value is not true
		if strings.ToLower(nextVar.Key) != loweredSearchString {
			continue
		}

		// Found a match
		// Only add if there isn't a match
		if cfg.AddKeyIfNotFound {
			return "", errors.New("addKeyIfNotFound was set to true but a variable already exists with key " + nextVar.Key)
		}

		tfVar := TFVar{Key: nextVar.Key, Value: cfg.NewValue, Hcl: false, Sensitive: cfg.SensitiveVariable}

		if !cfg.DryRunMode {
			UpdateV2Variable(cfg.Organization, cfg.Workspace, nextVar.ID, cfg.AtlasToken, tfVar)
		}
		return fmt.Sprintf("Replaced the value of %s from %s to %s", nextVar.Key, oldValue, cfg.NewValue), nil
	}

	// At this point, we haven't found a match
	if cfg.AddKeyIfNotFound {
		tfVar := TFVar{Key: cfg.SearchString, Value: cfg.NewValue, Hcl: false, Sensitive: cfg.SensitiveVariable}

		if !cfg.DryRunMode {
			CreateV2Variable(cfg.Organization, cfg.Workspace, cfg.AtlasToken, tfVar)
		}
		return fmt.Sprintf("Added variable %s = %s", cfg.SearchString, cfg.NewValue), nil
	}

	return "No match found and no variable added", nil
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

