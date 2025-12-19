package builder

import (
	"fmt"
	"strings"
)

// SQLBuilder provides a chainable API for building SQL queries.
// It supports SELECT, INSERT, UPDATE, and DELETE operations with
// a fluent interface for constructing queries dynamically.
//
// The builder generates PostgreSQL-compatible queries with positional
// placeholders ($1, $2, ...) for use with pgx and other PostgreSQL drivers.
//
// Example usage:
//
//	// SELECT query
//	query, args := NewSQLBuilder().
//		Select("id", "name", "price").
//		From("products").
//		Where("category = ?", "electronics").
//		OrderBy("price DESC").
//		Limit(10).
//		Build()
//	// Result: SELECT id, name, price FROM products WHERE category = $1 ORDER BY price DESC LIMIT 10
//
//	// INSERT query
//	query, args := NewSQLBuilder().
//		Insert("products").
//		Columns("name", "price", "category").
//		Values("Laptop", 999.99, "electronics").
//		Returning("id").
//		Build()
//	// Result: INSERT INTO products (name, price, category) VALUES ($1, $2, $3) RETURNING id
//
//	// UPDATE query
//	query, args := NewSQLBuilder().
//		Update("products").
//		Set("price = ?", 899.99).
//		Set("updated_at = ?", time.Now()).
//		Where("id = ?", 123).
//		Returning("id").
//		Build()
//	// Result: UPDATE products SET price = $1, updated_at = $2 WHERE id = $3 RETURNING id
//
//	// DELETE query
//	query, args := NewSQLBuilder().
//		Delete().
//		From("products").
//		Where("id = ?", 123).
//		Returning("id").
//		Build()
//	// Result: DELETE FROM products WHERE id = $1 RETURNING id
type SQLBuilder struct {
	queryType  string   // SELECT, INSERT, UPDATE, DELETE
	selectCols []string // Columns for SELECT
	tableName  string   // Table name
	insertCols []string // Columns for INSERT
	returning  []string
	values     []any
	setClauses []setClause
	whereConds []whereCondition
	orderByCol string
	limitVal   int
	offsetVal  int
}

type setClause struct {
	clause string
	args   []any
}

type whereCondition struct {
	condition string
	args      []any
}

// NewSQLBuilder creates a new SQLBuilder instance.
func NewSQLBuilder() *SQLBuilder {
	return &SQLBuilder{
		selectCols: make([]string, 0),
		insertCols: make([]string, 0),
		values:     make([]any, 0),
		setClauses: make([]setClause, 0),
		whereConds: make([]whereCondition, 0),
		limitVal:   -1,
		offsetVal:  -1,
	}
}

// Select specifies the columns to select in a SELECT query.
// Multiple columns can be provided as separate arguments.
//
// Example:
//
//	builder.Select("id", "name", "email")
func (b *SQLBuilder) Select(columns ...string) *SQLBuilder {
	b.queryType = "SELECT"
	b.selectCols = append(b.selectCols, columns...)
	return b
}

// From specifies the table name for the query.
//
// Example:
//
//	builder.From("users")
func (b *SQLBuilder) From(table string) *SQLBuilder {
	b.tableName = table
	return b
}

// Insert specifies the table name for an INSERT query.
//
// Example:
//
//	builder.Insert("users")
func (b *SQLBuilder) Insert(table string) *SQLBuilder {
	b.queryType = "INSERT"
	b.tableName = table
	return b
}

// Columns specifies the columns for an INSERT query.
//
// Example:
//
//	builder.Columns("name", "email", "age")
func (b *SQLBuilder) Columns(columns ...string) *SQLBuilder {
	b.insertCols = append(b.insertCols, columns...)
	return b
}

// Values specifies the values for an INSERT query.
// The number of values should match the number of columns.
//
// Example:
//
//	builder.Values("John Doe", "john@example.com", 30)
func (b *SQLBuilder) Values(values ...any) *SQLBuilder {
	b.values = append(b.values, values...)
	return b
}

// Update specifies the table name for an UPDATE query.
//
// Example:
//
//	builder.Update("users")
func (b *SQLBuilder) Update(table string) *SQLBuilder {
	b.queryType = "UPDATE"
	b.tableName = table
	return b
}

// Set adds a SET clause for an UPDATE query.
// Multiple Set calls can be chained to set multiple columns.
//
// Example:
//
//	builder.Set("name = ?", "Jane Doe").Set("age = ?", 31)
func (b *SQLBuilder) Set(clause string, args ...any) *SQLBuilder {
	b.setClauses = append(b.setClauses, setClause{
		clause: clause,
		args:   args,
	})
	return b
}

func (b *SQLBuilder) Returning(columns ...string) *SQLBuilder {
	b.returning = append(b.returning, columns...)
	return b
}

// Delete starts a DELETE query.
//
// Example:
//
//	builder.Delete().From("users")
func (b *SQLBuilder) Delete() *SQLBuilder {
	b.queryType = "DELETE"
	return b
}

// Where adds a WHERE condition to the query.
// Multiple Where calls are combined with AND.
//
// Example:
//
//	builder.Where("age > ?", 18).Where("status = ?", "active")
func (b *SQLBuilder) Where(condition string, args ...any) *SQLBuilder {
	b.whereConds = append(b.whereConds, whereCondition{
		condition: condition,
		args:      args,
	})
	return b
}

// OrderBy specifies the ORDER BY clause for a SELECT query.
//
// Example:
//
//	builder.OrderBy("created_at DESC")
func (b *SQLBuilder) OrderBy(column string) *SQLBuilder {
	b.orderByCol = column
	return b
}

// Limit specifies the LIMIT clause for a SELECT query.
//
// Example:
//
//	builder.Limit(10)
func (b *SQLBuilder) Limit(limit int) *SQLBuilder {
	b.limitVal = limit
	return b
}

