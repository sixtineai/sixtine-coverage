# Quickstart (CLI)

See the bundled interactive demo report: [demo-report.html](demo-report.html)

## Prerequisites

- Go installed
- A coverage profile (`cover.out`)
- Access to the source tree used to generate coverage

## Install

```bash
go install github.com/sixtineai/sixtine-coverage@latest
```

## Generate report

```bash
go test -coverprofile=cover.out ./...
sixtine-coverage -coverage cover.out -codebase . -out coverage.html
```

Open `coverage.html`.

## What the report looks like

![Dashboard](dashboard-overview.png)
![Files list](files-list.png)
![Explorer view](explorer-tree-main-go.png)

## Add route coverage (optional)

```bash
sixtine-coverage \
  -coverage cover.out \
  -codebase . \
  -route-coverage route-coverage.rcov \
  -out coverage.html
```

Route view preview:

![Route coverage overview](route-coverage-overview.png)

## Try with bundled demo data

```bash
go run . \
  -coverage ./sample-data/demo-cover.out \
  -codebase ./sample-data/demo-source \
  -route-coverage ./sample-data/demo-route-coverage.rcov \
  -out ./sample-data/demo-report.html
```

## Path handling behavior

`sixtine-coverage` supports these common `cover.out` path styles:

- `github.com/org/repo/pkg/x.go`
- `repo/pkg/x.go`
- `pkg/x.go`
- mixed relative paths without a single shared ancestor

If no module prefix can be inferred, the CLI now falls back to automatic path resolution.

## CI snippet (GitHub Actions)

```yaml
- name: Run tests with coverage
  run: go test -coverprofile=cover.out ./...

- name: Build coverage report
  run: |
    go install github.com/sixtineai/sixtine-coverage@latest
    sixtine-coverage -coverage cover.out -codebase . -out coverage.html

- name: Upload artifact
  uses: actions/upload-artifact@v4
  with:
    name: coverage-report
    path: coverage.html
```
