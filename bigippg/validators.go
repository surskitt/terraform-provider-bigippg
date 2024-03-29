/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigippg

import (
	"fmt"
	"net"
	"reflect"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Validate the incoming set only contains values from the specified set
func validateSetValues(valid *schema.Set) schema.SchemaValidateFunc {
	return func(value interface{}, field string) (ws []string, errors []error) {
		if valid.Intersection(value.(*schema.Set)).Len() != value.(*schema.Set).Len() {
			errors = append(errors, fmt.Errorf("%q can only contain %v", field, value.(*schema.Set).List()))
		}
		return
	}
}

func validateStringValue(values []string) schema.SchemaValidateFunc {
	return func(value interface{}, field string) (ws []string, errors []error) {
		for _, v := range values {
			if v == value.(string) {
				return
			}
		}
		errors = append(errors, fmt.Errorf("%q must be one of %v", field, values))
		return
	}
}

func validateF5Name(value interface{}, field string) (ws []string, errors []error) {
	var values []string
	switch value.(type) {
	case *schema.Set:
		values = setToStringSlice(value.(*schema.Set))
		break
	case []string:
		values = value.([]string)
		break
	case *[]string:
		values = *(value.(*[]string))
		break
	case string:
		values = []string{value.(string)}
		break
	default:
		errors = append(errors, fmt.Errorf("Unknown type %v in validateF5Name", reflect.TypeOf(value)))
	}

	for _, v := range values {
		match, _ := regexp.MatchString("^/[\\w_\\-.]+/[\\w_\\-.:]+$", v)
		if !match {
			errors = append(errors, fmt.Errorf("%q must match /Partition/Name and contain letters, numbers or [._-:]. e.g. /Common/my-pool", field))
		}
	}
	return
}

func validateF5NameWithDirectory(value interface{}, field string) (ws []string, errors []error) {
	var values []string
	switch value.(type) {
	case *schema.Set:
		values = setToStringSlice(value.(*schema.Set))
		break
	case []string:
		values = value.([]string)
		break
	case *[]string:
		values = *(value.(*[]string))
		break
	case string:
		values = []string{value.(string)}
		break
	default:
		errors = append(errors, fmt.Errorf("Unknown type %v in validateF5Name", reflect.TypeOf(value)))
	}

	for _, v := range values {
		match, _ := regexp.MatchString("(^/[\\w_\\-.]+/[\\w_\\-.:]+/[\\w_\\-.:]+$)|(^/[\\w_\\-.]+/[\\w_\\-.:]+$)", v)
		if !match {
			errors = append(errors, fmt.Errorf("%q must match /Partition/Name or /Partition/Directory/Name  e.g. /Common/my-node or /Common/test/my-node", field))
		}
	}
	return
}

func validatePartitionName(value interface{}, field string) (ws []string, errors []error) {
	var values []string
	switch value.(type) {
	case *schema.Set:
		values = setToStringSlice(value.(*schema.Set))
		break
	case []string:
		values = value.([]string)
		break
	case *[]string:
		values = *(value.(*[]string))
		break
	case string:
		values = []string{value.(string)}
		break
	default:
		errors = append(errors, fmt.Errorf("Unknown type %v in validatePartitionName", reflect.TypeOf(value)))
	}

	for _, v := range values {
		match, _ := regexp.MatchString(`^[^/][^\s]+$`, v)
		if !match {
			errors = append(errors, fmt.Errorf("%q name should not start with `/`, e.g Common [or] test-partition are valid ", field))
		}
	}
	return
}

func validatePoolMemberName(value interface{}, field string) (ws []string, errors []error) {
	var values []string
	switch value.(type) {
	case *schema.Set:
		values = setToStringSlice(value.(*schema.Set))
		break
	case []string:
		values = value.([]string)
		break
	case *[]string:
		values = *(value.(*[]string))
		break
	case string:
		values = []string{value.(string)}
		break
	default:
		errors = append(errors, fmt.Errorf("Unknown type %v in validatePoolMemberName", reflect.TypeOf(value)))
	}

	for _, v := range values {
		if strings.Count(v, ":") >= 2 {
			match, _ := regexp.MatchString("^\\/[\\w_\\-.]+\\/[\\w_\\-.:]+.\\d+$", v)
			if !match {
				errors = append(errors, fmt.Errorf("%q must match /Partition/Node_Name:Port and contain letters, numbers or [:._-]. e.g. /Common/node1:80", field))
			}
		} else {
			match, _ := regexp.MatchString("^[\\w_\\-.]+:\\d+$", v)
			if !match {
				errors = append(errors, fmt.Errorf("%q must match Node-address:Port and Node Address is IP/FQDN. e.g. 1.1.1.1:80/www.google.com:80", field))
			}
		}
	}
	return
}

// IsValidIP tests that the argument is a valid IP address.
func IsValidIP(value string) bool {
	if net.ParseIP(value) == nil {
		return false
	}
	return true
}

func validateEnabledDisabled(value interface{}, field string) (ws []string, errors []error) {
	var values []string
	switch value.(type) {
	case *schema.Set:
		values = setToStringSlice(value.(*schema.Set))
	case []string:
		values = value.([]string)
		break
	case *[]string:
		values = *(value.(*[]string))
		break
	case string:
		values = []string{value.(string)}
		break
	default:
		errors = append(errors, fmt.Errorf("Unknown type %v in validateEnabledDisabled", reflect.TypeOf(value)))
	}

	for _, v := range values {
		match, _ := regexp.MatchString("^enabled$|^disabled$", v)
		if !match {
			errors = append(errors, fmt.Errorf("%q must match as enabled or disabled", field))
		}
	}
	return
}

func validateReqPrefDisabled(value interface{}, field string) (ws []string, errors []error) {
	var values []string
	switch value.(type) {
	case *schema.Set:
		values = setToStringSlice(value.(*schema.Set))
	case []string:
		values = value.([]string)
	case *[]string:
		values = *(value.(*[]string))
	case string:
		values = []string{value.(string)}
	default:
		errors = append(errors, fmt.Errorf("Unknown type %v in validateReqPrefDisabled", reflect.TypeOf(value)))
	}

	for _, v := range values {
		match, _ := regexp.MatchString("^required$|^preferred$|^disabled$", v)
		if !match {
			errors = append(errors, fmt.Errorf("%q must match as required, preferred, or disabled", field))
		}
	}
	return
}

func validateDataGroupType(value interface{}, field string) (ws []string, errors []error) {
	var values []string
	switch value.(type) {
	case *schema.Set:
		values = setToStringSlice(value.(*schema.Set))
	case []string:
		values = value.([]string)
	case *[]string:
		values = *(value.(*[]string))
	case string:
		values = []string{value.(string)}
	default:
		errors = append(errors, fmt.Errorf("Unknown type %v in validateDataGroupType", reflect.TypeOf(value)))
	}

	for _, v := range values {
		match, _ := regexp.MatchString("^string$|^ip$|^integer$", v)
		if !match {
			errors = append(errors, fmt.Errorf("%q must match as string, ip, or integer", field))
		}
	}
	return
}
func validatePoolLicenseType(value interface{}, field string) (ws []string, errors []error) {
	var values []string
	switch value.(type) {
	case *schema.Set:
		values = setToStringSlice(value.(*schema.Set))
	case []string:
		values = value.([]string)
	case *[]string:
		values = *(value.(*[]string))
	case string:
		values = []string{value.(string)}
	default:
		errors = append(errors, fmt.Errorf("Unknown type %v in validatePoolLicenseType", reflect.TypeOf(value)))
	}
	for _, v := range values {
		match, _ := regexp.MatchString("(?mi)^Utility$|^regkey$", v)
		if !match {
			errors = append(errors, fmt.Errorf("%q must match as Utility (or) Regkey", field))
		}
	}
	return
}
func validateAssignmentType(value interface{}, field string) (ws []string, errors []error) {
	var values []string
	switch value.(type) {
	case *schema.Set:
		values = setToStringSlice(value.(*schema.Set))
	case []string:
		values = value.([]string)
	case *[]string:
		values = *(value.(*[]string))
	case string:
		values = []string{value.(string)}
	default:
		errors = append(errors, fmt.Errorf("Unknown type %v in validatePoolLicenseType", reflect.TypeOf(value)))
	}
	for _, v := range values {
		match, _ := regexp.MatchString("(?mi)^MANAGED$|^UNMANAGED$|^UNREACHABLE$", v)
		if !match {
			errors = append(errors, fmt.Errorf("%q must match as MANAGED/UNMANAGED/UNREACHABLE", field))
		}
	}
	return
}

func getDeviceUri(str string) ([]string, error) {
	re := regexp.MustCompile(`^(?:(?:(https?|s?ftp):)\/\/)([^:\/\s]+)(?::(\d*))?`)
	if len(re.FindStringSubmatch(str)) > 0 {
		return re.FindStringSubmatch(str), nil
	}
	return []string{}, nil
}
