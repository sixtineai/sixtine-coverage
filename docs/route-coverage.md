# Route coverage format (.rcov)

`sixtine-coverage` can merge route-level coverage with code coverage.

Try it directly in the bundled demo report: [demo-report.html](demo-report.html).

## Route coverage UI

![Route coverage summary cards](route-coverage.png)

## Example

```text
mode: status
# AuthenticationService
GET /api/rest/v1/authentication/modes 200:1
POST /api/rest/v1/authentication/login 2 200:3
# TenantService
GET /api/rest/v1/tenants 200:4
GET /api/rest/v1/tenants/{id} 0
```

## Rules

- Header: `mode: status`
- Service groups: `# ServiceName`
- Route line: `METHOD /path counts`
- `0` means no hits
- `200:3` means three responses with HTTP 200
- `2 200:3` means two unspecified + three with status 200
- `-1:2` means two connection failures

## Route line grammar

```text
<METHOD> <PATH> <TOKEN> [<TOKEN> ...]
```

Where each `<TOKEN>` is one of:

- `0`
- `<count>` (unspecified responses)
- `<status>:<count>`

Examples:

- `GET /api/rest/v1/tenants 200:12`
- `POST /api/rest/v1/login 3 200:10 401:2`
- `GET /api/rest/v1/admin/metrics 0`

## Recommended conventions

- Group by bounded context (`# TenantService`, `# AuthenticationService`)
- Keep HTTP methods uppercase
- Use normalized route templates (`/tenants/{id}`)
- Include failure statuses to reveal blind spots

## Notes on `-1`

`-1:count` is useful for transport or network failures where no HTTP status code was returned.

## Validation checklist

- First non-empty line is `mode: status`
- Every route belongs to a service section
- No duplicate `METHOD + PATH` within the same service (recommended)
- Counts are non-negative integers