// Offset specifies the OFFSET clause for a SELECT query.
//
// Example:
//
//	builder.Offset(20)
func (b *SQLBuilder) Offset(offset int) *SQLBuilder {
	b.offsetVal = offset
	return b
}

// Build constructs and returns the final SQL query and its arguments.
// Returns the query string and a slice of arguments for parameterized queries.
//
// Example:
//
//	query, args := builder.Build()
//	// Use with database/sql: db.Query(query, args...)
func (b *SQLBuilder) Build() (string, []any) {
	switch b.queryType {
	case "SELECT":
		return b.buildSelect()
	case "INSERT":
		return b.buildInsert()
	case "UPDATE":
		return b.buildUpdate()
	case "DELETE":
		return b.buildDelete()
	default:
		return "", nil
	}
}

// replacePlaceholders replaces ? placeholders with PostgreSQL-style $1, $2, etc.
// The placeholderNum is passed by reference and incremented for each placeholder found.
func replacePlaceholders(clause string, placeholderNum *int) string {
	result := ""
	for i := 0; i < len(clause); i++ {
		if clause[i] == '?' {
			result += fmt.Sprintf("$%d", *placeholderNum)
			*placeholderNum++
		} else {
			result += string(clause[i])
		}
	}
	return result
}

// buildSelect constructs a SELECT query.
func (b *SQLBuilder) buildSelect() (string, []any) {
	var query strings.Builder
	args := make([]any, 0)

	// SELECT clause
	query.WriteString("SELECT ")
	if len(b.selectCols) == 0 {
		query.WriteString("*")
	} else {
		query.WriteString(strings.Join(b.selectCols, ", "))
	}

	// FROM clause
	if b.tableName != "" {
		query.WriteString(" FROM ")
		query.WriteString(b.tableName)
	}

	// WHERE clause
	if len(b.whereConds) > 0 {
		query.WriteString(" WHERE ")
		conditions := make([]string, len(b.whereConds))
		placeholderNum := 1
		for i, cond := range b.whereConds {
			conditions[i] = replacePlaceholders(cond.condition, &placeholderNum)
			args = append(args, cond.args...)
		}
		query.WriteString(strings.Join(conditions, " AND "))
	}

	// ORDER BY clause
	if b.orderByCol != "" {
		query.WriteString(" ORDER BY ")
		query.WriteString(b.orderByCol)
	}

	// LIMIT clause
	if b.limitVal >= 0 {
		query.WriteString(fmt.Sprintf(" LIMIT %d", b.limitVal))
	}

	// OFFSET clause
	if b.offsetVal >= 0 {
		query.WriteString(fmt.Sprintf(" OFFSET %d", b.offsetVal))
	}

	return query.String(), args
}

// buildInsert constructs an INSERT query.
func (b *SQLBuilder) buildInsert() (string, []any) {
	var query strings.Builder

	query.WriteString("INSERT INTO ")
	query.WriteString(b.tableName)

	// Columns
	if len(b.insertCols) > 0 {
		query.WriteString(" (")
		query.WriteString(strings.Join(b.insertCols, ", "))
		query.WriteString(")")
	}

	// Values
	query.WriteString(" VALUES (")
	placeholders := make([]string, len(b.values))
	for i := range b.values {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}
	query.WriteString(strings.Join(placeholders, ", "))
	query.WriteString(")")

	// RETURNING clause
	if len(b.returning) > 0 {
		query.WriteString(" RETURNING ")
		query.WriteString(strings.Join(b.returning, ", "))
	}

	return query.String(), b.values
}

// buildUpdate constructs an UPDATE query.
func (b *SQLBuilder) buildUpdate() (string, []any) {
	var query strings.Builder
	args := make([]any, 0)

	query.WriteString("UPDATE ")
	query.WriteString(b.tableName)

	// SET clause
	placeholderNum := 1
	if len(b.setClauses) > 0 {
		query.WriteString(" SET ")
		clauses := make([]string, len(b.setClauses))
		for i, set := range b.setClauses {
			clauses[i] = replacePlaceholders(set.clause, &placeholderNum)
			args = append(args, set.args...)
		}
		query.WriteString(strings.Join(clauses, ", "))
	}

	// WHERE clause
	if len(b.whereConds) > 0 {
		query.WriteString(" WHERE ")
		conditions := make([]string, len(b.whereConds))
		for i, cond := range b.whereConds {
			conditions[i] = replacePlaceholders(cond.condition, &placeholderNum)
			args = append(args, cond.args...)
		}
		query.WriteString(strings.Join(conditions, " AND "))
	}

	// RETURNING clause
	if len(b.returning) > 0 {
		query.WriteString(" RETURNING ")
		query.WriteString(strings.Join(b.returning, ", "))
	}

	return query.String(), args
}

// buildDelete constructs a DELETE query.
func (b *SQLBuilder) buildDelete() (string, []any) {
	var query strings.Builder
	args := make([]any, 0)

	query.WriteString("DELETE FROM ")
	query.WriteString(b.tableName)

	// WHERE clause
	placeholderNum := 1
	if len(b.whereConds) > 0 {
		query.WriteString(" WHERE ")
		conditions := make([]string, len(b.whereConds))
		for i, cond := range b.whereConds {
			conditions[i] = replacePlaceholders(cond.condition, &placeholderNum)
			args = append(args, cond.args...)
		}
		query.WriteString(strings.Join(conditions, " AND "))
	}

	// RETURNING clause
	if len(b.returning) > 0 {
		query.WriteString(" RETURNING ")
		query.WriteString(strings.Join(b.returning, ", "))
	}

	return query.String(), args
}