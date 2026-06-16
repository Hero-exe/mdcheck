package rules

import "github.com/Hero-exe/mdcheck/internal/markdown"

type EmptyLinkRule struct{}

func (EmptyLinkRule) Name() string {
	return "empty_link"
}

func (r EmptyLinkRule) Check(ctx Context, doc markdown.Document) []Finding {
	var findings []Finding
	for _, link := range doc.Links {
		if link.Text == "" {
			findings = append(findings, Finding{
				Rule:    r.Name(),
				Message: "empty link text",
				Line:    link.Line,
			})
		}
		if link.URL == "" {
			findings = append(findings, Finding{
				Rule:    r.Name(),
				Message: "empty link URL",
				Line:    link.Line,
			})
		}
	}
	return findings
}
