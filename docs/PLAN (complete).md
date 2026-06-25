# PLAN.md

План исследования и реализации Go-клиента для приватного web API сайта ФГИС Аршин.

## Цель

Построить библиотеку для работы не только с реестром `16`, а со всеми доступными на сайте реестрами ФГИС Аршин, которые обслуживаются внутренним API портала:

```text
https://fgis.gost.ru/fundmetrology/api/
```

Реестр `16` (`Аттестованные методики (методы) измерений`) используем как первый проверенный пример и fixture для отладки, но архитектура должна быть общей:

- общий транспортный слой;
- общий клиент реестров;
- discovery доступных реестров и их метаданных;
- generic-представление записей и свойств;
- typed-адаптеры для отдельных важных реестров поверх generic API.

Важно явно документировать, что это не официальный публичный `eapi`, а API, которым пользуется web-интерфейс сайта. Схема и поведение могут измениться без предупреждения.

## Уже известный рабочий endpoint

Список данных реестра:

```text
GET /fundmetrology/api/registry/{registryId}/data
```

Известные параметры:

- `pageNumber` - номер страницы, начинается с `1`.
- `pageSize` - размер страницы.
- `orgID` - организация; для публичных страниц обычно `CURRENT_ORG`.
- `filterBy` - имя поля для фильтрации.
- `filterValues` - значение фильтра.
- `sortBy` - имя поля для сортировки.
- `sortValues` - направление сортировки.

Фильтры и сортировки передаются повторяющимися парами:

```text
filterBy={fieldName}&filterValues={value}
sortBy={fieldName}&sortValues={direction}
```

Проверенный пример для реестра `16`:

```text
/fundmetrology/api/registry/16/data?pageNumber=1&pageSize=5&orgID=CURRENT_ORG&filterBy=foei%3ANumRegCMM&filterValues=%D0%A4%D0%A0.1.31.2022.44733
```

## Исследование

### 1. Каталог доступных реестров

Первичный каталог похожих на registry `16` разделов собран в `REGISTRY.md`.
На его основе на старте исследуем следующие web-реестры:

| ID | Раздел |
| --- | --- |
| `1` | Нормативные правовые акты Российской Федерации |
| `2` | Шифры калибровочных клейм |
| `3` | Стандарты государственной системы обеспечения единства измерений |
| `6` | Первичные референтные методики (методы) измерений |
| `7` | Международные договоры |
| `8` | Референтные методики (методы) измерений |
| `9` | Международные сличения |
| `10` | Информация и данные ГСССД |
| `11` | Эталоны единиц величин |
| `12` | Государственные первичные эталоны Российской Федерации |
| `14` | Международные документы |
| `16` | Аттестованные методики (методы) измерений |
| `17` | Сведения об отнесении технических средств к средствам измерений |
| `18` | Перечень измерений, относящихся к сфере государственного регулирования |
| `19` | Утвержденные типы стандартных образцов |
| `47` | Уведомления об осуществлении деятельности по производству эталонов единиц величин, стандартных образцов и средств измерений |

Статус на 2026-06-25: все перечисленные разделы успешно отвечают через `/fundmetrology/api/registry/{id}/data`.
Машинная сводка сохранена в `docs/research/registry-survey.json`, краткий отчет - в `docs/research/registry-survey.md`.

Подтвержденные типы записей:

| ID | Type | Total count |
| --- | --- | ---: |
| `1` | `foei:ND_type` | 104 |
| `2` | `foei:SCM_type` | 783 |
| `3` | `foei:GSI_type` | 514 |
| `6` | `foei:CMM2_type` | 17 |
| `7` | `foei:MDG_type` | 12 |
| `8` | `foei:CMM3_type` | 7 |
| `9` | `foei:IC_type` | 437 |
| `10` | `foei:SSD_type` | 483 |
| `11` | `foei:SU_type` | 117537 |
| `12` | `foei:GPS_type` | 221 |
| `14` | `foei:MD_type` | 198 |
| `16` | `foei:CMM1_type` | 52635 |
| `17` | `foei:TSSI_type` | 1742 |
| `18` | `foei:EPI_type` | 1 |
| `19` | `foei:UTSO_type` | 11268 |
| `47` | `gost:p1wfRequestType4` | 4878 |

