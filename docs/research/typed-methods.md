# Typed measurement methods adapter

Дата: 2026-06-25.

## Область

Методики измерений на сайте представлены не одним реестром, а семейством близких CMM-реестров:

| Registry ID | Раздел | Item type |
| --- | --- | --- |
| `6` | Первичные референтные методики (методы) измерений | `foei:CMM2_type` |
| `8` | Референтные методики (методы) измерений | `foei:CMM3_type` |
| `16` | Аттестованные методики (методы) измерений | `foei:CMM1_type` |

Их схемы почти совпадают: ключевой номер методики хранится в `foei:NumRegCMM`, название - в `foei:NameCMM`, статус - в `foei:StatusCMM`, аттестующая организация - в `foei:NameCerOrgCMM`, разработчик - в `foei:NameDevMetdCMM`.

## Отличия

Registry `6` содержит поля приказа:

- `foei:NumberOrderMetCMM`;
- `foei:DataOrderMetCMM`;
- `foei:DocOrderMetCMM`.

Registry `8` содержит близкие поля приказа с другим префиксом:

- `foei:NumberRefOrderMetCMM`;
- `foei:DataRefOrderMetCMM`;
- `foei:DocRefOrderMetCMM`.

Registry `16` таких order-полей в metadata snapshot не содержит.

## Вложения

Общие attachment-поля семейства реестров методик измерений:

- `foei:DocCMM` - текст методики;
- `foei:RMDocCMM` - пределы измерений;
- `foei:ECDocCMM` - характеристика погрешности;
- `foei:CertValCMM` - свидетельство об аттестации;
- `foei:MetodDocCMM` - информационные данные.

Для order-документов используются:

- `foei:DocOrderMetCMM` в registry `6`;
- `foei:DocRefOrderMetCMM` в registry `8`.

В коде вложения нормализованы через generic `arshin.AttachmentRef` и `RegistryProperty.Attachments()`.

## Связи

Связи методик лежат в `foei:CMM1Relation_assoc` как `CHILD_OBJECT`.
Проверенный пример для registry `16`, item `1404476`:

- тип связи: `Заменяет`;
- связанная методика: `ФР.1.31.2013.16016`;
- link: `/registry/16/items/286196`.

Typed-модель `methods.Relation` сохраняет нормализованный тип связи, номер связанной методики и ссылку `arshin.RegistryItemRef`, если она есть.

## Реализация

Пакет `public/registries/methods` теперь покрывает registry `6`, `8`, `16`.

- `methods.New(client)` оставлен как совместимый дефолт для registry `16`.
- `methods.NewForRegistry(client, registryID)` выбирает конкретный CMM-реестр и валидирует `6`, `8`, `16`.
- `methods.NewForSpec(client, spec)` выбирает реестр через `RegistrySpec`.
- `methods.SpecForRegistry` и `methods.SpecForItemType` возвращают описание конкретного реестра методик измерений.
- `GetByRegistryNumber` ищет запись через `/registry/{id}/data`.
- `GetByRegistryNumberDetailed` дополнительно получает sectioned карточку `/items/{itemId}/data`, чтобы заполнить вложения и связи из карточки.
- `MapRecord` маппит плоский `RegistryRecord`.
- `MapItem` маппит sectioned `RegistryItem`.

Typed-модель сохраняет `RawProperties`, `RawRecord` и `RawItem`, чтобы изменения схемы сайта не приводили к потере данных.

Registry-specific константы для всех известных разделов хранятся в общем пакете `registries`:

- `registries.CMM1` - registry `16`;
- `registries.CMM2` - registry `6`;
- `registries.CMM3` - registry `8`.

Typed-пакет `public/registries/methods` использует эти specs как вход в общий маппер методик измерений.
Отдельные пакеты создаются только для реестров, где появляется полноценный typed API.
