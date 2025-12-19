# SQL Builder

A lightweight, chainable SQL query builder for Go that provides a fluent API for constructing SQL queries dynamically.

## Features

- **Chainable API**: Build queries using method chaining for better readability
- **Support for Basic SQL Commands**: SELECT, INSERT, UPDATE, DELETE
- **Type-safe**: Uses Go interfaces for parameter binding
- **Well-documented**: Comprehensive documentation and examples
- **Well-tested**: Extensive test coverage with real-world scenarios
- **No dependencies**: Pure Go implementation

## Installation

```bash
go get github.com/andro-kes/inventory_service/builder
```

## Usage

### Import

```go
import "github.com/andro-kes/inventory_service/internal/repo/builder"
```

### SELECT Queries

#### Basic SELECT

```go
query, args := builder.NewSQLBuilder().
    Select("id", "name", "email").
    From("users").
    Build()
// Result: SELECT id, name, email FROM users
```

#### SELECT with WHERE

```go
query, args := builder.NewSQLBuilder().
    Select("id", "name").
    From("users").
    Where("age > ?", 18).
    Where("status = ?", "active").
    Build()
// Result: SELECT id, name FROM users WHERE age > ? AND status = ?
// Args: [18, "active"]
```

#### SELECT with Ordering and Pagination

```go
query, args := builder.NewSQLBuilder().
    Select("id", "name", "created_at").
    From("products").
    Where("category = ?", "electronics").
    OrderBy("created_at DESC").
    Limit(10).
    Offset(20).
    Build()
// Result: SELECT id, name, created_at FROM products 
//         WHERE category = ? ORDER BY created_at DESC LIMIT 10 OFFSET 20
```

### INSERT Queries

#### Basic INSERT

```go
query, args := builder.NewSQLBuilder().
    Insert("users").
    Columns("name", "email", "age").
    Values("John Doe", "john@example.com", 30).
    Build()
// Result: INSERT INTO users (name, email, age) VALUES (?, ?, ?)
// Args: ["John Doe", "john@example.com", 30]
```

#### INSERT without Column Names

```go
query, args := builder.NewSQLBuilder().
    Insert("users").
    Values("John Doe", "john@example.com", 30).
    Build()
// Result: INSERT INTO users VALUES (?, ?, ?)
```

### UPDATE Queries

#### Basic UPDATE

```go
query, args := builder.NewSQLBuilder().
    Update("users").
    Set("name = ?", "Jane Doe").
    Set("age = ?", 31).
    Where("id = ?", 123).
    Build()
// Result: UPDATE users SET name = ?, age = ? WHERE id = ?
// Args: ["Jane Doe", 31, 123]
```

#### UPDATE with Multiple Conditions

```go
query, args := builder.NewSQLBuilder().
    Update("products").
    Set("price = ?", 899.99).
    Set("discount = ?", 10.0).
    Where("category = ?", "electronics").
    Where("stock > ?", 0).
    Build()
// Result: UPDATE products SET price = ?, discount = ? 
//         WHERE category = ? AND stock > ?
```

### DELETE Queries

#### Basic DELETE

```go
query, args := builder.NewSQLBuilder().
    Delete().
    From("users").
    Where("id = ?", 123).
    Build()
// Result: DELETE FROM users WHERE id = ?
// Args: [123]
```

#### DELETE with Multiple Conditions

```go
query, args := builder.NewSQLBuilder().
    Delete().
    From("users").
    Where("age < ?", 18).
    Where("status = ?", "inactive").
    Build()
// Result: DELETE FROM users WHERE age < ? AND status = ?
```

## Using with database/sql

The builder generates parameterized queries that work seamlessly with Go's `database/sql` package:

```go
import (
    "database/sql"
    "github.com/andro-kes/inventory_service/builder"
)

func getActiveUsers(db *sql.DB) ([]*User, error) {
    query, args := builder.NewSQLBuilder().
        Select("id", "name", "email").
        From("users").
        Where("status = ?", "active").
        OrderBy("created_at DESC").
        Limit(100).
        Build()
    
    rows, err := db.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    // Process rows...
    return users, nil
}
```

## Dynamic Query Building

The builder is perfect for constructing queries based on runtime conditions:

```go
func searchProducts(filters map[string]interface{}) (string, []interface{}) {
    b := builder.NewSQLBuilder().
        Select("id", "name", "price").
        From("products")
    
    // Add filters dynamically
    if category, ok := filters["category"]; ok {
        b.Where("category = ?", category)
    }
    
    if minPrice, ok := filters["min_price"]; ok {
        b.Where("price >= ?", minPrice)
    }
    
    if maxPrice, ok := filters["max_price"]; ok {
        b.Where("price <= ?", maxPrice)
    }
    
    // Add sorting if specified
    if sortBy, ok := filters["sort_by"]; ok {
        b.OrderBy(sortBy.(string))
    }
    
    return b.Build()
}
```

## API Reference

### Builder Methods

#### Query Type Methods

- `Select(columns ...string) *SQLBuilder` - Start a SELECT query
- `Insert(table string) *SQLBuilder` - Start an INSERT query
- `Update(table string) *SQLBuilder` - Start an UPDATE query
- `Delete() *SQLBuilder` - Start a DELETE query

#### Table Methods

- `From(table string) *SQLBuilder` - Specify the table name

#### Column Methods

- `Columns(columns ...string) *SQLBuilder` - Specify columns for INSERT

#### Value Methods

- `Values(values ...interface{}) *SQLBuilder` - Specify values for INSERT

#### Condition Methods

- `Where(condition string, args ...interface{}) *SQLBuilder` - Add WHERE condition (multiple calls are combined with AND)
- `Set(clause string, args ...interface{}) *SQLBuilder` - Add SET clause for UPDATE

#### Modifier Methods

- `OrderBy(column string) *SQLBuilder` - Add ORDER BY clause
- `Limit(limit int) *SQLBuilder` - Add LIMIT clause
- `Offset(offset int) *SQLBuilder` - Add OFFSET clause

#### Build Method

- `Build() (string, []interface{})` - Generate the final SQL query and arguments

## Testing

Run the tests:

```bash
go test -v ./builder/...
```

Run tests with coverage:

```bash
go test -cover ./builder/...
```

## Design Principles

1. **Modularity**: Each SQL component (SELECT, WHERE, ORDER BY, etc.) is handled by a separate method
2. **Chainability**: All methods return `*SQLBuilder` allowing for method chaining
3. **Type Safety**: Uses Go's type system and interfaces for safer query construction
4. **Simplicity**: Clean, intuitive API that follows SQL structure
5. **Flexibility**: Supports dynamic query construction based on runtime conditions

## Limitations

- Does not support JOIN operations (can be added in future versions)
- Does not support subqueries (can be added in future versions)
- Does not support UNION operations (can be added in future versions)
- Placeholders are always `?` (not database-specific like `$1`, `$2` for PostgreSQL)

## License

This is part of the inventory_service project.

## Contributing

Contributions are welcome! Please ensure:
- All tests pass
- New features include tests
- Code follows Go best practices
- Documentation is updated