Для каждого раздела нужно проверить:

- открывается ли страница `/fundmetrology/registry/{id}`;
- отвечает ли endpoint `/fundmetrology/api/registry/{id}/data`;
- совпадает ли формат ответа с registry `16`;
- какой `type` у записей;
- какие поля присутствуют в `properties[]`;
- какие поля реально фильтруются и сортируются;
- есть ли вложения и связи;
- совпадает ли название в `REGISTRY.md` с названием, возвращаемым сайтом.

Дополнительно нужно найти полный машинный каталог реестров, если он есть:

- найти endpoint, который отдает список реестров или структуру меню;
- если отдельного endpoint нет, извлечь список из HTML/JS bundle сайта;
- для каждого реестра сохранить:
  - `id`;
  - название;
  - URL страницы;
  - основной тип записи, например `foei:CMM1_type`;
  - доступность публичной страницы;
  - наличие данных через `/api/registry/{id}/data`.

Результат исследования сохранить как fixture/документацию, чтобы тесты не зависели от постоянного обхода сайта.

### 2. Метаданные реестра

Статус на 2026-06-25: endpoints метаданных найдены во frontend bundle и проверены live.
Машинная сводка сохранена в `docs/research/registry-metadata.json`, краткий отчет - в `docs/research/registry-metadata.md`.

Подтвержденные endpoints:

```text
GET /fundmetrology/api/registry/{registryId}
GET /fundmetrology/api/registry/{registryId}/fields
GET /fundmetrology/api/registry/{registryId}/displayfields
GET /fundmetrology/api/fieldtypes
```

Назначение:

- `/registry/{registryId}` - метаданные реестра: id, title, Alfresco type, nodeRef, timestamps, deleted, writable.
- `/fields` - полный набор известных полей реестра.
- `/displayfields` - поля, которые сайт использует как колонки таблицы.
- `/fieldtypes` - глобальная классификация типов полей: link/attach/id.

Для нескольких реестров дополнительно проверить:

- описание колонок таблицы;
- доступные фильтры;
- типы полей;
- варианты enum/constraint;
- сортируемые поля;
- связи и вложения;
- настройки формы карточки.

Оставшиеся вопросы:

- где раскрываются значения `LIST` constraints, например `foei:ObjectStatus_list`;
- есть ли отдельные настройки фильтруемости/сортируемости или сайт считает фильтруемыми все display/full fields;
- есть ли отдельная мета-структура для формы карточки записи;
- чем объясняется отличие metadata type и item type у registry `47`.

Fallback, если endpoint метаданных изменится:

- ответа `/registry/{id}/data`;
- frontend bundle;
- HTML страницы реестра;
- фактического набора `properties[]` в нескольких записях.

### 3. Поведение generic endpoint данных

Статус на 2026-06-25: базовое поведение `/registry/{id}/data` проверено live.
Машинная сводка сохранена в `docs/research/registry-data-behavior.json`, краткий отчет - в `docs/research/registry-data-behavior.md`.

Подтверждено:

- успешный ответ имеет корневые поля `status`, `result`, `message`, `trace`;
- `result` содержит `totalCount`, `currentPage`, `pageSize`, `items`;
- пустой результат фильтрации - успешный ответ с `totalCount: 0`;
- несколько фильтров через повторяющиеся `filterBy/filterValues` работают;
- сортировка через `sortBy/sortValues` работает для `asc` и `desc`;
- пагинация 1-based, `currentPage` возвращает запрошенную страницу;
- `pageSize=100` и `pageSize=500` работают для registry `16`;
- `pageSize=1000` в отдельной пробе ушел в timeout, поэтому дефолты должны быть консервативными;
- в sampled list responses поле `values` было `null`, данные находятся в `properties[]`;
- неверный registry id `999999` вернул HTTP `500 Internal Server Error`, а не аккуратный `404`;
- item type надо хранить отдельно от metadata type: у registry `47` metadata type `gost:p1wfRequest`, а item type `gost:p1wfRequestType4`.

Для разных реестров проверить:

