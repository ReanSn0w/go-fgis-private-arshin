# AGENTS.md

Памятка для Codex/агентов, продолжающих работу над `go-fgis-private-arshin`.

## Проект

Это Go-библиотека для приватного web API сайта ФГИС Аршин:

```text
https://fgis.gost.ru/fundmetrology/api/
```

Важно: это не официальный публичный `eapi`, а API, которым пользуется web-интерфейс сайта `https://fgis.gost.ru/fundmetrology/registry/{id}`. Схема может измениться без предупреждения, поэтому в проекте есть исследовательский CLI `arshin-probe` и snapshots в `docs/research`.

Не смешивать этот клиент с официальным публичным `eapi`. Здесь строится generic-клиент web-реестров и typed-адаптеры поверх него.

## Быстрый старт

Проверка без live-запросов:

```bash
go test ./...
```

Обычные тесты не должны ходить в сеть.

Исследовательская утилита:

```bash
go run ./cmd/arshin-probe -h
```

`arshin-probe` - не runtime API для пользователей, а диагностический инструмент. Он нужен, чтобы снимать текущую правду с сайта, обновлять fixtures/snapshots и быстро понимать, что сломалось, если API изменился.

## Главные файлы контекста

- `README.md` - текущее пользовательское описание API и команд.
- `PLAN.md` - план исследования/реализации, статусы шагов и оставшиеся вопросы.
- `NOTES.md` - исходные заметки по публичному `eapi`, registry `16` и проверенной записи.
- `REGISTRY.md` - первичный список публичных web-разделов, похожих на registry `16`.
- `docs/research/*.md` - человекочитаемые отчеты live-исследований.
- `docs/research/*.json` - машинные snapshots live-исследований.
- `testdata/*.json` - fixtures для тестов.

## Подтвержденные endpoints

Generic list/data:

```text
GET /fundmetrology/api/registry/{registryId}/data
```

Параметры:

- `pageNumber` - 1-based.
- `pageSize` - дефолт в клиенте 20.
- `orgID` - обычно `CURRENT_ORG`.
- повторяющиеся `filterBy` + `filterValues`.
- повторяющиеся `sortBy` + `sortValues`.

Метаданные:

```text
GET /fundmetrology/api/registry/{registryId}
GET /fundmetrology/api/registry/{registryId}/fields
GET /fundmetrology/api/registry/{registryId}/displayfields
GET /fundmetrology/api/fieldtypes
```

Карточки:

```text
GET /fundmetrology/api/registry/{registryId}/items/{itemId}/data
GET /fundmetrology/api/registry/{registryId}/items/{itemId}/plaindata
```

Вложения:

```text
GET /fundmetrology/api/downloadfile/{fileId}
```

## Реализованный Go API

Generic client package: `github.com/ReanSn0w/go-fgis-private-arshin/public/client`, package name `client`.
Registry catalog package: `github.com/ReanSn0w/go-fgis-private-arshin/public/registries`.

Основное:

- `NewClient(...)`
- `ListRegistryData(ctx, registryID, query)`
- `GetRegistry(ctx, registryID)`
- `ListRegistryFields(ctx, registryID)`
- `ListRegistryDisplayFields(ctx, registryID)`
- `ListFieldTypes(ctx)`
- `GetRegistryItem(ctx, registryID, itemID)`
- `GetRegistryItemPlainData(ctx, registryID, itemID)`
- `DownloadFile(ctx, fileID)`
- `DownloadFileLink(ctx, link)`

Generic helpers:

- `PropertiesByName`
- `RegistryRecord.PropertiesByName`
- `RegistryItem.Fields`
- `RegistryItem.FieldsByName`
- `ParseRegistryItemLink`
- `ParseDownloadFileLink`
- `RegistryProperty.Attachments`
- `RegistryProperty.LinkStrings`
- `RegistryProperty.ItemRefs`
- `RegistryProperty.ChildObjects`
- `ChildObject.FieldsByName`
- `ChildObject.ItemRefs`

Typed adapter:

- `public/registries` - общий каталог `Spec` для всех известных public web-реестров, включая named specs `registries.CMM1`, `registries.CMM2`, `registries.SU`, `registries.UTSO` и т.д.
- `public/registries/methods` - typed adapter for measurement method registries `6`, `8`, `16`
- Registries `6`, `8`, `16`
- `methods.New(client)` defaults to registry `16`
- `methods.NewForRegistry(client, registryID)`
- `methods.NewForSpec(client, spec)`
- `methods.SpecForRegistry`, `methods.SpecForItemType`, `methods.SupportedRegistries`
- `GetByRegistryNumber(ctx, "ФР.1.31.2022.44733")`
- `GetByRegistryNumberDetailed(ctx, "...")`

## Проверенная запись

Главная fixture-запись:

