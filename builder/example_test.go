package builder_test

import (
	"fmt"
	"time"

	"github.com/andro-kes/inventory_service/builder"
)

// ExampleSQLBuilder_Select demonstrates a basic SELECT query.
func ExampleSQLBuilder_Select() {
	query, args := builder.NewSQLBuilder().
		Select("id", "name", "email").
		From("users").
		Build()

	fmt.Println(query)
	fmt.Println(args)
	// Output:
	// SELECT id, name, email FROM users
	// []
}

// ExampleSQLBuilder_Where demonstrates a SELECT query with WHERE clause.
func ExampleSQLBuilder_Where() {
	query, args := builder.NewSQLBuilder().
		Select("id", "name").
		From("users").
		Where("age > ?", 18).
		Where("status = ?", "active").
		Build()

	fmt.Println(query)
	fmt.Printf("%v\n", args)
	// Output:
	// SELECT id, name FROM users WHERE age > ? AND status = ?
	// [18 active]
}

// ExampleSQLBuilder_Limit demonstrates pagination.
func ExampleSQLBuilder_Limit() {
	query, args := builder.NewSQLBuilder().
		Select("id", "name", "created_at").
		From("products").
		Where("category = ?", "electronics").
		OrderBy("created_at DESC").
		Limit(10).
		Offset(20).
		Build()

	fmt.Println(query)
	fmt.Printf("%v\n", args)
	// Output:
	// SELECT id, name, created_at FROM products WHERE category = ? ORDER BY created_at DESC LIMIT 10 OFFSET 20
	// [electronics]
}

// ExampleSQLBuilder_Insert demonstrates an INSERT query.
func ExampleSQLBuilder_Insert() {
	query, args := builder.NewSQLBuilder().
		Insert("users").
		Columns("name", "email", "age").
		Values("John Doe", "john@example.com", 30).
		Build()

	fmt.Println(query)
	fmt.Printf("%v\n", args)
	// Output:
	// INSERT INTO users (name, email, age) VALUES (?, ?, ?)
	// [John Doe john@example.com 30]
}

// ExampleSQLBuilder_Update demonstrates an UPDATE query.
func ExampleSQLBuilder_Update() {
	query, args := builder.NewSQLBuilder().
		Update("users").
		Set("name = ?", "Jane Doe").
		Set("age = ?", 31).
		Where("id = ?", 123).
		Build()

	fmt.Println(query)
	fmt.Printf("%v\n", args)
	// Output:
	// UPDATE users SET name = ?, age = ? WHERE id = ?
	// [Jane Doe 31 123]
}

// ExampleSQLBuilder_Delete demonstrates a DELETE query.
func ExampleSQLBuilder_Delete() {
	query, args := builder.NewSQLBuilder().
		Delete().
		From("users").
		Where("id = ?", 123).
		Build()

	fmt.Println(query)
	fmt.Printf("%v\n", args)
	// Output:
	// DELETE FROM users WHERE id = ?
	// [123]
}

// ExampleSQLBuilder_Build demonstrates a complex query with multiple conditions.
func ExampleSQLBuilder_Build() {
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	query, args := builder.NewSQLBuilder().
		Select("id", "title", "author", "published_at").
		From("books").
		Where("published_at > ?", now).
		Where("status = ?", "available").
		Where("price < ?", 50.00).
		OrderBy("published_at DESC").
		Limit(5).
		Build()

	fmt.Println(query)
	fmt.Printf("Args count: %d\n", len(args))
	// Output:
	// SELECT id, title, author, published_at FROM books WHERE published_at > ? AND status = ? AND price < ? ORDER BY published_at DESC LIMIT 5
	// Args count: 3
}

// Example_usageWithDatabase demonstrates how to use the builder with database/sql.
func Example_usageWithDatabase() {
	// This is a conceptual example showing how to use the builder with database/sql
	// Note: This won't actually run without a real database connection

	// Build a query
	query, args := builder.NewSQLBuilder().
		Select("id", "name", "email").
		From("users").
		Where("status = ?", "active").
		Limit(10).
		Build()

	// Use with database/sql:
	// rows, err := db.Query(query, args...)
	// if err != nil {
	//     log.Fatal(err)
	// }
	// defer rows.Close()

	fmt.Println("Query:", query)
	fmt.Printf("Args: %v\n", args)
	// Output:
	// Query: SELECT id, name, email FROM users WHERE status = ? LIMIT 10
	// Args: [active]
}

// Example_reusableBuilder demonstrates reusing a builder for similar queries.
func Example_reusableBuilder() {
	// Create a base builder for common query parts
	baseQuery := builder.NewSQLBuilder().
		Select("id", "name", "price").
		From("products").
		Where("status = ?", "active")

	// Build query with additional filter
	query1, args1 := baseQuery.
		Where("price > ?", 100.00).
		Build()

	fmt.Println("Query 1:", query1)
	fmt.Printf("Args 1: %v\n", args1)

	// Note: Once Build() is called, the builder's state is used.
	// For completely independent queries, create new builders.

	query2, args2 := builder.NewSQLBuilder().
		Select("id", "name", "price").
		From("products").
		Where("status = ?", "active").
		Where("category = ?", "electronics").
		Build()

	fmt.Println("Query 2:", query2)
	fmt.Printf("Args 2: %v\n", args2)
	// Output:
	// Query 1: SELECT id, name, price FROM products WHERE status = ? AND price > ?
	// Args 1: [active 100]
	// Query 2: SELECT id, name, price FROM products WHERE status = ? AND category = ?
	// Args 2: [active electronics]
}

// Example_dynamicQueryBuilding demonstrates building queries dynamically based on conditions.
func Example_dynamicQueryBuilding() {
	// Simulate dynamic filters
	filters := map[string]interface{}{
		"category": "electronics",
		"minPrice": 100.00,
	}

	b := builder.NewSQLBuilder().
		Select("id", "name", "price").
		From("products")

	// Add filters dynamically
	if category, ok := filters["category"]; ok {
		b.Where("category = ?", category)
	}

	if minPrice, ok := filters["minPrice"]; ok {
		b.Where("price >= ?", minPrice)
	}

	query, args := b.Build()

	fmt.Println(query)
	fmt.Printf("%v\n", args)
	// Output:
	// SELECT id, name, price FROM products WHERE category = ? AND price >= ?
	// [electronics 100]
}