- расширить live-проверку на больше registry id, если появятся подозрения на отличие формата;
- проверить несколько сортировок одновременно;
- проверить фильтрацию по `LIST`, `date`, `LINK_INTERNAL`, `ATTACH` полям;
- проверить поведение слишком большой страницы без увеличения клиентского timeout;
- проверить отрицательные/нулевые `pageNumber` и `pageSize`.

### 4. Карточки записей

Статус на 2026-06-25: endpoints карточки записи найдены во frontend bundle и проверены live на записи registry `16`, item `1404476`.
Краткий отчет сохранен в `docs/research/registry-item-details.md`, fixtures - в `testdata/item-1404476-data.json` и `testdata/item-1404476-plaindata.json`.

Подтвержденные endpoints:

```text
GET /fundmetrology/api/registry/{registryId}/items/{itemId}/data
GET /fundmetrology/api/registry/{registryId}/items/{itemId}/plaindata
```

Наблюдения:

- `/items/{itemId}/data` возвращает секционную карточку с `sections[].fields[]`;
- `/items/{itemId}/plaindata` возвращает плоский `RegistryRecord` с `properties[]`;
- sectioned `/data` лучше подходит для отображения карточки и исследования вложений/связей;
- `/plaindata` лучше подходит для typed-мапперов и машинной обработки;
- в карточке `1404476` найдены `LINK_INTERNAL`, `ATTACH` и `CHILD_OBJECT` поля;
- `RegistryProperty.link` может быть строкой, массивом или `null`, поэтому пока хранится как raw `any`.

Оставшиеся проверки:

- проверить карточки нескольких других реестров, особенно `47`, `11`, `19`;
- проверить поведение неверного item id;
- проверить, бывают ли карточки без `sections`;
- проверить, всегда ли `/plaindata` совпадает с list/data shape;
- проверить права/permissions на карточках, если сервер начнет их отдавать.

### 5. Вложения

Статус на 2026-06-25: механизм скачивания вложений проверен live на поле `foei:RMDocCMM` записи registry `16`, item `1404476`.
Краткий отчет сохранен в `docs/research/registry-attachments.md`.

Подтвержденный endpoint:

```text
GET /fundmetrology/api/downloadfile/{fileId}
```

Наблюдения:

- `ATTACH` поле содержит `value` с именем файла, `link` с `/api/downloadfile/{fileId}` и `mime`;
- успешная загрузка вернула `Content-Type: application/octet-stream`, хотя поле содержало `mime: application/pdf`;
- имя файла в `Content-Disposition` percent-encoded;
- отсутствующий file id вернул HTTP `500` с JSON body формата `status/result/message/trace`;
- публичное скачивание сработало без предварительной cookie-сессии, но сервер выставил `session-cookie`.

Реализовано:

- `DownloadFile(ctx, fileID)`;
- `DownloadFileLink(ctx, link)`;
- `DownloadedFile`;
- `arshin-probe -download`.

Оставшиеся проверки:

- проверить несколько других `ATTACH` полей и реестров;
- проверить multi-file attachments;
- проверить, бывает ли `mime` массивом;
- проверить поведение приватных/ограниченных файлов, если такие встретятся;
- добавить typed helpers для извлечения attachment descriptors из `RegistryProperty`.

### 6. Связи между записями

Статус на 2026-06-25: связи исследованы live на карточке registry `16`, item `1404476`.
Машинная сводка сохранена в `docs/research/registry-relations.json`, краткий отчет - в `docs/research/registry-relations.md`.

Подтвержденные формы связей:

- `LINK_INTERNAL` с link формата `/registry/{registryId}/items/{itemId}`;
- `CHILD_OBJECT` с массивом nested objects в `value`;
- nested child object содержит `id`, `type`, `nodeRef`, `columns`, `fields`;
- внутри child object могут быть поля `LINK_INTERNAL`, по которым можно перейти к полной записи generic-методом.

Подтвержденный пример:

- `foei:CMM1Relation_assoc`;
- child type: `foei:CMM1Relation`;
- relation type: `Заменяет`;
- related method: `ФР.1.31.2013.16016`;
- related link: `/registry/16/items/286196`;
- `/registry/16/items/286196/plaindata` успешно возвращает связанную запись.

Реализовано:

