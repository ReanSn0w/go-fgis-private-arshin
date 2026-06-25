# go-fgis-private-arshin

Go-библиотека для работы с приватным web API сайта ФГИС Аршин.

Важно: это не официальный публичный `eapi`, а API, которым пользуется web-интерфейс сайта `https://fgis.gost.ru/fundmetrology/registry/{id}`. Схема и поведение могут измениться без предупреждения, поэтому клиент по умолчанию использует осторожный rate limit и сохраняет raw-поля записей.

## Статус

Проект находится на этапе исследования и первичной реализации.

Уже есть:

- generic-клиент для `GET /fundmetrology/api/registry/{id}/data`;
- методы метаданных для `/registry/{id}`, `/fields`, `/displayfields`, `/fieldtypes`;
- стартовый каталог известных registry-разделов из `REGISTRY.md`;
- исследовательский CLI `arshin-probe`;
- typed-адаптер для семейства реестров методик измерений `6`, `8`, `16`;
- live-сводка по доступным разделам в `docs/research/registry-survey.md`.
- live-сводка по метаданным разделов в `docs/research/registry-metadata.md`.
- live-проверка поведения generic `/data` endpoint в `docs/research/registry-data-behavior.md`.
- методы карточки записи для `/items/{itemId}/data` и `/items/{itemId}/plaindata`.
- скачивание вложений через `/api/downloadfile/{fileId}`.
- generic helpers для `ATTACH`, `LINK_INTERNAL` и `CHILD_OBJECT` полей.
- live-проверка session/header/cookie поведения в `docs/research/registry-session-security.md`.
- typed-исследование семейства реестров методик измерений в `docs/research/typed-methods.md`.

## Generic usage

```go
client, err := arshinclient.NewClient()
if err != nil {
    return err
}

page, err := client.ListRegistryData(ctx, "16", arshinclient.RegistryQuery{
    PageNumber: 1,
    PageSize:   20,
    Filters: []arshinclient.Filter{
        {Field: "foei:NumRegCMM", Value: "ФР.1.31.2022.44733"},
    },
})
if err != nil {
    return err
}
```

## Metadata usage

```go
details, err := client.GetRegistry(ctx, "16")
fields, err := client.ListRegistryFields(ctx, "16")
displayFields, err := client.ListRegistryDisplayFields(ctx, "16")
fieldTypes, err := client.ListFieldTypes(ctx)
```

## Item details usage

```go
item, err := client.GetRegistryItem(ctx, "16", "1404476")
plainRecord, err := client.GetRegistryItemPlainData(ctx, "16", "1404476")
```

`GetRegistryItem` возвращает секционную карточку как на сайте. `GetRegistryItemPlainData` возвращает плоский `RegistryRecord` с `properties[]`.

## Attachment usage

```go
item, err := client.GetRegistryItem(ctx, "16", "1404476")
if err != nil {
    return err
}

attachments := item.FieldsByName()["foei:RMDocCMM"].Attachments()

file, err := client.DownloadFile(ctx, "c11f089b-22f4-4d66-87cd-7aff99909001")
if err != nil {
    return err
}
defer file.Body.Close()
```

## Relation usage

```go
item, err := client.GetRegistryItem(ctx, "16", "1404476")
if err != nil {
    return err
}

field := item.FieldsByName()["foei:CMM1Relation_assoc"]
children, err := field.ChildObjects()
if err != nil {
    return err
}

refs := children[0].ItemRefs()
```

## Measurement methods typed usage

```go
client, err := arshinclient.NewClient()
if err != nil {
    return err
}

methodsClient := methods.New(client)
method, err := methodsClient.GetByRegistryNumber(ctx, "ФР.1.31.2022.44733")
if err != nil {
    return err
}

detailed, err := methodsClient.GetByRegistryNumberDetailed(ctx, "ФР.1.31.2022.44733")
if err != nil {
    return err
}
_ = detailed.RangeDocuments
_ = detailed.Relations
```

`methods.New(client)` использует registry `16`. Для реестров первичных и референтных методик:

```go
primary, err := methods.NewForRegistry(client, methods.PrimaryReferenceRegistryID) // registry 6
reference, err := methods.NewForRegistry(client, methods.ReferenceRegistryID)      // registry 8
```

Все известные registry specs лежат в общем пакете `registries`:

```go
import (
    arshinclient "github.com/ReanSn0w/go-fgis-private-arshin/public/client"
    "github.com/ReanSn0w/go-fgis-private-arshin/public/registries"
    "github.com/ReanSn0w/go-fgis-private-arshin/public/registries/methods"
)

primary, err := methods.NewForSpec(client, registries.CMM2)
```

Для реестров без typed-адаптера используется generic API и тот же каталог:

```go
specs := registries.Known()
standardUnits := registries.SU   // registry 11
standardSamples := registries.UTSO // registry 19
```

## Examples

В директории [`example`](./example) лежит простая утилита, которая ищет записи
по названию документа во всех известных web-реестрах и выводит JSON со списком
найденных записей и внутренних связей.

```bash
go run ./example -query "ПНД Ф 14.1:2:3:4.48-2022"
```

## Probe CLI

Получить JSON первой страницы реестра:

```bash
go run ./cmd/arshin-probe -registry 16 -size 5
```

Вывести список полей по первым записям:

```bash
go run ./cmd/arshin-probe -registry 16 -size 2 -properties
```

Собрать сводку по всем известным разделам:

```bash
go run ./cmd/arshin-probe -all -size 2 -out docs/research/registry-survey.json
```

Собрать метаданные по всем известным разделам:

```bash
go run ./cmd/arshin-probe -all -metadata -out docs/research/registry-metadata.json
```

Проверить поведение generic `/data` endpoint:

```bash
go run ./cmd/arshin-probe -behavior -out docs/research/registry-data-behavior.json
```

Сохранить карточку записи:

```bash
go run ./cmd/arshin-probe -registry 16 -item 1404476 -out testdata/item-1404476-data.json
go run ./cmd/arshin-probe -registry 16 -item 1404476 -plain -out testdata/item-1404476-plaindata.json
```

Скачать вложение:

```bash
go run ./cmd/arshin-probe -download c11f089b-22f4-4d66-87cd-7aff99909001 -out /tmp/arshin-rmdoc.pdf
```

Вывести связи карточки:

```bash
go run ./cmd/arshin-probe -registry 16 -item 1404476 -relations
```

Проверить session/header/cookie поведение:

```bash
go run ./cmd/arshin-probe -session -out docs/research/registry-session-security.json
```

## Tests

Обычные тесты не делают live-запросы к сайту:

```bash
go test ./...
```

## License

Проект распространяется под лицензией MIT. См. [`LICENSE`](./LICENSE).
