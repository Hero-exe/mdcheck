package rules

import (
	"fmt"
	"regexp"
	"unicode/utf8"

	"github.com/Hero-exe/mdcheck/internal/markdown"
)

type WordCountRule struct{}

var markdownSyntaxPattern = regexp.MustCompile(`(?m)^[#>*\-\s]+|[*_` + "`" + `\[\]()]`)

func (WordCountRule) Name() string {
	return "word_count"
}

func (r WordCountRule) Check(ctx Context, doc markdown.Document) []Finding {
	count := utf8.RuneCountInString(markdownSyntaxPattern.ReplaceAllString(doc.Body, ""))
	findings := []Finding{{
		Rule:     r.Name(),
		Severity: SeverityInfo,
		Message:  fmt.Sprintf("character count: %d", count),
	}}

	if cfg, ok := ctx.Config.(interface{ GetWordCountRange() (int, int) }); ok {
		min, max := cfg.GetWordCountRange()
		if min > 0 && count < min {
			findings = append(findings, Finding{
				Rule:    r.Name(),
				Message: fmt.Sprintf("character count is below minimum: %d < %d", count, min),
			})
		}
		if max > 0 && count > max {
			findings = append(findings, Finding{
				Rule:    r.Name(),
				Message: fmt.Sprintf("character count exceeds maximum: %d > %d", count, max),
			})
		}
	}

	return findings
}
