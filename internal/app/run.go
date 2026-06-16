package app

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Hero-exe/mdcheck/internal/config"
	"github.com/Hero-exe/mdcheck/internal/markdown"
	"github.com/Hero-exe/mdcheck/internal/report"
	"github.com/Hero-exe/mdcheck/internal/rules"
)

type ExitCodeError struct {
	Code int
}

func (e ExitCodeError) Error() string {
	return fmt.Sprintf("exit code %d", e.Code)
}

func Run(ctx context.Context, args []string, stdout io.Writer, stderr io.Writer) error {
	fs := flag.NewFlagSet("mdcheck", flag.ContinueOnError)
	fs.SetOutput(stderr)

	format := fs.String("format", "text", "output format: text or json")
	configPath := fs.String("config", "mdcheck.yaml", "config file path")
	showVersion := fs.Bool("version", false, "show version")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if *showVersion {
		fmt.Fprintln(stdout, "mdcheck 0.1.0")
		return nil
	}
	if fs.NArg() == 0 {
		return fmt.Errorf("usage: mdcheck [--format text|json] [--config mdcheck.yaml] <file-or-directory>...")
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		return err
	}

	files, err := collectMarkdownFiles(fs.Args(), cfg.Ignore)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		fmt.Fprintln(stderr, "no markdown files found")
		return nil
	}

	ruleSet := rules.DefaultRules()
	var results []report.FileResult
	hasError := false

	for _, file := range files {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		doc, err := markdown.ParseFile(file)
		if err != nil {
			return err
		}

		var findings []rules.Finding
		checkCtx := rules.Context{Config: cfg}
		for _, rule := range ruleSet {
			severity := cfg.SeverityFor(rule.Name())
			if severity == rules.SeverityOff {
				continue
			}
			for _, finding := range rule.Check(checkCtx, doc) {
				if finding.Severity == "" {
					finding.Severity = severity
				}
				findings = append(findings, finding)
				if finding.Severity == rules.SeverityError {
					hasError = true
				}
			}
		}

		results = append(results, report.FileResult{File: file, Findings: findings})
	}

	switch *format {
	case "text":
		report.WriteText(stdout, results)
	case "json":
		enc := json.NewEncoder(stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(results); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown format %q: use text or json", *format)
	}

	if hasError {
		return ExitCodeError{Code: 1}
	}
	return nil
}

func collectMarkdownFiles(paths []string, ignore []string) ([]string, error) {
	var files []string
	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			return nil, err
		}
		if !info.IsDir() {
			if isMarkdown(path) && !isIgnored(path, ignore) {
				files = append(files, path)
			}
			continue
		}
		err = filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				name := d.Name()
				if name == ".git" || name == "vendor" || name == "node_modules" || isIgnored(p, ignore) {
					return filepath.SkipDir
				}
				return nil
			}
			if isMarkdown(p) && !isIgnored(p, ignore) {
				files = append(files, p)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return files, nil
}

func isMarkdown(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".md" || ext == ".markdown"
}

func isIgnored(path string, patterns []string) bool {
	cleanPath := filepath.ToSlash(filepath.Clean(path))
	for _, pattern := range patterns {
		pattern = filepath.ToSlash(strings.TrimSpace(pattern))
		if pattern == "" {
			continue
		}
		pattern = strings.TrimSuffix(pattern, "/")
		if cleanPath == pattern || strings.HasPrefix(cleanPath, pattern+"/") || strings.Contains(cleanPath, "/"+pattern+"/") {
			return true
		}
	}
	return false
}
