package builder

import (
	"fmt"
	"strings"
)

// SQLBuilder provides a chainable API for building SQL queries.
// It supports SELECT, INSERT, UPDATE, and DELETE operations with
// a fluent interface for constructing queries dynamically.
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
//
//	// INSERT query
//	query, args := NewSQLBuilder().
//		Insert("products").
//		Columns("name", "price", "category").
//		Values("Laptop", 999.99, "electronics").
//		Build()
//
//	// UPDATE query
//	query, args := NewSQLBuilder().
//		Update("products").
//		Set("price = ?", 899.99).
//		Set("updated_at = ?", time.Now()).
//		Where("id = ?", 123).
//		Build()
//
//	// DELETE query
//	query, args := NewSQLBuilder().
//		Delete().
//		From("products").
//		Where("id = ?", 123).
//		Build()
type SQLBuilder struct {
	queryType   string   // SELECT, INSERT, UPDATE, DELETE
	selectCols  []string // Columns for SELECT
	tableName   string   // Table name
	insertCols  []string // Columns for INSERT
	values      []interface{}
	setClauses  []setClause
	whereConds  []whereCondition
	orderByCol  string
	limitVal    int
	offsetVal   int
}

type setClause struct {
	clause string
	args   []interface{}
}

type whereCondition struct {
	condition string
	args      []interface{}
}

// NewSQLBuilder creates a new SQLBuilder instance.
func NewSQLBuilder() *SQLBuilder {
	return &SQLBuilder{
		selectCols: make([]string, 0),
		insertCols: make([]string, 0),
		values:     make([]interface{}, 0),
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
func (b *SQLBuilder) Values(values ...interface{}) *SQLBuilder {
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
func (b *SQLBuilder) Set(clause string, args ...interface{}) *SQLBuilder {
	b.setClauses = append(b.setClauses, setClause{
		clause: clause,
		args:   args,
	})
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
func (b *SQLBuilder) Where(condition string, args ...interface{}) *SQLBuilder {
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
func (b *SQLBuilder) Build() (string, []interface{}) {
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

// buildSelect constructs a SELECT query.
func (b *SQLBuilder) buildSelect() (string, []interface{}) {
	var query strings.Builder
	args := make([]interface{}, 0)

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
		for i, cond := range b.whereConds {
			conditions[i] = cond.condition
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
func (b *SQLBuilder) buildInsert() (string, []interface{}) {
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
		placeholders[i] = "?"
	}
	query.WriteString(strings.Join(placeholders, ", "))
	query.WriteString(")")

	return query.String(), b.values
}

// buildUpdate constructs an UPDATE query.
func (b *SQLBuilder) buildUpdate() (string, []interface{}) {
	var query strings.Builder
	args := make([]interface{}, 0)

	query.WriteString("UPDATE ")
	query.WriteString(b.tableName)

	// SET clause
	if len(b.setClauses) > 0 {
		query.WriteString(" SET ")
		clauses := make([]string, len(b.setClauses))
		for i, set := range b.setClauses {
			clauses[i] = set.clause
			args = append(args, set.args...)
		}
		query.WriteString(strings.Join(clauses, ", "))
	}

	// WHERE clause
	if len(b.whereConds) > 0 {
		query.WriteString(" WHERE ")
		conditions := make([]string, len(b.whereConds))
		for i, cond := range b.whereConds {
			conditions[i] = cond.condition
			args = append(args, cond.args...)
		}
		query.WriteString(strings.Join(conditions, " AND "))
	}

	return query.String(), args
}

// buildDelete constructs a DELETE query.
func (b *SQLBuilder) buildDelete() (string, []interface{}) {
	var query strings.Builder
	args := make([]interface{}, 0)

	query.WriteString("DELETE FROM ")
	query.WriteString(b.tableName)

	// WHERE clause
	if len(b.whereConds) > 0 {
		query.WriteString(" WHERE ")
		conditions := make([]string, len(b.whereConds))
		for i, cond := range b.whereConds {
			conditions[i] = cond.condition
			args = append(args, cond.args...)
		}
		query.WriteString(strings.Join(conditions, " AND "))
	}

	return query.String(), args
}
