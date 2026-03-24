package main

import (
	"fmt"
	"strings"

	"github.com/googleapis/api-linter/v2/lint"
)

// ruleNameToBufPluginID converts an api-linter rule name like
// "core::0203::field-behavior-required" to a bufplugin rule ID like
// "AIP_0203_FIELD_BEHAVIOR_REQUIRED". The group prefix (core::, client-libraries::)
// is dropped in favor of the AIP_ prefix.
func ruleNameToBufPluginID(name lint.RuleName) string {
	parts := strings.SplitN(string(name), "::", 3)
	if len(parts) < 3 {
		s := strings.ReplaceAll(string(name), "::", "_")
		s = strings.ReplaceAll(s, "-", "_")
		return "AIP_" + strings.ToUpper(s)
	}
	// Drop the group, keep AIP number and rule name
	rule := strings.ReplaceAll(parts[2], "-", "_")
	return "AIP_" + parts[1] + "_" + strings.ToUpper(rule)
}

// ruleNameToURL converts an api-linter rule name like
// "core::0203::field-behavior-required" to a URL like
// "https://linter.aip.dev/203/field-behavior-required".
func ruleNameToURL(name lint.RuleName) string {
	s := string(name)
	parts := strings.SplitN(s, "::", 3)
	if len(parts) >= 2 {
		aip := strings.TrimLeft(parts[1], "0")
		return fmt.Sprintf("https://linter.aip.dev/%s/%s", aip, parts[2])
	}
	return ""
}

// aipCategoryID extracts the AIP category from a rule name.
// For "core::0203::field-behavior-required", returns "AIP_0203".
func aipCategoryID(name lint.RuleName) string {
	parts := strings.SplitN(string(name), "::", 3)
	if len(parts) >= 2 {
		return "AIP_" + parts[1]
	}
	return "AIP_UNKNOWN"
}

// aipCategoryPurpose returns a human-readable purpose for an AIP category.
func aipCategoryPurpose(aipNumber string) string {
	return "Rules from AIP-" + strings.TrimLeft(aipNumber, "0") + "."
}
