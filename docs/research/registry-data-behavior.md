# Registry data behavior

Дата live-проверки: 2026-06-25.

Команда:

```bash
go run ./cmd/arshin-probe -behavior -out docs/research/registry-data-behavior.json
```

Проверяемый endpoint:

```text
GET /fundmetrology/api/registry/{registryId}/data
```

## Summary

| Check | Registry | OK | Total | Page | Page size | Items | Types | First ID | First number |
| --- | --- | --- | ---: | ---: | ---: | ---: | --- | --- | --- |
| `first_page_registry_16` | `16` | yes | 52635 | 1 | 3 | 3 | `foei:CMM1_type` | `1430402` | `54152` |
| `pagination_registry_16_page_2` | `16` | yes | 52635 | 2 | 3 | 3 | `foei:CMM1_type` | `1430355` | `54122` |
| `empty_filter_registry_16` | `16` | yes | 0 | 1 | 3 | 0 |  |  |  |
| `single_filter_registry_16` | `16` | yes | 1 | 1 | 3 | 1 | `foei:CMM1_type` | `1404476` | `44733` |
| `multiple_filters_registry_16` | `16` | yes | 1 | 1 | 3 | 1 | `foei:CMM1_type` | `1404476` | `44733` |
| `sort_asc_registry_16` | `16` | yes | 52635 | 1 | 3 | 3 | `foei:CMM1_type` | `273736` | `25826` |
| `sort_desc_registry_16` | `16` | yes | 52635 | 1 | 3 | 3 | `foei:CMM1_type` | `302490` |  |
| `larger_page_size_registry_16` | `16` | yes | 52635 | 1 | 100 | 100 | `foei:CMM1_type` | `1430402` | `54152` |
| `large_page_size_500_registry_16` | `16` | yes | 52635 | 1 | 500 | 500 | `foei:CMM1_type` | `1430402` | `54152` |
| `small_registry_18` | `18` | yes | 1 | 1 | 10 | 1 | `foei:EPI_type` | `493452` |  |
| `largest_registry_11` | `11` | yes | 117537 | 1 | 3 | 3 | `foei:SU_type` | `1430363` | `196230` |
| `registry_47_namespace` | `47` | yes | 4878 | 1 | 3 | 3 | `gost:p1wfRequestType4` | `500262` |  |
| `invalid_registry` | `999999` | expected failure | 0 | 0 | 0 | 0 |  |  |  |

## Confirmed behavior

- Root response shape is stable for successful data requests: `status`, `result`, `message`, `trace`.
- `result` contains `totalCount`, `currentPage`, `pageSize`, `items`.
- Empty search results are successful responses with `totalCount: 0` and no items.
- Repeated `filterBy/filterValues` pairs work for multiple filters.
- Repeated `sortBy/sortValues` pairs are accepted; single-field ascending and descending sort both returned successful pages.
- Pagination is 1-based and returns the requested `currentPage`.
- `pageSize=100` and `pageSize=500` worked for registry `16`; `pageSize=1000` timed out in a separate manual probe, so the client should keep conservative defaults.
- In sampled list responses, `RegistryRecord.values` was null and useful data was in `properties[]`.
- Invalid registry id `999999` returned HTTP `500 Internal Server Error`, not a neat `404`.
- Registry `47` list items use `gost:p1wfRequestType4`, while metadata reports `gost:p1wfRequest`.

## Implementation notes

- Keep default `pageSize=20`.
- Treat `totalCount=0` as a normal empty result, not an error.
- Do not assume invalid ids return `404`; preserve HTTP status errors.
- Store record `type` from data responses separately from registry metadata `type`.
- Prefer small pages in live tests; large pages can be slow and produce multi-megabyte responses.
