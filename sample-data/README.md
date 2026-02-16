# Sample data pack (instant realistic demo)

Use this pack to test `sixtine-coverage` without any external project.

This demo now includes:

- 27 Go files
- 8 package directories (well above the requested 7)
- Rich route coverage across multiple services
- Mixed covered / uncovered blocks for better visual contrast

## Included files

- [sample-data/demo-cover.out](demo-cover.out)
- [sample-data/demo-route-coverage.rcov](demo-route-coverage.rcov)
- [sample-data/demo-source](demo-source)
- [sample-data/docker-compose.demo.yml](docker-compose.demo.yml)

## Package map

- `cmd/api`
- `cmd/worker`
- `internal/httpapi`
- `internal/service`
- `internal/store`
- `pkg/model`
- `pkg/validator`
- `pkg/utils`
- `pkg/report`

## Local CLI demo

From repository root:

```bash
go run . \
  -coverage ./sample-data/demo-cover.out \
  -codebase ./sample-data/demo-source \
  -route-coverage ./sample-data/demo-route-coverage.rcov \
  -out ./sample-data/demo-report.html
```

Open `./sample-data/demo-report.html`.

## Docker demo

```bash
cd sample-data
docker compose -f docker-compose.demo.yml up -d
```

Open http://localhost:8088.

## Troubleshooting: module prefix detection

Current behavior auto-resolves file paths even when no common prefix is detected.

## Screenshot checklist

Capture these screens in order:

1. Global dashboard
2. Folder/package tree expansion (`internal/*`, `pkg/*`)
3. `cmd/api/main.go` line highlights
4. `internal/service/health.go` line highlights
5. A mostly uncovered file (`pkg/report/route.go`)
6. Route coverage panel with service grouping
