package psql

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
)

type repository[M any] struct {
	db    *sql.DB
	table string
}

func NewRepository[M any](db *sql.DB, table string) Repository[M] {
	return &repository[M]{db: db}
}

func (r *repository[M]) Select(ctx context.Context, filters map[string][]any) ([]M, error) {
	query, args, err := queryBuilder(r.table, filters)
	if err != nil {
		return nil, err
	}

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	results := make([]M, 0)

	for rows.Next() {
		m := new(M)
		if err := rows.Scan(&m); err != nil {
			return nil, err
		}

		results = append(results, *m)
	}

	return results, nil
}

const (
	DISJUNCTION string = "disjunction"
	CONJUNCTION string = "conjuction"
)

func NewConnective(c string) (string, error) {
	switch strings.ToLower(c) {
	case "disjunction":
		return DISJUNCTION, nil
	case "conjuction":
		return CONJUNCTION, nil
	default:
		return "", ErrNoConnective
	}
}

var ErrNoConnective = errors.New("no connective")

func queryBuilder(table string, input map[string][]any, connective string) (string, []any, error) {
	builder := strings.Builder{}

	builder.WriteString("SELECT * FROM ")
	builder.WriteString(table)
	builder.WriteString(" WHERE ")

	args := make([]any, len(input))
	count := 1

	for k, v := range input {
		if len(v) == 0 {
			continue
		}

		if count > 1 {
			connective, err := NewConnective(connective)
			if err != nil {
				break
			}

			if connective == DISJUNCTION {
				builder.WriteString(" OR ")
			} else {
				builder.WriteString(" AND ")
			}
		}

		switch {
		case len(v) == 1:
			builder.WriteString(k)
			builder.WriteString(" = $")
			builder.WriteString(strconv.Itoa(count))
			args[count-1] = v[0]
		case len(v) > 1:
			builder.WriteString(k)
			builder.WriteString(" IN (")
			argsForIn := make([]any, len(v))
			for i, arg := range v {
				argsForIn[i] = arg
				args = append(args, arg)
			}
			builder.WriteString(strings.Repeat("$"+strconv.Itoa(count+i)+",", len(v))[:len(strings.Repeat("$"+strconv.Itoa(count+i)+",", len(v))-1])
			count += len(v)
			builder.WriteString(")")
		}

		count++
	}

	builder.WriteString(";")

	return builder.String(), nil, nil
}