- `ParseRegistryItemLink`;
- `RegistryProperty.LinkStrings`;
- `RegistryProperty.ItemRefs`;
- `RegistryProperty.ChildObjects`;
- `ChildObject.FieldsByName`;
- `ChildObject.ItemRefs`;
- `arshin-probe -item ... -relations`.

Оставшиеся проверки:

- проверить связи в других реестрах, особенно `19`, `47`, `11`;
- проверить, бывают ли одиночные child object вместо массива;
- решить, добавлять ли скрытые/reference registry `13`, `20`, `23`, `45` в статический каталог;
- проверить, есть ли отдельные endpoints раскрытия relations помимо embedded child objects.

### 7. Сессия, заголовки и защита

Статус на 2026-06-25: session/header/cookie поведение проверено live.
Машинная сводка сохранена в `docs/research/registry-session-security.json`, краткий отчет - в `docs/research/registry-session-security.md`.

Подтверждено для публичных GET endpoints:

- `/registry/16/data` работает без предварительной cookie, без Referer и без browser-like headers;
- browser-like headers также работают;
- HTML-страница `/fundmetrology/registry/16` выставляет `session-cookie`;
- HTML-страница содержит `meta name='hmac-token-name' content='Ajax-Token'`;
- GET data/item/download endpoints работают после preflight с cookie jar;
- `Ajax-Token` не потребовался для проверенного GET data request;
- invalid registry id вернул HTTP `500`, но JSON body содержал `status: 404`, `message: "not found"`;
- HTML-страница с `Accept: application/json` все равно возвращает HTML со статусом `200`, поэтому проверка `Content-Type` обязательна.

Реализовано:

- `arshin-probe -session`;
- общий JSON decoder теперь пытается распарсить non-2xx JSON body как `APIError`.

Не проверено намеренно:

- POST/PUT/DELETE endpoints, потому что безопасная публичная мутационная ручка не нужна для текущей библиотеки;
- агрессивный rate-limit, чтобы не стрессовать публичный сервис.

В реализации по умолчанию использовать осторожный профиль:

- cookie jar;
- `Referer: https://fgis.gost.ru/fundmetrology/registry/{id}`;
- `Accept: application/json, text/plain, */*`;
- `Accept-Language: ru-RU,ru;q=0.9`;
- явный `User-Agent`;
- rate limit не быстрее 1 запроса в секунду.

### 8. Typed-исследование отдельных реестров

Статус на 2026-06-25: первый typed-адаптер расширен с registry `16` до семейства реестров методик измерений `6`, `8`, `16`.
Краткий отчет сохранен в `docs/research/typed-methods.md`.

Первое typed-семейство:

- registry `6`: `Первичные референтные методики (методы) измерений`, item type `foei:CMM2_type`;
- registry `8`: `Референтные методики (методы) измерений`, item type `foei:CMM3_type`;
- registry `16`: `Аттестованные методики (методы) измерений`, item type `foei:CMM1_type`;
- проверенная запись: `ФР.1.31.2022.44733`.

Реализовано:

- `public/registries/methods` поддерживает `6`, `8`, `16`;
- общий package `registries` содержит `Spec` для всех 16 известных web-реестров;
- общий `registries.Spec`, named specs `registries.ND`, `registries.CMM1`, `registries.SU`, `registries.UTSO` и т.д.;
- lookup helpers `registries.Known`, `registries.SpecForRegistry`, `registries.SpecForItemType`;
- Methods-specific `RegistrySpec`, `SpecForRegistry`, `SpecForItemType`, `SupportedRegistries`;
- `methods.New(client)` оставлен как совместимый дефолт для registry `16`;
- `methods.NewForRegistry(client, registryID)` валидирует registry id семейства методик измерений;
- `methods.NewForSpec(client, spec)` создает клиент из registry-specific описания;
- `MapRecord` маппит плоский `RegistryRecord`;
- `MapItem` маппит sectioned `RegistryItem`;
- `GetByRegistryNumberDetailed` получает карточку записи и заполняет вложения/связи;
- generic `AttachmentRef`, `ParseDownloadFileLink`, `RegistryProperty.Attachments`;
- Methods-модель сохраняет `RawProperties`, `RawRecord`, `RawItem`.

Ключевые поля семейства реестров методик измерений:

