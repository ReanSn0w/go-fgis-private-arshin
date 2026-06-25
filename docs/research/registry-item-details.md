# Registry item details

Дата live-проверки: 2026-06-25.

Проверенная запись:

- registry: `16`
- item id: `1404476`
- registry number: `ФР.1.31.2022.44733`

## Confirmed endpoints

```text
GET /fundmetrology/api/registry/{registryId}/items/{itemId}/data
GET /fundmetrology/api/registry/{registryId}/items/{itemId}/plaindata
```

Fixtures:

- `testdata/item-1404476-data.json`
- `testdata/item-1404476-plaindata.json`

## `/data`

`/items/{itemId}/data` returns a sectioned card shape:

```json
{
  "status": 200,
  "result": {
    "id": 1404476,
    "alfrescoId": "...",
    "registryId": 16,
    "registryTitle": null,
    "deleted": false,
    "sections": []
  }
}
```

For item `1404476`, the card has 5 sections:

| Section | Fields |
| --- | ---: |
| Аттестующая организация | 6 |
| Сведения о разработчике МВИ | 4 |
| Сведения об аттестованной методике (методе) измерений | 14 |
| Вложения | 0 |
| Общие сведения | 3 |

The sectioned response is the better source for UI/card rendering and for discovering links, attachments, and child objects as the site presents them.

## `/plaindata`

`/items/{itemId}/plaindata` returns the flat `RegistryRecord` shape:

```json
{
  "status": 200,
  "result": {
    "values": null,
    "properties": [],
    "id": "1404476",
    "alfrescoId": "...",
    "type": "foei:CMM1_type"
  }
}
```

For item `1404476`, `plaindata` returned 30 properties.

The plain response is the better source for typed mappers and machine processing, because it matches the list/data `properties[]` model.

## Observed links and nested data

Sectioned `/data` exposed:

- `LINK_INTERNAL` to organization records:
  - `/registry/13/items/394779`
  - `/registry/20/items/272627`
  - `/registry/23/items/239512`
  - `/registry/45/items/269711`
- `ATTACH` field:
  - `foei:RMDocCMM`
  - value: `["ПРИЛОЖЕНИЕ_048.pdf"]`
  - link: `["/api/downloadfile/c11f089b-22f4-4d66-87cd-7aff99909001"]`
  - mime: `application/pdf`
- `CHILD_OBJECT` relation:
  - `foei:CMM1Relation_assoc`
  - relation type: `Заменяет`
  - related method: `ФР.1.31.2013.16016`
  - related link: `/registry/16/items/286196`

## Implementation notes

- Keep separate models for sectioned item cards and flat records.
- Use `GetRegistryItem` for sectioned `/data`.
- Use `GetRegistryItemPlainData` for flat `/plaindata`.
- Do not force `RegistryProperty.Link` into a string yet: in observed data it can be a string, an array, or null.
- Child object values should stay raw until relation support is implemented.
