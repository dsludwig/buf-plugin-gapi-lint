package main

import (
	"strings"

	"github.com/googleapis/api-linter/v2/lint"
)

// ruleNameToBufPluginID converts an api-linter rule name like
// "core::0203::field-behavior-required" to a bufplugin rule ID like
// "CORE_0203_FIELD_BEHAVIOR_REQUIRED".
func ruleNameToBufPluginID(name lint.RuleName) string {
	s := string(name)
	s = strings.ReplaceAll(s, "::", "_")
	s = strings.ReplaceAll(s, "-", "_")
	return strings.ToUpper(s)
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
