# Registry metadata

Дата live-проверки: 2026-06-25.

Команда:

```bash
go run ./cmd/arshin-probe -all -metadata -out docs/research/registry-metadata.json
```

## Confirmed endpoints

```text
GET /fundmetrology/api/registry/{registryId}
GET /fundmetrology/api/registry/{registryId}/fields
GET /fundmetrology/api/registry/{registryId}/displayfields
GET /fundmetrology/api/fieldtypes
```

`/registry/{registryId}` returns registry-level metadata: id, title, Alfresco type, nodeRef, timestamps, deletion flag.

`/registry/{registryId}/fields` returns all known fields for the registry.

`/registry/{registryId}/displayfields` returns fields used as table columns in the site UI.

`/fieldtypes` returns the global field type classification.

## Field types

| Type | Link | Attach | ID |
| --- | --- | --- | --- |
| `any` | no | no | no |
| `id` | no | no | yes |
| `text` | no | no | no |
| `content` | no | no | no |
| `number` | no | no | no |
| `date` | no | no | no |
| `ATTACH` | no | yes | no |
| `LIST` | yes | no | no |
| `LINK_INTERNAL` | yes | no | no |
| `LINK_EXTERNAL` | yes | no | no |
| `CHILD_OBJECT` | no | no | no |

## Registry metadata summary

| ID | OK | Metadata type | Fields | Display fields |
| --- | --- | --- | ---: | ---: |
| `1` | yes | `foei:ND_type` | 17 | 8 |
| `2` | yes | `foei:SCM_type` | 14 | 4 |
| `3` | yes | `foei:GSI_type` | 14 | 6 |
| `6` | yes | `foei:CMM2_type` | 46 | 5 |
| `7` | yes | `foei:MDG_type` | 15 | 6 |
| `8` | yes | `foei:CMM3_type` | 46 | 5 |
| `9` | yes | `foei:IC_type` | 25 | 7 |
| `10` | yes | `foei:SSD_type` | 15 | 6 |
| `11` | yes | `foei:SU_type` | 44 | 5 |
| `12` | yes | `foei:GPS_type` | 66 | 7 |
| `14` | yes | `foei:MD_type` | 19 | 6 |
| `16` | yes | `foei:CMM1_type` | 43 | 5 |
| `17` | yes | `foei:TSSI_type` | 14 | 6 |
| `18` | yes | `foei:EPI_type` | 2 | 2 |
| `19` | yes | `foei:UTSO_type` | 21 | 6 |
| `47` | yes | `gost:p1wfRequest` | 28 | 6 |

## Notes

- Registry `47` reports metadata type `gost:p1wfRequest`, while list data previously observed item type `gost:p1wfRequestType4`. The client should keep both registry metadata type and item type observations.
- `fields` is richer than fields observed on the first data page; for example registry `16` has 43 metadata fields, while a first-page data sample showed 32 populated/listed fields.
- `displayfields` matches table-column metadata and is the best source for default column selection.
- `LIST` fields expose constraint names, but this survey did not find enum value dictionaries yet. That remains a separate research item.
- `ATTACH`, `LINK_INTERNAL`, `LINK_EXTERNAL`, and `CHILD_OBJECT` are visible at metadata level, so attachment/relation support can be driven from field metadata instead of guessed from record values.
