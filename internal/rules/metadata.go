package rules

import (
	"fmt"
	"strings"

	"github.com/Hero-exe/mdcheck/internal/markdown"
)

type MetadataRule struct{}

func (MetadataRule) Name() string {
	return "metadata"
}

func (r MetadataRule) Check(ctx Context, doc markdown.Document) []Finding {
	required := []string{"title", "description", "date", "tags"}
	if cfg, ok := ctx.Config.(interface{ RequiredMetadata() []string }); ok {
		required = cfg.RequiredMetadata()
	}

	var findings []Finding
	if len(doc.FrontMatter) == 0 {
		return []Finding{{
			Rule:    r.Name(),
			Message: "missing YAML front matter",
		}}
	}

	for _, key := range required {
		value, ok := doc.FrontMatter[key]
		if !ok || strings.TrimSpace(value) == "" {
			findings = append(findings, Finding{
				Rule:    r.Name(),
				Message: fmt.Sprintf("missing metadata field: %s", key),
			})
		}
	}
	return findings
}
