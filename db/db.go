package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"wflow/config"
	"wflow/core"
	"fmt"
	"strconv"
	"strings"

	"github.com/mattn/go-sqlite3"
)

const (
	whereIdInSql = `WHERE id IN (%s)`
)

var (
	ErrDuplicate    = errors.New("Record already exists")
	ErrDeleteFailed = errors.New("Delete failed")
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		db: db,
	}
}

func Repo(schemaSql string) (*SQLiteRepository, error) {
	db, err := sql.Open("sqlite3", config.DbPath())
	if err != nil {
		return nil, err
	}
	sqlRepository := NewSQLiteRepository(db)

	if err = sqlRepository.Migrate(schemaSql); err != nil {
		return nil, err
	}

	return sqlRepository, nil
}

func (r *SQLiteRepository) Migrate(query string) error {
	_, err := r.db.Exec(query)
	return err
}

func (r *SQLiteRepository) Create(query string, args ...any) (id int64, err error) {
	res, err := r.db.Exec(query, args...)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return id, ErrDuplicate
			}
		}
		return id, err
	}

	id, err = res.LastInsertId()
	if err != nil {
		return id, err
	}

	return id, nil
}

func (r *SQLiteRepository) Select(query string, ids ...[]int) ([]string, error) {
	if len(ids) > 0 && ids[0] != nil {
		query = fmt.Sprintf("%s %s", query, fmt.Sprintf(whereIdInSql, core.ArrayToString(ids[0], ",")))
	}

	rows, err := r.db.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return RowsToJson(rows), nil
}

func (r *SQLiteRepository) Delete(query string, ids ...[]int) error {
	if len(ids) > 0 && ids[0] != nil {
		query = fmt.Sprintf("%s %s", query, fmt.Sprintf(whereIdInSql, core.ArrayToString(ids[0], ",")))
	}

	res, err := r.db.Exec(query)

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrDeleteFailed
	}

	return err
}

func RowsToJson(rows *sql.Rows) []string {
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}

	values := make([]interface{}, len(columns))

	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	c := 0
	results := make(map[string]interface{})
	data := []string{}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}

		for i, value := range values {
			switch value.(type) {
			case nil:
				results[columns[i]] = nil

			case []byte:
				s := string(value.([]byte))
				x, err := strconv.Atoi(s)

				if err != nil {
					results[columns[i]] = s
				} else {
					results[columns[i]] = x
				}

			default:
				results[columns[i]] = value
			}
		}

		b, _ := json.Marshal(results)
		data = append(data, strings.TrimSpace(string(b)))
		c++
	}

	return data
}
