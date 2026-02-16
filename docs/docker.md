# Docker mode (live team dashboard)

Prefer a quick no-setup tour first? Open [demo-report.html](demo-report.html).

Run the dashboard with mounted coverage and source files.

```bash
docker run -d -p 8080:80 \
  -v $(pwd)/cover.out:/data/coverage/cover.out:ro \
  -v $(pwd):/data/code:ro \
  -e SIXTINEAI_COVERAGE_URL=/data/coverage/cover.out \
  -e SIXTINEAI_SOURCE_ROOT=/data/code/ \
  ghcr.io/sixtineai/sixtine-coverage:latest
```

Open http://localhost:8080.

## UI preview

![Dashboard](dashboard-overview.png)

## With route coverage

Mount route file and set:

`SIXTINEAI_ROUTE_COVERAGE_URL=/data/coverage/route-coverage.rcov`

Example:

```bash
docker run -d -p 8080:80 \
  -v $(pwd)/cover.out:/data/coverage/cover.out:ro \
  -v $(pwd)/route-coverage.rcov:/data/coverage/route-coverage.rcov:ro \
  -v $(pwd):/data/code:ro \
  -e SIXTINEAI_COVERAGE_URL=/data/coverage/cover.out \
  -e SIXTINEAI_ROUTE_COVERAGE_URL=/data/coverage/route-coverage.rcov \
  -e SIXTINEAI_SOURCE_ROOT=/data/code/ \
  ghcr.io/sixtineai/sixtine-coverage:latest
```

## Best practices

- Mount coverage files as read-only
- Keep source root aligned with paths in `cover.out`
- Regenerate `cover.out` per commit or per CI run
- Use tagged images in CI for reproducibility

## Team workflow

1. CI generates `cover.out` and `route-coverage.rcov`.
2. Files are mounted in a shared environment.
3. Team uses one dashboard URL for every build.
