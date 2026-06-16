# mdcheck

`mdcheck` is a small Markdown article checker for blogs.

## Checks

- Heading structure: missing or duplicated H1, empty headings, skipped heading levels.
- Empty links: empty link text or empty link URLs.
- Image paths: missing local image files.
- Character count: reports body character count and optional min/max violations.
- Metadata: required YAML front matter fields.

## Usage

```sh
go run ./cmd/mdcheck article.md
go run ./cmd/mdcheck posts/
go run ./cmd/mdcheck --format json posts/
go run ./cmd/mdcheck --config mdcheck.yaml posts/
```

## Configuration

Create `mdcheck.yaml` in the working directory.

```yaml
rules:
  heading_structure: error
  empty_link: error
  image_path: error
  word_count: warn
  metadata: warn

metadata:
  required:
    - title
    - description
    - date
    - tags

word_count:
  min: 800
  max: 5000

ignore:
  - drafts/
  - node_modules/
```

Rule severities are `error`, `warn`, `info`, or `off`.

## Exit Codes

- `0`: no error findings
- `1`: at least one error finding
- `2`: CLI usage or runtime error