- registry: `16`
- item id: `1404476`
- registry number: `ФР.1.31.2022.44733`
- document: `ПНД Ф 14.1:2:3:4.48-2022`

Fixtures:

- `testdata/fr-1-31-2022-44733.json` - list/data search result.
- `testdata/item-1404476-data.json` - sectioned card.
- `testdata/item-1404476-plaindata.json` - flat card.

## Research snapshots

Current reports:

- `docs/research/registry-survey.md`
- `docs/research/registry-metadata.md`
- `docs/research/registry-data-behavior.md`
- `docs/research/registry-item-details.md`
- `docs/research/registry-attachments.md`
- `docs/research/registry-relations.md`
- `docs/research/registry-session-security.md`
- `docs/research/typed-methods.md`

Corresponding JSON:

- `docs/research/registry-survey.json`
- `docs/research/registry-metadata.json`
- `docs/research/registry-data-behavior.json`
- `docs/research/registry-relations.json`
- `docs/research/registry-session-security.json`

## Useful probe commands

Survey known registries:

```bash
go run ./cmd/arshin-probe -all -size 2 -out docs/research/registry-survey.json
```

Fetch metadata:

```bash
go run ./cmd/arshin-probe -all -metadata -out docs/research/registry-metadata.json
```

Check `/data` behavior:

```bash
go run ./cmd/arshin-probe -behavior -out docs/research/registry-data-behavior.json
```

Fetch item card:

```bash
go run ./cmd/arshin-probe -registry 16 -item 1404476 -out testdata/item-1404476-data.json
go run ./cmd/arshin-probe -registry 16 -item 1404476 -plain -out testdata/item-1404476-plaindata.json
```

Download attachment:

```bash
go run ./cmd/arshin-probe -download c11f089b-22f4-4d66-87cd-7aff99909001 -out /tmp/arshin-rmdoc.pdf
```

Print item relations:

```bash
go run ./cmd/arshin-probe -registry 16 -item 1404476 -relations
```

Check session/header/cookie behavior:

```bash
go run ./cmd/arshin-probe -session -out docs/research/registry-session-security.json
```

## Known behavior and caveats

- `/data` empty filter result is a successful response with `totalCount: 0`.
- Invalid registry id returned HTTP `500 Internal Server Error`, not clean `404`.
- `pageSize=100` and `500` worked for registry `16`; `1000` timed out. Keep default conservative.
- List responses observed `values: null`; useful data is in `properties[]`.
- `/items/{id}/data` returns sectioned `sections[].fields[]`.
- `/items/{id}/plaindata` returns flat `RegistryRecord` with `properties[]`.
- `RegistryProperty.Link` and `RegistryProperty.Value` intentionally stay `any`: observed values include strings, arrays, objects, and null.
- `ATTACH` download returns `application/octet-stream`; better MIME can be in field `mime`.
- Download filename in `Content-Disposition` is percent-encoded.
- Missing download file id returns HTTP `500` with JSON `APIError`.
- Registry `47` metadata type is `gost:p1wfRequest`, but list item type is `gost:p1wfRequestType4`; keep registry type and item type separate.
- Internal links can point to reference registries not in `REGISTRY.md`, e.g. `13`, `20`, `23`, `45`.
- Public GET endpoints currently work without cookie preflight and without `Ajax-Token`, but server sets `session-cookie`; keep cookie jar anyway.
- HTML registry pages contain `Ajax-Token` token name, but checked GET API calls did not require it.
- HTML pages can return status `200` even with `Accept: application/json`; keep content-type checks.

## Implemented plan status

Completed or mostly completed:

1. Catalog/survey of initial registries from `REGISTRY.md`.
2. Registry metadata endpoints.
3. Generic `/data` behavior checks.
4. Item card endpoints.
5. Attachment downloads.
6. Generic relation helpers.
7. Session/header/cookie behavior for public GET endpoints.
8. First typed adapter for CMM registries `6`, `8`, `16`.

Likely next steps:

- Schema layer: local prepared metadata from `docs/research/registry-metadata.json`, so runtime code does not need `/fields` every time.
- Broaden fixtures beyond registry `16`: especially `11`, `19`, `47`, and reference registries reached by links.
- Add live fixtures for CMM registries `6` and `8`, especially order fields and relation fields.
- Decide whether registries `11`, `19`, `47` deserve typed adapters.

## Style and constraints

- Prefer generic client first; typed packages should wrap generic API, not duplicate HTTP logic.
- Keep live network checks in `arshin-probe`, not normal unit tests.
- Keep conservative rate limit. The default client rate period is 1 second.
- Preserve raw records/properties in typed models because schema can change.
- Do not assume all registries are public menu entries; linked reference registries may still be accessible.
- Use `go test ./...` before handing off.
