package lib

import (
	"encoding/json"
	"fmt"
	"os"
)

// TFMeta matches the meta element
type TFMeta struct {
	Total int `json:"total"`
}

// TFEnv matches the contents of the environment element
type TFEnv struct {
	Username string `json:"username"`
	Name     string `json:"name"`
}

// TFState matches one entry in the states list
type TFState struct {
	UpdatedAt   string `json:"updated_at"`
	Environment TFEnv  `json:"environment"`
}

// TFAllStates matches the return value of a call to the v1 terraform state api
type TFAllStates struct {
	States []TFState `json:"states"`
	Meta   TFMeta    `json:"meta"`
}

// TFVar matches the attributes of a terraform environment/workspace's variable
type TFVar struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Hcl   bool   `json:"hcl"`
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

// GetAllEnvNamesFromV1API calls the v1 terraform state api.
// @param tfToken The user's Terraform Enterprise Token
// @return a slice of strings - the environment names from the v1 api
func GetAllEnvNamesFromV1API(tfToken string) []string {
	baseURL := "https://atlas.hashicorp.com/api/v1/terraform/state?page="
	names := []string{}

	for page := 1; ; page++ {
		url := fmt.Sprintf("%s%d", baseURL, page)
		resp := CallAPI("GET", url, "", map[string]string{"X-Atlas-Token": tfToken})

		defer resp.Body.Close()

		var statesFromPage TFAllStates

		if err := json.NewDecoder(resp.Body).Decode(&statesFromPage); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		for _, nextState := range statesFromPage.States {
			names = append(names, nextState.Environment.Name)
		}

		// If there isn't a whole page of contents, then we're on the last one.
		if len(statesFromPage.States) < 20 {
			break
		}
	}

	return names
}

// GetTFVarsFromV1Config calls the v1 terraform configurations api and
// returns a list of Terraform variables for a given environment
func GetTFVarsFromV1Config(organization, envName, tfToken string) ([]TFVar, error) {

	url := fmt.Sprintf(
		"https://atlas.hashicorp.com/api/v1/terraform/configurations/%s/%s/versions/latest",
		organization,
		envName,
	)

	headers := map[string]string{
		"X-Atlas-Token": tfToken,
	}
	resp := CallAPI("GET", url, "", headers)

	defer resp.Body.Close()
	// bodyBytes, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(bodyBytes))

	var tfConfig TFConfig

	if err := json.NewDecoder(resp.Body).Decode(&tfConfig); err != nil {
		return nil, err
	}

	return tfConfig.Version.TfVars, nil
}
