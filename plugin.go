package main

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"buf.build/go/bufplugin/check"
	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/rules"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type contextKey struct{}

type lintResults struct {
	problems map[string][]problemWithFile
}

type problemWithFile struct {
	filePath string
	problem  lint.Problem
}

func buildSpec() *check.Spec {
	registry := lint.NewRuleRegistry()
	if err := rules.Add(registry); err != nil {
		panic(fmt.Sprintf("failed to register api-linter rules: %v", err))
	}

	// Build the ID mapping and collect categories.
	idToRuleName := make(map[string]lint.RuleName)
	categories := make(map[string]string) // categoryID -> aipNumber
	var ruleSpecs []*check.RuleSpec

	for name := range registry {
		bufID := ruleNameToBufPluginID(name)
		idToRuleName[bufID] = name

		catID := aipCategoryID(name)
		parts := strings.SplitN(string(name), "::", 3)
		if len(parts) >= 2 {
			categories[catID] = parts[1]
		}

		isDefault := strings.HasPrefix(string(name), "core::")

		ruleID := bufID // capture for closure
		ruleSpecs = append(ruleSpecs, &check.RuleSpec{
			ID:          ruleID,
			CategoryIDs: []string{catID},
			Default:     isDefault,
			Purpose:     fmt.Sprintf("Enforces %s.", name),
			Type:        check.RuleTypeLint,
			Handler: check.RuleHandlerFunc(func(ctx context.Context, responseWriter check.ResponseWriter, request check.Request) error {
				results, ok := ctx.Value(contextKey{}).(*lintResults)
				if !ok {
					return fmt.Errorf("lint results not found in context")
				}
				for _, p := range results.problems[ruleID] {
					opts := []check.AddAnnotationOption{
						check.WithMessage(p.problem.Message),
					}
					if p.problem.Location != nil {
						sourcePath := make(protoreflect.SourcePath, len(p.problem.Location.GetPath()))
						for i, v := range p.problem.Location.GetPath() {
							sourcePath[i] = v
						}
						opts = append(opts, check.WithFileNameAndSourcePath(p.filePath, sourcePath))
					} else if p.problem.Descriptor != nil {
						opts = append(opts, check.WithDescriptor(p.problem.Descriptor))
					}
					responseWriter.AddAnnotation(opts...)
				}
				return nil
			}),
		})
	}

	// Sort rule specs for deterministic output.
	sort.Slice(ruleSpecs, func(i, j int) bool {
		return ruleSpecs[i].ID < ruleSpecs[j].ID
	})

	// Build category specs.
	var categorySpecs []*check.CategorySpec
	for catID, aipNumber := range categories {
		categorySpecs = append(categorySpecs, &check.CategorySpec{
			ID:      catID,
			Purpose: aipCategoryPurpose(aipNumber),
		})
	}
	sort.Slice(categorySpecs, func(i, j int) bool {
		return categorySpecs[i].ID < categorySpecs[j].ID
	})

	return &check.Spec{
		Rules:      ruleSpecs,
		Categories: categorySpecs,
		Before: func(ctx context.Context, request check.Request) (context.Context, check.Request, error) {
			protoFiles := make([]protoreflect.FileDescriptor, 0, len(request.FileDescriptors()))
			for _, fd := range request.FileDescriptors() {
				if fd.IsImport() {
					continue
				}
				protoFiles = append(protoFiles, fd.ProtoreflectFileDescriptor())
			}

			configs := lint.Configs{
				{EnabledRules: []string{"all"}},
			}
			l := lint.New(registry, configs)

			responses, err := l.LintProtos(protoFiles...)
			if err != nil {
				return ctx, request, fmt.Errorf("api-linter failed: %w", err)
			}

			results := &lintResults{
				problems: make(map[string][]problemWithFile),
			}
			for _, resp := range responses {
				for _, prob := range resp.Problems {
					bufID := ruleNameToBufPluginID(prob.RuleID)
					results.problems[bufID] = append(results.problems[bufID], problemWithFile{
						filePath: resp.FilePath,
						problem:  prob,
					})
				}
			}

			ctx = context.WithValue(ctx, contextKey{}, results)
			return ctx, request, nil
		},
	}
}
