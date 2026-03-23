package main

import (
	"testing"

	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/rules"
)

func TestRuleNameToBufPluginID(t *testing.T) {
	tests := []struct {
		input lint.RuleName
		want  string
	}{
		{"core::0203::field-behavior-required", "CORE_0203_FIELD_BEHAVIOR_REQUIRED"},
		{"core::0122::name-suffix", "CORE_0122_NAME_SUFFIX"},
		{"cloud::0203::field-behavior-required", "CLOUD_0203_FIELD_BEHAVIOR_REQUIRED"},
	}
	for _, tt := range tests {
		t.Run(string(tt.input), func(t *testing.T) {
			got := ruleNameToBufPluginID(tt.input)
			if got != tt.want {
				t.Errorf("ruleNameToBufPluginID(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestAIPCategoryID(t *testing.T) {
	tests := []struct {
		input lint.RuleName
		want  string
	}{
		{"core::0203::field-behavior-required", "AIP_0203"},
		{"cloud::0122::name-suffix", "AIP_0122"},
	}
	for _, tt := range tests {
		t.Run(string(tt.input), func(t *testing.T) {
			got := aipCategoryID(tt.input)
			if got != tt.want {
				t.Errorf("aipCategoryID(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestAIPCategoryPurpose(t *testing.T) {
	if got := aipCategoryPurpose("0203"); got != "Rules from AIP-203." {
		t.Errorf("aipCategoryPurpose(\"0203\") = %q, want %q", got, "Rules from AIP-203.")
	}
}

func TestNoBufPluginIDCollisions(t *testing.T) {
	registry := lint.NewRuleRegistry()
	if err := rules.Add(registry); err != nil {
		t.Fatalf("failed to add rules: %v", err)
	}

	seen := make(map[string]lint.RuleName)
	for name := range registry {
		bufID := ruleNameToBufPluginID(name)
		if prev, ok := seen[bufID]; ok {
			t.Errorf("collision: %q and %q both map to %q", prev, name, bufID)
		}
		seen[bufID] = name
	}
}
