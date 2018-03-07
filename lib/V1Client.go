package lib

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type TFMeta struct {
	Total int `json:"total"`
}

type TFEnv struct {
	Username string `json:"username"`
	Name     string `json:"name"`
}

type TFState struct {
	UpdatedAt   string `json:"updated_at"`
	Environment TFEnv  `json:"environment"`
}

type TFAllStates struct {
	States []TFState `json:"states"`
	Meta   TFMeta    `json:"meta"`
}

type TFVar struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Hcl   bool   `json:"hcl"`
}

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

func getJsonFromFile(jsonFile string) TFAllStates {
	raw, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var contents TFAllStates
	json.Unmarshal(raw, &contents)
	return contents
}

/*
 * @param jsonFile The path and/or file name of a json file the contents
 *        of which match what is returned from the terraform api
 * @return a slice of strings -the environment names
 */
func GetAllEnvNamesFromJson(jsonFile string) []string {
	envNames := []string{}
	allStates := getJsonFromFile(jsonFile)

	for _, nextState := range allStates.States {
		envNames = append(envNames, nextState.Environment.Name)
	}

	return envNames
}

/*
 * @param tfToken The user's Terraform Enterprise Token
 * @return a slice of strings - the environment names from the v1 api
 */
func GetAllEnvNamesFromV1API(tfToken string) []string {
	baseURL := "https://atlas.hashicorp.com/api/v1/terraform/state?page="
	names := []string{}

	for page := 1; ; page++ {
		url := fmt.Sprintf("%s%d", baseURL, page)
		resp := CallApi("GET", url, "", map[string]string{"X-Atlas-Token": tfToken})

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

/*
 * @param filePathName - The path and/or name of the config file.
 * @param envNames - A slice of the names of the V1 environments
 */
func CreatePlanFromV1EnvNames(filePathName string, envNames []string) {

	file, err := os.Create(filePathName)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	line1 := []string{"V1 Env Name", "V2 Workspace Name", "Terraform Version",
		"VCS Repo Id", "VCS Branch", "Terraform Working Directory"}

	err = writer.Write(line1)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	for _, value := range envNames {
		newData := []string{value, value, "", "", "", ""}
		err := writer.Write(newData)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
}

func GetTFVarsFromV1Config(organization, envName, tfToken string) ([]TFVar, error) {

	url := fmt.Sprintf(
		"https://atlas.hashicorp.com/api/v1/terraform/configurations/%s/%s/versions/latest",
		organization,
		envName,
	)

	headers := map[string]string{
		"X-Atlas-Token": tfToken,
	}
	resp := CallApi("POST", url, "", headers)

	defer resp.Body.Close()

	var tfConfig TFConfig

	if err := json.NewDecoder(resp.Body).Decode(&tfConfig); err != nil {
		return nil, err
	}

	return tfConfig.Version.TfVars, nil
}
