package main

import (
	"testing"

	"buf.build/go/bufplugin/check/checktest"
)

func TestSpec(t *testing.T) {
	t.Parallel()
	checktest.SpecTest(t, buildSpec())
}

func TestEnumUnspecified(t *testing.T) {
	t.Parallel()
	checktest.CheckTest{
		Spec: buildSpec(),
		Request: &checktest.RequestSpec{
			Files: &checktest.ProtoFileSpec{
				DirPaths:  []string{"testdata/enum_unspecified"},
				FilePaths: []string{"test.proto"},
			},
			RuleIDs: []string{"AIP_0126_UNSPECIFIED"},
		},
		ExpectedAnnotations: []checktest.ExpectedAnnotation{
			{
				RuleID:  "AIP_0126_UNSPECIFIED",
				Message: `The first enum value should be "STATE_UNSPECIFIED" (https://linter.aip.dev/126/unspecified)`,
				FileLocation: &checktest.ExpectedFileLocation{
					FileName:    "test.proto",
					StartLine:   7,
					StartColumn: 2,
					EndLine:     7,
					EndColumn:   8,
				},
			},
		},
	}.Run(t)
}

func TestFieldLowerSnake(t *testing.T) {
	t.Parallel()
	checktest.CheckTest{
		Spec: buildSpec(),
		Request: &checktest.RequestSpec{
			Files: &checktest.ProtoFileSpec{
				DirPaths:  []string{"testdata/field_names"},
				FilePaths: []string{"test.proto"},
			},
			RuleIDs: []string{"AIP_0140_LOWER_SNAKE"},
		},
		ExpectedAnnotations: []checktest.ExpectedAnnotation{
			{
				RuleID:  "AIP_0140_LOWER_SNAKE",
				Message: "Field `badName` must use lower_snake_case. (https://linter.aip.dev/140/lower-snake)",
				FileLocation: &checktest.ExpectedFileLocation{
					FileName:    "test.proto",
					StartLine:   7,
					StartColumn: 9,
					EndLine:     7,
					EndColumn:   16,
				},
			},
		},
	}.Run(t)
}

// TestDisableRule verifies that when specific RuleIDs are requested,
// only those rules produce annotations. The enum_unspecified proto
// triggers AIP_0126_UNSPECIFIED and several AIP_0191/0192 rules,
// but requesting only AIP_0191_JAVA_PACKAGE should yield just that one.
func TestDisableRule(t *testing.T) {
	t.Parallel()
	checktest.CheckTest{
		Spec: buildSpec(),
		Request: &checktest.RequestSpec{
			Files: &checktest.ProtoFileSpec{
				DirPaths:  []string{"testdata/enum_unspecified"},
				FilePaths: []string{"test.proto"},
			},
			RuleIDs: []string{"AIP_0191_JAVA_PACKAGE"},
		},
		ExpectedAnnotations: []checktest.ExpectedAnnotation{
			{
				RuleID:  "AIP_0191_JAVA_PACKAGE",
				Message: "Proto files must set `option java_package`. (https://linter.aip.dev/191/java-package)",
				FileLocation: &checktest.ExpectedFileLocation{
					FileName:    "test.proto",
					StartLine:   2,
					StartColumn: 0,
					EndLine:     2,
					EndColumn:   16,
				},
			},
		},
	}.Run(t)
}
