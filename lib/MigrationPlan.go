package lib

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"strings"
)

// MigrationPlan represents one row of the plan.csv file's contents
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

// NewMigrationPlan creates a new MigrationPlan from a csv row
func NewMigrationPlan(values []string) (MigrationPlan, error) {
	var mp MigrationPlan
	colNames := mp.getColNames()
	columnCount := len(colNames)
	if len(values) < columnCount {
		return mp, fmt.Errorf("Too few values to create MigrationPlan. Need: %d, but only got %d",
			columnCount, len(values),
		)
	}

	for index, nextValue := range values {
		if nextValue == "" {
			return mp, fmt.Errorf("Empty string not allowed in column %d", index+1)
		}
	}

	version := values[4]
	versionParts := strings.Split(version, ".")

	if len(versionParts) != 3 {
		return mp, fmt.Errorf("The version value should have three sets of digits separated by dots. Got %s",
			version,
		)
	}

	repoID := values[5]
	repoIDParts := strings.Split(repoID, `/`)

	if len(repoIDParts) < 2 {
		return mp, fmt.Errorf("The RepoID value should include a slash. Got %s", repoID)
	}

	mp = MigrationPlan{
		LegacyOrg:        values[0],
		LegacyName:       values[1],
		NewOrg:           values[2],
		NewName:          values[3],
		TerraformVersion: version,
		RepoID:           repoID,
		Branch:           values[6],
		Directory:        values[7],
	}

	return mp, nil
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
