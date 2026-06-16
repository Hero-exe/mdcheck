package rules

import "github.com/Hero-exe/mdcheck/internal/markdown"

type Severity string

const (
	SeverityError Severity = "error"
	SeverityWarn  Severity = "warn"
	SeverityInfo  Severity = "info"
	SeverityOff   Severity = "off"
)

type Context struct {
	Config interface{}
}

type Rule interface {
	Name() string
	Check(ctx Context, doc markdown.Document) []Finding
}

type Finding struct {
	Rule     string   `json:"rule"`
	Severity Severity `json:"severity"`
	Message  string   `json:"message"`
	Line     int      `json:"line,omitempty"`
}

func ParseSeverity(value string) Severity {
	switch Severity(value) {
	case SeverityError, SeverityWarn, SeverityInfo, SeverityOff:
		return Severity(value)
	default:
		return SeverityWarn
	}
}

func DefaultRules() []Rule {
	return []Rule{
		HeadingStructureRule{},
		EmptyLinkRule{},
		ImagePathRule{},
		WordCountRule{},
		MetadataRule{},
	}
}
