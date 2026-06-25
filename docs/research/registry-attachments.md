# Registry attachments

Дата live-проверки: 2026-06-25.

Проверенная запись:

- registry: `16`
- item id: `1404476`
- field: `foei:RMDocCMM`
- filename from field value: `ПРИЛОЖЕНИЕ_048.pdf`
- link from field value: `/api/downloadfile/c11f089b-22f4-4d66-87cd-7aff99909001`
- mime from field value: `application/pdf`

## Confirmed endpoint

```text
GET /fundmetrology/api/downloadfile/{fileId}
```

The same endpoint can be reached through the relative link stored in an `ATTACH` field:

```text
/api/downloadfile/{fileId}
```

## Successful download

Observed headers:

```text
HTTP/1.1 200 OK
Content-Type: application/octet-stream
Content-Length: 101187
content-disposition: form-data; name="attachment"; filename="%D0%9F%D0%A0%D0%98%D0%9B%D0%9E%D0%96%D0%95%D0%9D%D0%98%D0%95_048.pdf"
cache-control: max-age=604800
```

Observed file:

```text
PDF document, version 1.7, 1 pages
sha256: c69deb1576a8efbf5dcc126011e2ee7c97f980a0b1579547cadf07bb3bd64946
```

Notes:

- Server returns `application/octet-stream` even for PDF.
- Better MIME can be available in the original `ATTACH` property (`application/pdf` in this case).
- Filename in `Content-Disposition` is percent-encoded and must be decoded.
- Server sets a `session-cookie`, but this public download worked without a pre-existing cookie.

## Missing file

An unknown file id returned:

```text
HTTP/1.1 500 Internal Server Error
Content-Type: application/json;charset=UTF-8
```

The JSON body contains:

```json
{
  "status": 500,
  "result": null,
  "message": "File with ... not found in the internal DB!",
  "trace": "..."
}
```

So download code should attempt to parse non-2xx JSON responses as `APIError`.

## Implementation

Added:

- `DownloadFile(ctx, fileID)`
- `DownloadFileLink(ctx, link)`
- `DownloadedFile`
- `arshin-probe -download`

Examples:

```bash
go run ./cmd/arshin-probe -download c11f089b-22f4-4d66-87cd-7aff99909001 -out /tmp/arshin-rmdoc.pdf
go run ./cmd/arshin-probe -download /api/downloadfile/c11f089b-22f4-4d66-87cd-7aff99909001 -out /tmp/arshin-rmdoc.pdf
```

## Open questions

- Whether private/authenticated files behave differently.
- Whether all `ATTACH` fields use `/api/downloadfile/{id}`.
- Whether multi-file attachments always use parallel arrays in `value` and `link`.
- Whether `mime` can be an array for multi-file fields.
