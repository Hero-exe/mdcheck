package markdown

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	headingPattern = regexp.MustCompile(`^(#{1,6})(?:\s+|$)(.*?)\s*#*\s*$`)
	imagePattern   = regexp.MustCompile(`!\[([^\]]*)\]\(([^)]*)\)`)
	linkPattern    = regexp.MustCompile(`(^|[^!])\[([^\]]*)\]\(([^)]*)\)`)
)

func ParseFile(path string) (Document, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Document{}, err
	}
	return Parse(path, string(data)), nil
}

func Parse(path string, source string) Document {
	frontMatter, body, bodyLine := splitFrontMatter(source)
	doc := Document{
		Path:        path,
		Dir:         filepath.Dir(path),
		Source:      source,
		Body:        body,
		BodyLine:    bodyLine,
		FrontMatter: frontMatter,
	}

	lines := strings.Split(body, "\n")
	inFence := false
	for i, line := range lines {
		lineNo := i + bodyLine
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}

		if match := headingPattern.FindStringSubmatch(line); match != nil {
			doc.Headings = append(doc.Headings, Heading{
				Level: len(match[1]),
				Text:  strings.TrimSpace(match[2]),
				Line:  lineNo,
			})
		}

		for _, match := range imagePattern.FindAllStringSubmatch(line, -1) {
			doc.Images = append(doc.Images, Image{
				Alt:  strings.TrimSpace(match[1]),
				Path: strings.TrimSpace(match[2]),
				Line: lineNo,
			})
		}

		for _, match := range linkPattern.FindAllStringSubmatch(line, -1) {
			doc.Links = append(doc.Links, Link{
				Text: strings.TrimSpace(match[2]),
				URL:  strings.TrimSpace(match[3]),
				Line: lineNo,
			})
		}
	}

	return doc
}

func splitFrontMatter(source string) (map[string]string, string, int) {
	frontMatter := map[string]string{}
	normalized := strings.ReplaceAll(source, "\r\n", "\n")
	if !strings.HasPrefix(normalized, "---\n") {
		return frontMatter, source, 1
	}

	rest := strings.TrimPrefix(normalized, "---\n")
	end := strings.Index(rest, "\n---\n")
	if end < 0 {
		return frontMatter, source, 1
	}

	block := rest[:end]
	body := rest[end+len("\n---\n"):]
	bodyLine := frontMatterLineCount(block) + 3
	currentKey := ""
	for _, rawLine := range strings.Split(block, "\n") {
		line := strings.TrimSpace(rawLine)
		if strings.HasPrefix(line, "- ") && currentKey != "" {
			item := strings.Trim(strings.TrimSpace(strings.TrimPrefix(line, "- ")), `"'`)
			if item != "" {
				if frontMatter[currentKey] == "" {
					frontMatter[currentKey] = item
				} else {
					frontMatter[currentKey] += "," + item
				}
			}
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
		if key != "" {
			frontMatter[key] = value
			currentKey = key
		}
	}

	return frontMatter, body, bodyLine
}

func frontMatterLineCount(block string) int {
	if block == "" {
		return 0
	}
	return strings.Count(block, "\n") + 1
}
