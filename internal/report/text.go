package report

import (
	"fmt"
	"io"

	"github.com/Hero-exe/mdcheck/internal/rules"
)

type FileResult struct {
	File     string          `json:"file"`
	Findings []rules.Finding `json:"findings"`
}

func WriteText(w io.Writer, results []FileResult) {
	for i, result := range results {
		if i > 0 {
			fmt.Fprintln(w)
		}
		fmt.Fprintln(w, result.File)
		if len(result.Findings) == 0 {
			fmt.Fprintln(w, "  ok")
			continue
		}
		for _, finding := range result.Findings {
			line := ""
			if finding.Line > 0 {
				line = fmt.Sprintf(":%d", finding.Line)
			}
			fmt.Fprintf(w, "  [%s] %s%s\n", finding.Severity, finding.Rule, line)
			fmt.Fprintf(w, "    %s\n", finding.Message)
		}
	}
}
