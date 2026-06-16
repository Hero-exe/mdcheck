package rules

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/Hero-exe/mdcheck/internal/markdown"
)

type ImagePathRule struct{}

func (ImagePathRule) Name() string {
	return "image_path"
}

func (r ImagePathRule) Check(ctx Context, doc markdown.Document) []Finding {
	var findings []Finding
	for _, image := range doc.Images {
		if image.Path == "" {
			findings = append(findings, Finding{
				Rule:    r.Name(),
				Message: "empty image path",
				Line:    image.Line,
			})
			continue
		}
		if isRemote(image.Path) || strings.HasPrefix(image.Path, "#") {
			continue
		}

		path := image.Path
		if hash := strings.Index(path, "#"); hash >= 0 {
			path = path[:hash]
		}
		if query := strings.Index(path, "?"); query >= 0 {
			path = path[:query]
		}
		if filepath.IsAbs(path) {
			if _, err := os.Stat(path); err != nil {
				findings = append(findings, missingImageFinding(r.Name(), image.Line, image.Path))
			}
			continue
		}

		fullPath := filepath.Join(doc.Dir, path)
		if _, err := os.Stat(fullPath); err != nil {
			findings = append(findings, missingImageFinding(r.Name(), image.Line, image.Path))
		}
	}
	return findings
}

func isRemote(path string) bool {
	u, err := url.Parse(path)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https")
}

func missingImageFinding(rule string, line int, path string) Finding {
	return Finding{
		Rule:    rule,
		Message: fmt.Sprintf("image path does not exist: %s", path),
		Line:    line,
	}
}
