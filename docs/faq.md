# FAQ

## Is there a ready-to-open demo?

Yes. Open [demo-report.html](demo-report.html) in this docs folder.

## Does it require a running server to view reports?

No. CLI mode produces a single static HTML file.

## Can I use it in CI?

Yes. Generate `coverage.html` and upload as an artifact.

## Can I use it with private repositories?

Yes. Everything can run locally or in your own CI/CD.

## Can I include API route coverage?

Yes, with an optional `.rcov` file.

## Is Docker required?

No. Docker is optional for shared team dashboards.

## I got: "could not detect module prefix from coverage file"

This was a common issue with mixed or prefix-less paths in `cover.out`.

Current behavior uses automatic path resolution even when no common module prefix is detected.

## Can I keep reports private?

Yes. Generate reports in private CI and publish artifacts only to authorized users.

## Does route coverage need OpenAPI?

No. `.rcov` is a simple text format and can be produced by any test tooling.
