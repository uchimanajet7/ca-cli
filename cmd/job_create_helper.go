package cmd

import (
	"errors"
	"reflect"
	"strings"
)

// ca job rule_type
type jobRuleType struct {
	Name        string
	IsRuleValue bool
}

// see also:
// https://cloudautomator.com/api_docs/v1/api.html#header--1
var jobRuleTypeList = []jobRuleType{
	jobRuleType{
		Name:        "cron",
		IsRuleValue: true,
	},
	jobRuleType{
		Name:        "immediate_execution",
		IsRuleValue: false,
	},
	jobRuleType{
		Name:        "webhook",
		IsRuleValue: false,
	},
	jobRuleType{
		Name:        "sqs",
		IsRuleValue: true,
	},
	jobRuleType{
		Name:        "amazon_snssqs",
		IsRuleValue: false,
	},
}

func isValidJobRuleType(name string) (bool, bool) {
	for _, v := range jobRuleTypeList {
		if strings.ToLower(name) == strings.ToLower(v.Name) {
			return true, v.IsRuleValue
		}
	}
	return false, false
}

// see also:
// https://cloudautomator.com/api_docs/v1/api.html#header--2
func createJobParseObjectParameter(raw string, msg string) (*map[string]interface{}, error) {
	// parse raw string
	rawSlice := strings.Split(raw, ",")
	parsedMap := make(map[string]interface{}, len(rawSlice))
	prevKey := ""
	for _, v := range rawSlice {

		if strings.Index(v, "=") > 0 {
			slice := strings.Split(v, "=")

			if len(slice) < 2 {
				return nil, errors.New(msg)
			}
			mapKey := strings.TrimSpace(slice[0])
			mapVal := strings.TrimSpace(slice[1])
			parsedMap[mapKey] = mapVal
			prevKey = mapKey
			continue
		}

		// not exsit "="
		prevVal := parsedMap[prevKey]
		if reflect.TypeOf(prevVal).Name() == "string" {
			prevVal = []string{prevVal.(string)}
		}
		parsedMap[prevKey] = append(prevVal.([]string), strings.TrimSpace(v))
	}

	return &parsedMap, nil
}
