package lib

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
)

type MigrationPlan struct {
	LegacyName       string
	NewName          string
	TerraformVersion string
	RepoID           string
	Branch           string
	Directory        string
}

func (p *MigrationPlan) Array() []string {
	return []string{
		p.LegacyName,
		p.NewName,
		p.TerraformVersion,
		p.RepoID,
		p.Branch,
		p.Directory,
	}
}

func (f *MigrationPlan) getColNames() []string {
	val := reflect.ValueOf(f).Elem()
	cols := []string{}

	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		cols = append(cols, typeField.Name)
	}

	return cols
}

/*
 * CreatePlanFile generates a CSV file with plan details
 * @param filePathName - The path and/or name of the config file.
 * @param envNames - A slice of the names of the V1 environments
 */
func CreatePlanFile(filePathName string, plans []MigrationPlan) {
	file, err := os.Create(filePathName)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	emptyPlan := MigrationPlan{}
	line1 := emptyPlan.getColNames()

	err = writer.Write(line1)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	for _, plan := range plans {
		newData := plan.Array()
		err := writer.Write(newData)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}

}

func GetBasePlansFromEnvNames(envNames []string) []MigrationPlan {
	basePlans := []MigrationPlan{}
	for _, env := range envNames {
		newPlan := MigrationPlan{
			LegacyName:       env,
			NewName:          env,
			TerraformVersion: "",
			RepoID:           "",
			Branch:           "",
			Directory:        "",
		}
		basePlans = append(basePlans, newPlan)
	}

	return basePlans
}
