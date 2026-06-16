package rules

import (
	"fmt"
	"strings"

	"github.com/Hero-exe/mdcheck/internal/markdown"
)

type HeadingStructureRule struct{}

func (HeadingStructureRule) Name() string {
	return "heading_structure"
}

func (r HeadingStructureRule) Check(ctx Context, doc markdown.Document) []Finding {
	var findings []Finding
	h1Count := 0
	previousLevel := 0

	for _, heading := range doc.Headings {
		if heading.Level == 1 {
			h1Count++
		}
		if strings.TrimSpace(heading.Text) == "" {
			findings = append(findings, Finding{
				Rule:    r.Name(),
				Message: "empty heading text",
				Line:    heading.Line,
			})
		}
		if previousLevel > 0 && heading.Level > previousLevel+1 {
			findings = append(findings, Finding{
				Rule:    r.Name(),
				Message: fmt.Sprintf("heading level jumps from H%d to H%d", previousLevel, heading.Level),
				Line:    heading.Line,
			})
		}
		previousLevel = heading.Level
	}

	if h1Count == 0 {
		findings = append(findings, Finding{
			Rule:    r.Name(),
			Message: "missing H1 heading",
		})
	}
	if h1Count > 1 {
		findings = append(findings, Finding{
			Rule:    r.Name(),
			Message: fmt.Sprintf("multiple H1 headings found: %d", h1Count),
		})
	}

	return findings
}
