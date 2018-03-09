package lib

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
)

type MigrationPlan struct {
	LegacyOrg        string
	LegacyName       string
	NewOrg           string
	NewName          string
	TerraformVersion string
	RepoID           string
	Branch           string
	Directory        string
}

// AsArray returns the values of the MigrationPlan attributes
func (p *MigrationPlan) AsArray() []string {
	return []string{
		p.LegacyOrg,
		p.LegacyName,
		p.NewOrg,
		p.NewName,
		p.TerraformVersion,
		p.RepoID,
		p.Branch,
		p.Directory,
	}
}

func (p *MigrationPlan) getColNames() []string {
	val := reflect.ValueOf(p).Elem()
	cols := []string{}

	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		cols = append(cols, typeField.Name)
	}

	return cols
}

// CreatePlanFile generates a CSV file with plan details
// @param filePathName - The path and/or name of the config file.
// @param envNames - A slice of the names of the V1 environments
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
		newData := plan.AsArray()
		err := writer.Write(newData)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
}

// GetBasePlansFromEnvNames creates the initial MigrationPlan structs to
// be used to create a plan.csv file
func GetBasePlansFromEnvNames(envNames []string, legacyOrg, newOrg string) []MigrationPlan {
	basePlans := []MigrationPlan{}
	for _, env := range envNames {
		newPlan := MigrationPlan{
			LegacyOrg:  legacyOrg,
			LegacyName: env,
			NewOrg:     newOrg,
			NewName:    env,
		}
		basePlans = append(basePlans, newPlan)
	}

	return basePlans
}