- `foei:NumRegCMM` - номер в реестре;
- `foei:NameCMM` - наименование документа;
- `foei:StatusCMM` - статус;
- `foei:DataCertCMM`, `foei:NumCertValCMM` - данные свидетельства;
- `foei:NameDevMetdCMM` - разработчик;
- `foei:NameCerOrgCMM` - аттестующая организация;
- `foei:MeasuredVarCMM`, `foei:RangeMeasuCMM`, `foei:ErrorCharactCMM`;
- `foei:DocCMM`, `foei:RMDocCMM`, `foei:ECDocCMM`, `foei:CertValCMM`, `foei:MetodDocCMM` - вложения;
- `foei:CMM1Relation_assoc` - связи.

Отличия:

- registry `6` использует order-поля `foei:NumberOrderMetCMM`, `foei:DataOrderMetCMM`, `foei:DocOrderMetCMM`;
- registry `8` использует `foei:NumberRefOrderMetCMM`, `foei:DataRefOrderMetCMM`, `foei:DocRefOrderMetCMM`;
- registry `16` этих order-полей в metadata snapshot не содержит.

Для каждого typed-реестра нужно:

- собрать набор `property.name -> title/type/multiple/example`;
- определить ключевые поля;
- написать mapper из generic `RegistryRecord` в typed struct;
- сохранить raw properties для совместимости с будущими изменениями схемы.

Оставшиеся проверки:

- собрать live fixtures для registry `6` и `8`, чтобы проверить order-поля на реальных записях;
- проверить, одинаково ли устроены relations в `6` и `8`;
- решить, нужны ли отдельные typed-пакеты для `11`, `19`, `47` или достаточно generic API.

## Архитектура библиотеки

### Пакеты

Предварительная структура:

```text
.
├── public
│   ├── client
│   │   ├── client.go
│   │   ├── errors.go
│   │   ├── registry.go
│   │   ├── record.go
│   │   ├── property.go
│   │   └── ratelimit.go
│   └── registries
│       ├── spec.go
│       └── methods
│           ├── client.go
│           ├── model.go
│           ├── spec.go
│           └── mapper.go
├── testdata
├── docs
│   └── research
└── cmd
    └── arshin-probe
```

Публичный generic-клиент живет в `public/client`. Каталог реестров и typed-адаптеры живут в `public/registries/...` и зависят от generic-клиента, но не дублируют транспорт и HTTP-логику.

### Generic API

Базовые сущности:

- `Client`;
- `RegistryID`;
- `RegistryInfo`;
- `ListRegistries`;
- `ListRegistryData`;
- `GetRegistryRecord`;
- `DownloadAttachment`;
- `RegistryRecord`;
- `RegistryProperty`;
- `RegistryQuery`;
- `Filter`;
- `Sort`.

Пример целевого публичного интерфейса:

```go
client, err := arshinclient.NewClient()
if err != nil {
    return err
}

records, err := client.ListRegistryData(ctx, "16", arshinclient.RegistryQuery{
    PageNumber: 1,
    PageSize:   20,
    OrgID:      "CURRENT_ORG",
    Filters: []arshinclient.Filter{
        {Field: "foei:NumRegCMM", Value: "ФР.1.31.2022.44733"},
    },
})
```

### Typed API поверх generic API

Typed-адаптеры не должны дублировать транспорт:

```go
methodsClient := methods.New(client)

method, err := methodsClient.GetByRegistryNumber(ctx, "ФР.1.31.2022.44733")
```

Typed-модели должны хранить:

- нормализованные поля;
- id записи;
- тип записи;
- raw properties;
- raw record при необходимости.

## Реализация по этапам

### Этап 1. Базовый проект

- Создать Go module.
- Добавить README с предупреждением про private web API.
- Добавить базовый `Client`.
- Добавить настройку base URL, HTTP client, user agent, cookie jar.
- Добавить общий JSON transport и типы ошибок.

### Этап 2. Generic registry data

- Реализовать `/registry/{id}/data`.
- Поддержать пагинацию.
- Поддержать несколько фильтров.
- Поддержать несколько сортировок.
- Реализовать generic-типы ответа.
- Добавить unit-тесты query encoding.

