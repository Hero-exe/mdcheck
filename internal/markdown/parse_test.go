package markdown

import "testing"

func TestParseFrontMatterAndLineNumbers(t *testing.T) {
	doc := Parse("post.md", `---
title: Hello
description: Intro
tags:
  - go
---
# Title

## Section
[empty]()
![missing](image.png)
`)

	if doc.FrontMatter["title"] != "Hello" {
		t.Fatalf("expected title metadata, got %q", doc.FrontMatter["title"])
	}
	if doc.FrontMatter["tags"] != "go" {
		t.Fatalf("expected list metadata, got %q", doc.FrontMatter["tags"])
	}
	if len(doc.Headings) != 2 {
		t.Fatalf("expected 2 headings, got %d", len(doc.Headings))
	}
	if doc.Headings[0].Line != 7 {
		t.Fatalf("expected first heading on line 7, got %d", doc.Headings[0].Line)
	}
	if len(doc.Links) != 1 {
		t.Fatalf("expected 1 link, got %d", len(doc.Links))
	}
	if len(doc.Images) != 1 {
		t.Fatalf("expected 1 image, got %d", len(doc.Images))
	}
}

func TestParseIgnoresFencedCode(t *testing.T) {
	doc := Parse("post.md", "# Title\n\n```md\n# Not a heading\n[]()\n```\n")

	if len(doc.Headings) != 1 {
		t.Fatalf("expected 1 heading, got %d", len(doc.Headings))
	}
	if len(doc.Links) != 0 {
		t.Fatalf("expected no links from fenced code, got %d", len(doc.Links))
	}
}
