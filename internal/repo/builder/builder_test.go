package builder

import (
	"testing"
	"time"
)

// TestSelectBasic tests a basic SELECT query with all columns.
func TestSelectBasic(t *testing.T) {
	query, args := NewSQLBuilder().
		Select().
		From("users").
		Build()

	expected := "SELECT * FROM users"
	if query != expected {
		t.Errorf("Expected query: %s, got: %s", expected, query)
	}
	if len(args) != 0 {
		t.Errorf("Expected 0 args, got: %d", len(args))
	}
}

// TestSelectWithColumns tests a SELECT query with specific columns.
func TestSelectWithColumns(t *testing.T) {
	query, args := NewSQLBuilder().
		Select("id", "name", "email").
		From("users").
		Build()

	expected := "SELECT id, name, email FROM users"
	if query != expected {
		t.Errorf("Expected query: %s, got: %s", expected, query)
	}
	if len(args) != 0 {
		t.Errorf("Expected 0 args, got: %d", len(args))
	}
}

// TestSelectWithWhere tests a SELECT query with WHERE clause.
func TestSelectWithWhere(t *testing.T) {
	query, args := NewSQLBuilder().
		Select("id", "name").
		From("users").
		Where("age > ?", 18).
		Build()

	expected := "SELECT id, name FROM users WHERE age > ?"
	if query != expected {
		t.Errorf("Expected query: %s, got: %s", expected, query)
	}
	if len(args) != 1 || args[0] != 18 {
		t.Errorf("Expected args: [18], got: %v", args)
	}
}

// TestSelectWithMultipleWhere tests a SELECT query with multiple WHERE clauses.
func TestSelectWithMultipleWhere(t *testing.T) {
	query, args := NewSQLBuilder().
		Select("id", "name").
		From("users").
		Where("age > ?", 18).
		Where("status = ?", "active").
		Build()

	expected := "SELECT id, name FROM users WHERE age > ? AND status = ?"
	if query != expected {
		t.Errorf("Expected query: %s, got: %s", expected, query)
	}
	if len(args) != 2 {
		t.Errorf("Expected 2 args, got: %d", len(args))
	}
	if args[0] != 18 || args[1] != "active" {
		t.Errorf("Expected args: [18, active], got: %v", args)
	}
}

// TestSelectWithOrderBy tests a SELECT query with ORDER BY clause.
func TestSelectWithOrderBy(t *testing.T) {
	query, args := NewSQLBuilder().
		Select("id", "name", "created_at").
		From("users").
		OrderBy("created_at DESC").
		Build()

	expected := "SELECT id, name, created_at FROM users ORDER BY created_at DESC"
	if query != expected {
		t.Errorf("Expected query: %s, got: %s", expected, query)
	}
	if len(args) != 0 {
		t.Errorf("Expected 0 args, got: %d", len(args))
	}
}

// TestSelectWithLimit tests a SELECT query with LIMIT clause.
func TestSelectWithLimit(t *testing.T) {
	query, args := NewSQLBuilder().
		Select("id", "name").
		From("users").
		Limit(10).
		Build()

	expected := "SELECT id, name FROM users LIMIT 10"
	if query != expected {
		t.Errorf("Expected query: %s, got: %s", expected, query)
	}
	if len(args) != 0 {
		t.Errorf("Expected 0 args, got: %d", len(args))
	}
}

// TestSelectWithOffset tests a SELECT query with OFFSET clause.
func TestSelectWithOffset(t *testing.T) {
	query, args := NewSQLBuilder().
		Select("id", "name").
		From("users").
		Offset(20).
		Build()

	expected := "SELECT id, name FROM users OFFSET 20"
	if query != expected {
		t.Errorf("Expected query: %s, got: %s", expected, query)
	}
	if len(args) != 0 {
		t.Errorf("Expected 0 args, got: %d", len(args))
	}
}

// TestSelectComplex tests a complex SELECT query with all clauses.
func TestSelectComplex(t *testing.T) {
	query, args := NewSQLBuilder().
		Select("id", "name", "email", "created_at").
		From("users").
		Where("age > ?", 18).
		Where("status = ?", "active").
		OrderBy("created_at DESC").
		Limit(10).
		Offset(20).
		Build()

	expected := "SELECT id, name, email, created_at FROM users WHERE age > ? AND status = ? ORDER BY created_at DESC LIMIT 10 OFFSET 20"
	if query != expected {
		t.Errorf("Expected query: %s, got: %s", expected, query)
	}
	if len(args) != 2 {
		t.Errorf("Expected 2 args, got: %d", len(args))
	}
}

// TestInsertBasic tests a basic INSERT query.
func TestInsertBasic(t *testing.T) {
	query, args := NewSQLBuilder().
		Insert("users").
		Columns("name", "email", "age").
		Values("John Doe", "john@example.com", 30).
		Build()

	expected := "INSERT INTO users (name, email, age) VALUES (?, ?, ?)"
	if query != expected {
		t.Errorf("Expected query: %s, got: %s", expected, query)
	}
	if len(args) != 3 {
		t.Errorf("Expected 3 args, got: %d", len(args))
	}
	if args[0] != "John Doe" || args[1] != "john@example.com" || args[2] != 30 {
		t.Errorf("Expected args: [John Doe, john@example.com, 30], got: %v", args)
	}
}

