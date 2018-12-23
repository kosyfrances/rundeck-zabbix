package lib

import (
	"errors"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
)

type Helper struct {
}

func FileExists(name string) bool {
	if _, err := os.Stat(name); !os.IsNotExist(err) {
		return true
	}
	return false
}

func GetStringMapKeys(data map[string]string) []string {
	keys := []string{}
	for key, _ := range data {
		keys = append(keys, key)
	}
	return keys
}

func GetStringValueIndex(data []string, value string) (int, error) {
	for index, item := range data {
		if item == value {
			return index, nil
		}
	}
	return 0, errors.New("value not found in array")
}

func GetStringDelimitedMapKeys(data map[string]string, delimter string) string {
	dataKeys := GetStringMapKeys(data)
	return strings.Join(dataKeys, delimter)
}

func ValidateComponent(component string, validComponent map[string]string) string {
	if _, ok := validComponent[component]; ok == false {
		return "Invalid component set. Supported value(s): " + GetStringDelimitedMapKeys(validComponent, " | ")
	}
	return ""
}

func TableMaker(header []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	// table.SetAutoMergeCells(true)
	table.SetRowLine(true)

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}
