# Deploy to GitHub Pages

This guide gives you two deployment styles:

- `docs/` branch-folder publishing (very simple)
- GitHub Actions Pages deployment (recommended for automation)

The docs folder already includes a runnable demo page at [demo-report.html](demo-report.html), plus screenshots used by the documentation.

## Option A: publish from docs/

1. Copy markdown pages from `plan/github-pages/` to `docs/`.
2. Ensure `docs/index.md` exists.
3. In repository settings:
   - Open **Pages**
   - Source: **Deploy from a branch**
   - Branch: `main`
   - Folder: `/docs`

### Fast local helper

Use [scripts/publish-github-pages.sh](../../scripts/publish-github-pages.sh) to sync docs from the planning folder.

```bash
bash scripts/publish-github-pages.sh
```

## Option B: GitHub Actions deploy

Create `.github/workflows/pages.yml`:

```yaml
name: Deploy Pages

on:
  push:
    branches: [main]
    paths:
      - docs/**
      - .github/workflows/pages.yml

permissions:
  contents: read
  pages: write
  id-token: write

concurrency:
  group: "pages"
  cancel-in-progress: false

jobs:
  deploy:
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/configure-pages@v5
      - uses: actions/upload-pages-artifact@v3
        with:
          path: docs
      - id: deployment
        uses: actions/deploy-pages@v4
```

## Suggested docs structure

- `docs/index.md`
- `docs/quickstart.md`
- `docs/route-coverage.md`
- `docs/docker.md`
- `docs/deploy.md`
- `docs/faq.md`
- `docs/demo-report.html`
- `docs/dashboard-overview.png`
- `docs/files-list.png`
- `docs/explorer-tree-main-go.png`
- `docs/route-coverage-overview.png`

## Recommended publish flow

1. Update content in `plan/github-pages/`
2. Sync to `docs/` with the helper script
3. Commit and push
4. Verify published pages URL

## Recommended homepage sections

- "Try in 2 minutes" quick command
- Feature list (code + route coverage)
- Demo link
- Deployment guide link
