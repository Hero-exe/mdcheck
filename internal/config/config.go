package config

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/Hero-exe/mdcheck/internal/rules"
)

type Config struct {
	Rules    map[string]rules.Severity
	Metadata MetadataConfig
	Word     WordCountConfig
	Ignore   []string
}

type MetadataConfig struct {
	Required []string
}

type WordCountConfig struct {
	Min int
	Max int
}

func Default() Config {
	return Config{
		Rules: map[string]rules.Severity{
			"heading_structure": rules.SeverityError,
			"empty_link":        rules.SeverityError,
			"image_path":        rules.SeverityError,
			"word_count":        rules.SeverityWarn,
			"metadata":          rules.SeverityWarn,
		},
		Metadata: MetadataConfig{
			Required: []string{"title", "description", "date", "tags"},
		},
		Word: WordCountConfig{
			Min: 0,
			Max: 0,
		},
		Ignore: []string{".git", "vendor", "node_modules"},
	}
}

func Load(path string) (Config, error) {
	cfg := Default()
	if path == "" {
		return cfg, nil
	}
	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return cfg, nil
	}
	if err != nil {
		return cfg, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var section string
	var subsection string
	metadataRequiredSeen := false

	for scanner.Scan() {
		raw := scanner.Text()
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if !strings.HasPrefix(raw, " ") && strings.HasSuffix(line, ":") {
			section = strings.TrimSuffix(line, ":")
			subsection = ""
			continue
		}

		if strings.HasPrefix(raw, "  ") && strings.HasSuffix(line, ":") {
			subsection = strings.TrimSuffix(line, ":")
			continue
		}

		switch section {
		case "rules":
			key, value, ok := splitKeyValue(line)
			if ok {
				cfg.Rules[key] = rules.ParseSeverity(value)
			}
		case "metadata":
			if subsection == "required" {
				if item, ok := parseListItem(line); ok {
					if !metadataRequiredSeen {
						cfg.Metadata.Required = nil
						metadataRequiredSeen = true
					}
					cfg.Metadata.Required = appendIfMissing(cfg.Metadata.Required, item)
				}
			}
		case "word_count":
			key, value, ok := splitKeyValue(line)
			if !ok {
				continue
			}
			n, err := strconv.Atoi(value)
			if err != nil {
				continue
			}
			if key == "min" {
				cfg.Word.Min = n
			}
			if key == "max" {
				cfg.Word.Max = n
			}
		case "ignore":
			if item, ok := parseListItem(line); ok {
				cfg.Ignore = appendIfMissing(cfg.Ignore, item)
			}
		}
	}

	return cfg, scanner.Err()
}

func (c Config) SeverityFor(name string) rules.Severity {
	if severity, ok := c.Rules[name]; ok {
		return severity
	}
	return rules.SeverityWarn
}

func splitKeyValue(line string) (string, string, bool) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	return strings.TrimSpace(parts[0]), strings.Trim(strings.TrimSpace(parts[1]), `"'`), true
}

func parseListItem(line string) (string, bool) {
	if !strings.HasPrefix(line, "- ") {
		return "", false
	}
	return strings.Trim(strings.TrimSpace(strings.TrimPrefix(line, "- ")), `"'`), true
}

func appendIfMissing(values []string, value string) []string {
	for _, existing := range values {
		if existing == value {
			return values
		}
	}
	return append(values, value)
}
