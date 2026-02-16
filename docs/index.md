---
title: sixtine-coverage
description: Interactive Go coverage dashboard from cover.out
---

# sixtine-coverage

From `cover.out` to a beautiful, shareable coverage dashboard in one command.

`sixtine-coverage` turns raw Go coverage into a visual report your team can read in seconds.

It works great for local development, CI artifacts, and review workflows.

## Live demo report

- Open the bundled interactive demo: [demo-report.html](demo-report.html)

## Why use it

- Single self-contained HTML report (great for CI artifacts)
- Folder and file drill-down with source highlighting
- Optional HTTP route coverage view
- Docker mode for a persistent team dashboard
- Works with mixed coverage path styles (module prefix optional)

## What you can show in demos

- Overall statement coverage and drill-down by package
- Covered vs non-covered blocks at source line level
- Route-level API coverage grouped by service
- Uncovered endpoints and failure responses (`4xx`, `5xx`, `-1`)

## Screenshots

### Dashboard

![Dashboard](dashboard-overview.png)

### File and folder drill-down

![Files list](files-list.png)
![Explorer and highlighted source](explorer-tree-main-go.png)
![Files sorted by coverage](files-list-sorted.png)

### Routes view

![Route coverage summary](route-coverage.png)

## 2-minute quickstart

```bash
go test -coverprofile=cover.out ./...
go install github.com/sixtineai/sixtine-coverage@latest
sixtine-coverage -coverage cover.out -codebase . -out coverage.html
```

Open `coverage.html` in your browser.

## Route coverage quick add-on

```bash
sixtine-coverage \
	-coverage cover.out \
	-codebase . \
	-route-coverage route-coverage.rcov \
	-out coverage.html
```

## Demo package

Want to test now without your own project?

Use [sample-data/README.md](../sample-data/README.md).

## Links

- Quickstart: [quickstart.md](quickstart.md)
- Docker mode: [docker.md](docker.md)
- Route coverage format: [route-coverage.md](route-coverage.md)
- FAQ: [faq.md](faq.md)

---

If this saves you time, star the repo and share it with your Go team.
