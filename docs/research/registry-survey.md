# Registry survey

Дата live-проверки: 2026-06-25.

Команда:

```bash
go run ./cmd/arshin-probe -all -size 2 -out docs/research/registry-survey.json
```

Все разделы из `REGISTRY.md` ответили через generic endpoint:

```text
GET https://fgis.gost.ru/fundmetrology/api/registry/{registryId}/data
```

## Summary

| ID | OK | Total count | Record type | Properties |
| --- | --- | ---: | --- | ---: |
| `1` | yes | 104 | `foei:ND_type` | 14 |
| `2` | yes | 783 | `foei:SCM_type` | 13 |
| `3` | yes | 514 | `foei:GSI_type` | 10 |
| `6` | yes | 17 | `foei:CMM2_type` | 33 |
| `7` | yes | 12 | `foei:MDG_type` | 12 |
| `8` | yes | 7 | `foei:CMM3_type` | 35 |
| `9` | yes | 437 | `foei:IC_type` | 16 |
| `10` | yes | 483 | `foei:SSD_type` | 12 |
| `11` | yes | 117537 | `foei:SU_type` | 22 |
| `12` | yes | 221 | `foei:GPS_type` | 46 |
| `14` | yes | 198 | `foei:MD_type` | 16 |
| `16` | yes | 52635 | `foei:CMM1_type` | 32 |
| `17` | yes | 1742 | `foei:TSSI_type` | 11 |
| `18` | yes | 1 | `foei:EPI_type` | 2 |
| `19` | yes | 11268 | `foei:UTSO_type` | 18 |
| `47` | yes | 4878 | `gost:p1wfRequestType4` | 20 |

## Notes

- The generic response shape is usable across all surveyed registry IDs.
- Registry `47` uses the `gost:` namespace, while most other surveyed sections use `foei:`.
- Several registries expose `ATTACH`, `LINK_INTERNAL`, and `CHILD_OBJECT` fields in list responses, so attachments and relations should be implemented in the generic layer before adding many typed adapters.
- Registry `11` is the largest surveyed section and should be used for pagination/rate-limit testing.
- Registries `6`, `8`, and `16` are related method registries with similar `CMM` field families and are good candidates for shared typed mapping helpers.