### Этап 3. Fixtures и probe CLI

- Добавить `cmd/arshin-probe`.
- Команды:
  - получить страницу реестра;
  - запросить данные реестра;
  - сохранить JSON fixture;
  - вывести список найденных property names.
- Использовать CLI для всех первичных registry id из `REGISTRY.md`: `1`, `2`, `3`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `14`, `16`, `17`, `18`, `19`, `47`.

### Этап 4. Discovery реестров

- Перенести список из `REGISTRY.md` в код как начальный статический каталог.
- Найти способ получить список всех реестров с сайта.
- Реализовать `ListRegistries`, если есть стабильный endpoint.
- Если стабильного endpoint нет, сделать best-effort discovery и documented fixtures.
- Добавить тесты на парсинг discovery fixtures.

### Этап 5. Детали, вложения, связи

- Найти и реализовать endpoint карточки записи, если он есть.
- Найти и реализовать download endpoint.
- Описать association-поля.
- Добавить generic helpers для вложений и связей.

### Этап 6. Typed measurement methods registries `6`, `8`, `16`

- Реализовать typed-пакет для методик измерений.
- Замаппить ключевые поля:
  - `foei:NumRegCMM`;
  - `foei:number`;
  - `foei:NameCMM`;
  - `foei:StatusCMM`;
  - `foei:status`;
  - `foei:date`;
  - `foei:DataCertCMM`;
  - `foei:NumCertValCMM`;
  - `foei:NameDevMetdCMM`;
  - `foei:NameCerOrgCMM`;
  - `foei:MeasuredVarCMM`;
  - `foei:RangeMeasuCMM`;
  - `foei:ErrorCharactCMM`;
  - `foei:TypeMetCMM`;
  - `foei:TypeMeasuCMM`;
  - `foei:SFCMM`;
  - `foei:DocCMM`;
  - `foei:RMDocCMM`;
  - `foei:ECDocCMM`;
  - `foei:CertValCMM`;
  - `foei:MetodDocCMM`;
  - `foei:CMM1Relation_assoc`.
- Добавить fixture на запись `ФР.1.31.2022.44733`.

### Этап 7. Остальные typed-адаптеры

Добавлять typed-пакеты только для реестров, где есть реальная потребность и понятная схема. Для остальных оставить generic API.

Критерии добавления typed-адаптера:

- часто используемый реестр;
- есть стабильные ключевые поля;
- есть понятный сценарий поиска;
- generic `properties[]` слишком неудобен для потребителя.

### Этап 8. Документация и примеры

- Документировать generic usage.
- Документировать registry discovery.
- Документировать typed usage для CMM registry `6`, `8`, `16`.
- Описать ограничения и совместимость.
- Добавить примеры обработки ошибок и rate limit.

## Тестирование

### Unit-тесты

- Кодирование query-параметров.
- Парсинг generic JSON.
- Маппинг `properties[]`.
- Typed mapping registry `16`.
- Ошибки API.
- HTML вместо JSON.

### Fixture-тесты

- Ответ registry `16` по `ФР.1.31.2022.44733`.
- Несколько ответов других реестров.
- Пустой результат.
- Ответ с вложением.
- Ответ со связью.

### Live-тесты

Live-тесты запускать только явно, например через env:

```text
ARSHIN_LIVE_TESTS=1 go test ./...
```

Live-тесты должны:

- быть медленными;
- соблюдать rate limit;
- не полагаться на большое число запросов;
- пропускаться по умолчанию.

## Риски

- API не является публично документированным.
- Схема полей может меняться.
- Некоторые реестры могут иметь отличающийся формат.
- Доступ к деталям или файлам может требовать cookie/token.
- Возможны ограничения по частоте запросов.
- Часть метаданных может быть доступна только из frontend bundle.

## Ближайший следующий шаг

Начать с исследовательского `arshin-probe`, который умеет:

1. выполнять generic-запрос `/registry/{id}/data`;
2. сохранять сырой JSON;
3. печатать список `properties[].name/title/type`;
4. работать с любым `registryId`;
5. использовать осторожные browser-like headers, cookie jar и rate limit.

После этого собрать fixtures для `16` и еще нескольких реестров, чтобы подтвердить, насколько универсален формат API.