// TestInsertWithoutColumns tests an INSERT query without specifying columns.
func TestInsertWithoutColumns(t *testing.T) {
	query, args := NewSQLBuilder().
		Insert("users").
		Values("John Doe", "john@example.com", 30).
		Build()

	expected := "INSERT INTO users VALUES (?, ?, ?)"
	if query != expected {
		t.Errorf("Expected query: %s, got: %s", expected, query)
	}
	if len(args) != 3 {
		t.Errorf("Expected 3 args, got: %d", len(args))
	}
}

// TestInsertMultipleValues tests an INSERT query with multiple value sets.
func TestInsertMultipleValues(t *testing.T) {
	query, args := NewSQLBuilder().
		Insert("products").
		Columns("name", "price").
		Values("Laptop", 999.99).
		Build()

	expected := "INSERT INTO products (name, price) VALUES (?, ?)"
	if query != expected {
		t.Errorf("Expected query: %s, got: %s", expected, query)
	}
	if len(args) != 2 {
		t.Errorf("Expected 2 args, got: %d", len(args))
	}
}

// TestUpdateBasic tests a basic UPDATE query.
func TestUpdateBasic(t *testing.T) {
	query, args := NewSQLBuilder().
		Update("users").
		Set("name = ?", "Jane Doe").
		Where("id = ?", 123).
		Build()

	expected := "UPDATE users SET name = ? WHERE id = ?"
	if query != expected {
		t.Errorf("Expected query: %s, got: %s", expected, query)
	}
	if len(args) != 2 {
		t.Errorf("Expected 2 args, got: %d", len(args))
	}
	if args[0] != "Jane Doe" || args[1] != 123 {
		t.Errorf("Expected args: [Jane Doe, 123], got: %v", args)
	}
}

// TestUpdateMultipleSets tests an UPDATE query with multiple SET clauses.
func TestUpdateMultipleSets(t *testing.T) {
	now := time.Now()
	query, args := NewSQLBuilder().
		Update("users").
		Set("name = ?", "Jane Doe").
		Set("age = ?", 31).
		Set("updated_at = ?", now).
		Where("id = ?", 123).
		Build()

	expected := "UPDATE users SET name = ?, age = ?, updated_at = ? WHERE id = ?"
	if query != expected {
		t.Errorf("Expected query: %s, got: %s", expected, query)
	}
	if len(args) != 4 {
		t.Errorf("Expected 4 args, got: %d", len(args))
	}
}

// TestUpdateWithMultipleWhere tests an UPDATE query with multiple WHERE clauses.
func TestUpdateWithMultipleWhere(t *testing.T) {
	query, args := NewSQLBuilder().
		Update("users").
		Set("status = ?", "inactive").
		Where("age < ?", 18).
		Where("last_login < ?", "2020-01-01").
		Build()

	expected := "UPDATE users SET status = ? WHERE age < ? AND last_login < ?"
	if query != expected {
		t.Errorf("Expected query: %s, got: %s", expected, query)
	}
	if len(args) != 3 {
		t.Errorf("Expected 3 args, got: %d", len(args))
	}
}

// TestUpdateWithoutWhere tests an UPDATE query without WHERE clause.
func TestUpdateWithoutWhere(t *testing.T) {
	query, args := NewSQLBuilder().
		Update("users").
		Set("status = ?", "active").
		Build()

	expected := "UPDATE users SET status = ?"
	if query != expected {
		t.Errorf("Expected query: %s, got: %s", expected, query)
	}
	if len(args) != 1 {
		t.Errorf("Expected 1 arg, got: %d", len(args))
	}
}

// TestDeleteBasic tests a basic DELETE query.
func TestDeleteBasic(t *testing.T) {
	query, args := NewSQLBuilder().
		Delete().
		From("users").
		Where("id = ?", 123).
		Build()

	expected := "DELETE FROM users WHERE id = ?"
	if query != expected {
		t.Errorf("Expected query: %s, got: %s", expected, query)
	}
	if len(args) != 1 {
		t.Errorf("Expected 1 arg, got: %d", len(args))
	}
	if args[0] != 123 {
		t.Errorf("Expected arg: 123, got: %v", args[0])
	}
}

// TestDeleteWithMultipleWhere tests a DELETE query with multiple WHERE clauses.
func TestDeleteWithMultipleWhere(t *testing.T) {
	query, args := NewSQLBuilder().
		Delete().
		From("users").
		Where("age < ?", 18).
		Where("status = ?", "inactive").
		Build()

	expected := "DELETE FROM users WHERE age < ? AND status = ?"
	if query != expected {
		t.Errorf("Expected query: %s, got: %s", expected, query)
	}
	if len(args) != 2 {
		t.Errorf("Expected 2 args, got: %d", len(args))
	}
}

