# Registry relations

Дата live-проверки: 2026-06-25.

Проверенная запись:

- registry: `16`
- item id: `1404476`
- registry number: `ФР.1.31.2022.44733`

Machine-readable summary:

- `docs/research/registry-relations.json`

## Relation forms

Two relation-like field types were observed in the sectioned item card:

- `LINK_INTERNAL`
- `CHILD_OBJECT`

## LINK_INTERNAL

`LINK_INTERNAL` fields store a human-readable value and a site-relative link:

```json
{
  "name": "foei:CertifOrgCMM",
  "type": "LINK_INTERNAL",
  "value": "АХУ УРО РАН",
  "link": "/registry/13/items/394779"
}
```

The link format can be parsed as:

```text
/registry/{registryId}/items/{itemId}
```

Confirmed links from item `1404476`:

| Field | Registry | Item | Value |
| --- | --- | --- | --- |
| `foei:CertifOrgCMM` | `13` | `394779` | АХУ УРО РАН |
| `foei:DevMetdCMM` | `20` | `272627` | ФГБУ «ФЦАО» |
| `foei:SFCMM` | `23` | `239512` | Вне сферы обороны |
| `foei:TypeMeasuCMM` | `45` | `269711` | Физико-химические измерения |

The linked organization record `/registry/13/items/394779/plaindata` was reachable and returned type `foei:AL_type`.

## CHILD_OBJECT

`CHILD_OBJECT` fields store nested objects in `value`.

Observed field:

```text
foei:CMM1Relation_assoc
```

Observed child object shape:

```json
{
  "id": "...",
  "type": "foei:CMM1Relation",
  "nodeRef": "workspace://SpacesStore/...",
  "columns": [],
  "fields": []
}
```

For item `1404476`, the child object means:

| Child type | Relation type | Related method | Related link |
| --- | --- | --- | --- |
| `foei:CMM1Relation` | `Заменяет` | `ФР.1.31.2013.16016` | `/registry/16/items/286196` |

The related item `/registry/16/items/286196/plaindata` was reachable and returned:

- type: `foei:CMM1_type`
- registry number: `ФР.1.31.2013.16016`
- record number: `38286`

## Implementation

Added generic helpers:

- `ParseRegistryItemLink`
- `RegistryProperty.LinkStrings`
- `RegistryProperty.ItemRefs`
- `RegistryProperty.ChildObjects`
- `ChildObject.FieldsByName`
- `ChildObject.ItemRefs`

These helpers intentionally do not assume registry `16`; they work with any `/registry/{id}/items/{itemId}` link.

## Open questions

- Whether all internal links use `/registry/{id}/items/{itemId}`.
- Whether `CHILD_OBJECT` values can be single objects instead of arrays.
- Whether `columns` always matches the `fields` set.
- Whether there are relation endpoints beyond embedded child objects.
- Whether hidden/reference registries such as `13`, `20`, `23`, `45` should be added to the static catalog or kept as discovered references.
