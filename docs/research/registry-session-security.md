# Registry session and headers

Дата live-проверки: 2026-06-25.

Команда:

```bash
go run ./cmd/arshin-probe -session -out docs/research/registry-session-security.json
```

## Checked endpoints

- data: `/fundmetrology/api/registry/16/data?...`
- page preflight: `/fundmetrology/registry/16`
- item card: `/fundmetrology/api/registry/16/items/1404476/data`
- download: `/fundmetrology/api/downloadfile/c11f089b-22f4-4d66-87cd-7aff99909001`
- invalid registry: `/fundmetrology/api/registry/999999/data?...`

## Summary

| Check | Status | Content type | Set-Cookie | Cookie after | Ajax token | JSON | HTML |
| --- | ---: | --- | --- | --- | --- | --- | --- |
| `data_plain_no_browser_headers_no_cookie` | 200 | `application/json;charset=UTF-8` | `session-cookie` |  |  | yes | no |
| `data_browser_headers_no_cookie` | 200 | `application/json;charset=UTF-8` | `session-cookie` |  |  | yes | no |
| `registry_page_preflight_sets_cookie_and_token` | 200 | `text/html; charset=utf-8` | `session-cookie` | `session-cookie` | `Ajax-Token` | no | yes |
| `data_after_preflight_cookie` | 200 | `application/json;charset=UTF-8` |  | `session-cookie` |  | yes | no |
| `data_after_preflight_cookie_with_ajax_token_header` | 200 | `application/json;charset=UTF-8` |  | `session-cookie` |  | yes | no |
| `item_after_preflight_cookie` | 200 | `application/json;charset=UTF-8` |  | `session-cookie` |  | yes | no |
| `download_after_preflight_cookie` | 200 | `application/octet-stream` |  | `session-cookie` |  | no | no |
| `invalid_registry_json_error` | 500 | `application/json;charset=UTF-8` | `session-cookie` |  |  | yes | no |
| `html_page_when_json_expected` | 200 | `text/html; charset=utf-8` | `session-cookie` |  | `Ajax-Token` | no | yes |

## Findings

- Public GET data endpoint worked without prior cookie, without Referer, and without browser-like headers.
- Browser-like headers also worked.
- Registry HTML page sets `session-cookie` and contains:

```html
<meta name='hmac-token-name' content='Ajax-Token'/>
```

- After HTML preflight, data, item card, and download endpoints still worked with the cookie jar.
- Supplying an `Ajax-Token` header for a GET request did not appear necessary; the checked GET request worked both with and without it.
- Download endpoint worked after preflight and returned binary `application/octet-stream`.
- Invalid registry id returned HTTP `500`, but JSON body had `status: 404`, `message: "not found"`, and `trace`.
- Requesting the HTML registry page with `Accept: application/json` still returned HTML with status `200`; content-type checking is required.

## Implementation notes

- Keep the default client profile conservative:
  - cookie jar;
  - explicit `User-Agent`;
  - `Accept: application/json, text/plain, */*`;
  - `Accept-Language: ru-RU,ru;q=0.9`;
  - registry `Referer` for registry-bound API calls;
  - 1 request/second default rate limit.
- Do not require cookie preflight for public GET calls.
- Do not require `Ajax-Token` for public GET calls.
- Continue checking `Content-Type` to catch HTML returned where JSON was expected.
- Parse JSON error bodies on non-2xx responses as `APIError`; HTTP status and API body status can differ.

## Not tested

- Mutating POST/PUT/DELETE endpoints. No safe public mutation endpoint has been identified, and the library currently only needs GET operations.
- Aggressive rate-limit behavior. Do not stress the public service; keep conservative request pacing.
