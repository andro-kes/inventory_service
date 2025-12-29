# Inventory Service

Сервис для хранения информации о товарах и управления их доступностью по gRPC.

## Кратко
- gRPC API: CRUD для товаров (Product).
- PostgreSQL + `pgxpool` с автоконфигурацией пула.
- Chainable SQL builder (`internal/repo/builder`) с плейсхолдерами `$n`.
- Генерация protobuf: скрипт `./proto/make_proto.sh`.
- Логирование: `zap` (консоль по умолчанию).

## Требования
- Go 1.21+
- PostgreSQL
- protoc + плагины `protoc-gen-go`, `protoc-gen-go-grpc` (для Regeneration protobuf)
- Опционально: `grpcurl`/`evans` для ручного вызова gRPC

## Переменные окружения
| Переменная  | Описание                          | Обязательна | Пример                              |
|-------------|-----------------------------------|-------------|-------------------------------------|
| `DB_URL`    | DSN PostgreSQL                    | да          | `postgres://user:pass@localhost:5432/inventory?sslmode=disable` |
| `GRPC_ADDR` | Адрес gRPC-сервера                | да          | `:50051`                            |

Пул соединений (`pgxpool`):
- `MaxConns=20`, `MinConns=2`
- `MaxConnLifetime=30m`, `HealthCheckPeriod=1m`

## Запуск
### Локально
```bash
go mod download
DB_URL=... GRPC_ADDR=:50051 go run ./cmd/server
```

### Тесты
```bash
go test ./...
```

## gRPC API (proto)
Файл: [`proto/inventory.proto`](proto/inventory.proto)

Сервис `InventoryService`:
- `ListProducts(ListRequest) returns (ListResponse)`
- `GetProduct(GetRequest) returns (GetResponse)`
- `CreateProduct(CreateRequest) returns (CreateResponse)`
- `UpdateProduct(UpdateRequest) returns (UpdateResponse)` — частичное обновление через `FieldMask`
- `DeleteProduct(DeleteRequest) returns (DeleteResponse)`

Структура `Product`:
- `id, name, description, price, quantity, tags[], available, created_at, updated_at`

### Пример вызовов через grpcurl
```bash
# Health-check отсутствует; используем любой метод
grpcurl -plaintext -d '{}' localhost:50051 inventory.InventoryService.ListProducts

grpcurl -plaintext -d '{
  "product": {
    "name": "Ноутбук",
    "description": "14\"",
    "price": 79990,
    "quantity": 10,
    "tags": ["electronics", "laptop"],
    "available": true
  }
}' localhost:50051 inventory.InventoryService.CreateProduct
```

### Фильтрация и сортировка
- `ListRequest.filter`: строка, используется в условии `tags @> ARRAY[?]::text[]`
- `ListRequest.order_by`: поддерживаются `price`, `price DESC|ASC`, `created_at`, `created_at DESC|ASC`
- Пагинация: `prev_size` (offset), `page_size` (limit)

## Структура проекта (основное)
```
cmd/server/main.go       # входная точка, gRPC server, init logger + DB
internal/logger          # zap-конфиг с ротацией (опционально)
internal/repo/builder    # SQL builder (SELECT/INSERT/UPDATE/DELETE)
internal/repo            # доступ к БД (products)
internal/services        # бизнес-логика (ProductService)
internal/rpc             # gRPC handlers
proto/                   # protobuf схемы и сгенерированные go-файлы
```

## Работа с protobuf
Сгенерировать заново (из корня проекта):
```bash
chmod +x proto/make_proto.sh
./proto/make_proto.sh
```

## SQL Builder
- Используйте `?` в where/set, билдер сам пронумерует как `$1, $2, ...`.
- Документация и примеры: [`internal/repo/builder/README.md`](internal/repo/builder/README.md).

## Логирование
По умолчанию: уровень `debug`, формат `console`, вывод в stdout (`cmd/server/main.go`). При необходимости настройте `internal/logger.Config` (JSON, ротация, файлы).

## TODO
- [] Добавить health-check endpoint/метод.
- [] Добавить пример docker-compose и миграций под PostgreSQL.
- [] Описать схемы БД (DDL) и реальный фильтр по тегам.
- [] Добавить секцию об авторизации/ACL (если потребуется).
- [] Добавить тесты.