// TestDeleteWithoutWhere tests a DELETE query without WHERE clause.
func TestDeleteWithoutWhere(t *testing.T) {
	query, args := NewSQLBuilder().
		Delete().
		From("users").
		Build()

	expected := "DELETE FROM users"
	if query != expected {
		t.Errorf("Expected query: %s, got: %s", expected, query)
	}
	if len(args) != 0 {
		t.Errorf("Expected 0 args, got: %d", len(args))
	}
}

// TestChainability tests that methods can be chained in any order.
func TestChainability(t *testing.T) {
	// Test chaining in different orders produces the same result
	query1, args1 := NewSQLBuilder().
		Select("id", "name").
		From("users").
		Where("status = ?", "active").
		Limit(10).
		Build()

	query2, args2 := NewSQLBuilder().
		Select("id", "name").
		Where("status = ?", "active").
		From("users").
		Limit(10).
		Build()

	if query1 != query2 {
		t.Errorf("Chaining order should not affect query. Got: %s vs %s", query1, query2)
	}
	if len(args1) != len(args2) || args1[0] != args2[0] {
		t.Errorf("Chaining order should not affect args. Got: %v vs %v", args1, args2)
	}
}

// TestEmptyBuilder tests that building without setting any values returns empty.
func TestEmptyBuilder(t *testing.T) {
	query, args := NewSQLBuilder().Build()

	if query != "" {
		t.Errorf("Expected empty query, got: %s", query)
	}
	if args != nil {
		t.Errorf("Expected nil args, got: %v", args)
	}
}

// TestSelectNoFrom tests SELECT without FROM clause.
func TestSelectNoFrom(t *testing.T) {
	query, args := NewSQLBuilder().
		Select("1", "2", "3").
		Build()

	expected := "SELECT 1, 2, 3"
	if query != expected {
		t.Errorf("Expected query: %s, got: %s", expected, query)
	}
	if len(args) != 0 {
		t.Errorf("Expected 0 args, got: %d", len(args))
	}
}

// TestMultipleWhereArgs tests WHERE clause with multiple arguments.
func TestMultipleWhereArgs(t *testing.T) {
	query, args := NewSQLBuilder().
		Select("id", "name").
		From("users").
		Where("id IN (?, ?, ?)", 1, 2, 3).
		Build()

	expected := "SELECT id, name FROM users WHERE id IN (?, ?, ?)"
	if query != expected {
		t.Errorf("Expected query: %s, got: %s", expected, query)
	}
	if len(args) != 3 {
		t.Errorf("Expected 3 args, got: %d", len(args))
	}
	if args[0] != 1 || args[1] != 2 || args[2] != 3 {
		t.Errorf("Expected args: [1, 2, 3], got: %v", args)
	}
}

// TestComplexRealWorldScenarios tests real-world complex scenarios.
func TestComplexRealWorldScenarios(t *testing.T) {
	// Scenario 1: Pagination query with filters
	query, args := NewSQLBuilder().
		Select("id", "name", "email", "created_at").
		From("users").
		Where("status = ?", "active").
		Where("role IN (?, ?)", "admin", "moderator").
		OrderBy("created_at DESC").
		Limit(20).
		Offset(40).
		Build()

	expectedQuery := "SELECT id, name, email, created_at FROM users WHERE status = ? AND role IN (?, ?) ORDER BY created_at DESC LIMIT 20 OFFSET 40"
	if query != expectedQuery {
		t.Errorf("Expected query: %s, got: %s", expectedQuery, query)
	}
	if len(args) != 3 {
		t.Errorf("Expected 3 args, got: %d", len(args))
	}

	// Scenario 2: Complex update with timestamp
	now := time.Now()
	query2, args2 := NewSQLBuilder().
		Update("products").
		Set("price = ?", 899.99).
		Set("discount = ?", 10.0).
		Set("updated_at = ?", now).
		Where("category = ?", "electronics").
		Where("stock > ?", 0).
		Build()

	expectedQuery2 := "UPDATE products SET price = ?, discount = ?, updated_at = ? WHERE category = ? AND stock > ?"
	if query2 != expectedQuery2 {
		t.Errorf("Expected query: %s, got: %s", expectedQuery2, query2)
	}
	if len(args2) != 5 {
		t.Errorf("Expected 5 args, got: %d", len(args2))
	}

	// Scenario 3: Conditional delete
	query3, args3 := NewSQLBuilder().
		Delete().
		From("sessions").
		Where("expires_at < ?", now).
		Where("user_id IS NOT NULL").
		Build()

	expectedQuery3 := "DELETE FROM sessions WHERE expires_at < ? AND user_id IS NOT NULL"
	if query3 != expectedQuery3 {
		t.Errorf("Expected query: %s, got: %s", expectedQuery3, query3)
	}
	if len(args3) != 1 {
		t.Errorf("Expected 1 arg, got: %d", len(args3))
	}
}
