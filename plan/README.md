# mdcheck implementation plan

This document is a handoff plan for future agents working on `mdcheck`.

`mdcheck` is a Go CLI that checks Markdown articles for blog publishing quality:
heading structure, empty links, image paths, character count, and front matter metadata.

## Current Status

- Initial CLI is implemented in Go with only the standard library.
- Entry point: `cmd/mdcheck/main.go`
- Main orchestration: `internal/app/run.go`
- Markdown parser: `internal/markdown`
- Rules: `internal/rules`
- Report output: `internal/report`
- Config loading: `internal/config`
- Tests exist for Markdown parsing and core rules.
- Module path is `github.com/Hero-exe/mdcheck`.

Run checks with:

```sh
go test ./...
go run ./cmd/mdcheck README.md
go run ./cmd/mdcheck --format json README.md
```

## Public Repository Notes

- Do not commit `.cache/`, `outputs/`, `.DS_Store`, or built binaries.
- No secrets or private config should be added.
- A `LICENSE` file is not added yet. Add one before presenting this as reusable open source.

## Design Principles

- Keep rules independent and easy to add.
- Prefer clear findings over clever parsing.
- Avoid external dependencies until they solve a real parser/config problem.
- Keep CLI output stable enough for CI usage.
- Text output is for humans; JSON output is for automation.

## Rule Interface

New checks should be implemented as rules under `internal/rules`.

Each rule should:

- implement `Name() string`
- implement `Check(ctx Context, doc markdown.Document) []Finding`
- return precise line numbers when possible
- avoid exiting or printing directly
- be added to `DefaultRules()` in `internal/rules/rule.go`
- include focused tests

## Priority Tasks

### 1. Add License

Goal:
Make the public repository clearly reusable.

Suggested work:

- Add `LICENSE`.
- MIT is a reasonable default for a small CLI unless the owner wants another license.
- Mention the license in `README.md`.

Acceptance criteria:

- `LICENSE` exists.
- README includes a short license section.

### 2. Improve README for Public Users

Goal:
Make the project understandable from GitHub.

Suggested work:

- Add install instructions:
  - `go install github.com/Hero-exe/mdcheck/cmd/mdcheck@latest`
  - local development commands
- Add example output.
- Add config file example.
- Add exit code behavior.
- Add a short roadmap section.

Acceptance criteria:

- A new user can install, run, and configure the CLI from the README alone.

### 3. Add GitHub Actions CI

Goal:
Run tests automatically on pull requests and pushes.

Suggested work:

- Add `.github/workflows/ci.yml`.
- Run `go test ./...`.
- Use a stable Go version.
- Optionally run `gofmt` check.

Acceptance criteria:

- CI runs on `push` and `pull_request`.
- CI fails if tests fail.
- CI fails if formatting is not gofmt-compliant.

### 4. Replace Ad Hoc Config Parsing

Goal:
Support real YAML config more reliably.

Current limitation:
`internal/config` uses a small line-based parser. It handles the current examples but is not a full YAML parser.

Suggested work:

- Add `gopkg.in/yaml.v3` or another maintained YAML parser.
- Preserve the current config shape.
- Keep defaults when no config exists.
- Add config tests for:
  - missing file
  - rule severity override
  - metadata required override
  - word count min/max
  - ignore patterns

Acceptance criteria:

- Existing config examples still work.
- More realistic YAML, including comments and lists, is handled correctly.

### 5. Improve Markdown Parsing

Goal:
Reduce false positives and support more Markdown syntax.

Current limitation:
`internal/markdown` uses regular expressions and simple line scanning.

Suggested work:

- Consider `github.com/yuin/goldmark` for Markdown AST parsing.
- Preserve current `markdown.Document` as the internal contract if possible.
- Support:
  - reference links
  - autolinks
  - Markdown links with titles
  - images with titles
  - nested emphasis without confusing link detection

Acceptance criteria:

- Existing tests pass.
- New tests cover reference-style links and image titles.
- Fenced code blocks are still ignored.

### 6. Add Image Alt Text Rule

Goal:
Warn when blog images have no alt text.

Suggested work:

- Add `image_alt` rule.
- Detect `![](path.png)` and `![   ](path.png)`.
- Make default severity `warn`.
- Add config support via existing severity mechanism.

Acceptance criteria:

- Empty alt text produces a finding.
- Non-empty alt text does not.
- Rule can be disabled with `image_alt: off`.

### 7. Add Title Length Rule

Goal:
Help blog posts keep metadata titles readable.

Suggested work:

- Add `title_length` rule.
- Check front matter `title`.
- Add config:
  - `title_length.min`
  - `title_length.max`
- Default suggestion: max 60 characters.

Acceptance criteria:

- Missing title remains covered by metadata rule.
- Too-long title produces a finding.
- Config can change max length.

### 8. Add Internal Link Rule

Goal:
Catch broken local Markdown links.

Suggested work:

- Add `internal_link` rule.
- Ignore remote URLs.
- Resolve relative paths from the current document directory.
- Support anchor-only links like `#section`.
- For file links with anchors, check file existence first.

Acceptance criteria:

- Missing local file links are reported.
- Remote links are ignored.
- Existing files are accepted.

### 9. Add Better Exit Code Tests

Goal:
Make CI behavior stable.

Suggested work:

- Add tests around `internal/app.Run`.
- Verify:
  - no error findings returns nil
  - error findings return `ExitCodeError{Code: 1}`
  - invalid format returns a normal error
  - JSON output is valid

Acceptance criteria:

- App-level behavior is tested without invoking `os.Exit`.

### 10. Prepare Releases

Goal:
Make the CLI easy to install as a binary.

Suggested work:

- Add GoReleaser config or a simple GitHub Actions release workflow.
- Build binaries for:
  - macOS arm64
  - macOS amd64
  - Linux amd64
  - Linux arm64
- Add checksums.

Acceptance criteria:

- Tagged release builds downloadable binaries.
- README explains install options.

## Nice-To-Have Tasks

- Add `--quiet`.
- Add `--fail-on warn`.
- Add `--no-color` and colored text output.
- Add `--ignore` CLI flag.
- Add `--rule` CLI flag for running a subset of rules.
- Add SARIF output for GitHub code scanning.
- Add `mdcheck init` to generate `mdcheck.yaml`.
- Add Japanese README section or `README.ja.md`.

## Suggested First Issue Set

For another agent, a good first batch is:

1. Add `LICENSE`.
2. Expand `README.md`.
3. Add GitHub Actions CI.
4. Add `image_alt` rule.

These are independent, low-risk, and improve the public repository immediately.
