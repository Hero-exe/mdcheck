package rules

import (
	"strings"
	"testing"

	"github.com/Hero-exe/mdcheck/internal/markdown"
)

func TestHeadingStructureRule(t *testing.T) {
	doc := markdown.Parse("post.md", "# Title\n\n### Skipped\n\n# Again\n")
	findings := HeadingStructureRule{}.Check(Context{}, doc)

	if len(findings) != 2 {
		t.Fatalf("expected 2 findings, got %d: %#v", len(findings), findings)
	}
}

func TestEmptyLinkRule(t *testing.T) {
	doc := markdown.Parse("post.md", "# Title\n\n[]()\n")
	findings := EmptyLinkRule{}.Check(Context{}, doc)

	if len(findings) != 2 {
		t.Fatalf("expected 2 findings, got %d", len(findings))
	}
}

func TestMetadataRule(t *testing.T) {
	doc := markdown.Parse("post.md", "---\ntitle: Hello\n---\n# Title\n")
	findings := MetadataRule{}.Check(Context{}, doc)

	if len(findings) == 0 {
		t.Fatal("expected metadata findings")
	}
	if !strings.Contains(findings[0].Message, "description") {
		t.Fatalf("expected missing description first, got %q", findings[0].Message)
	}
}
