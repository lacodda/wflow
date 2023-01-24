package api

import (
	"database/sql"
	"errors"
	"finlab/apps/time-tool/config"
	"finlab/apps/time-tool/core"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const (
	insertSql    = `INSERT INTO timestamps (timestamp, type) VALUES (?, ?);`
	selectSql    = `SELECT * FROM timestamps`
	deleteSql    = `DELETE FROM timestamps`
	whereIdInSql = `WHERE id IN (%s)`
	schemaSql    = `CREATE TABLE IF NOT EXISTS timestamps (
        id INTEGER NOT NULL PRIMARY KEY,
        timestamp DATETIME NOT NULL,
        type VARCHAR(32));`
)

type DB struct {
	sql    *sql.DB
	stmt   *sql.Stmt
	buffer []core.Timestamp
}

func Db() (*DB, error) {
	sqlDB, err := sql.Open("sqlite3", config.DbPath())
	if err != nil {
		return nil, err
	}

	if _, err = sqlDB.Exec(schemaSql); err != nil {
		return nil, err
	}

	stmt, err := sqlDB.Prepare(insertSql)
	if err != nil {
		return nil, err
	}

	db := DB{
		sql:    sqlDB,
		stmt:   stmt,
		buffer: make([]core.Timestamp, 0, 1024),
	}
	return &db, nil
}

func (db *DB) Add(timestamp core.Timestamp) error {
	if len(db.buffer) == cap(db.buffer) {
		return errors.New("Timestamps buffer is full")
	}

	db.buffer = append(db.buffer, timestamp)
	if len(db.buffer) == cap(db.buffer) {
		if err := db.Flush(); err != nil {
			return err
		}
	}

	return nil
}

func (db *DB) Flush() error {
	tx, err := db.sql.Begin()
	if err != nil {
		return err
	}

	for _, timestamp := range db.buffer {
		_, err := tx.Stmt(db.stmt).Exec(timestamp.Timestamp, timestamp.Type)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	db.buffer = db.buffer[:0]
	return tx.Commit()
}

func (db *DB) Close() error {
	defer func() {
		db.stmt.Close()
		db.sql.Close()
	}()

	if err := db.Flush(); err != nil {
		return err
	}

	return nil
}

func (db *DB) Select(ids ...[]int) ([]core.Timestamp, error) {
	var timestamps []core.Timestamp
	var query = selectSql

	if len(ids) > 0 && ids[0] != nil {
		query = fmt.Sprintf("%s %s", selectSql, fmt.Sprintf(whereIdInSql, core.ArrayToString(ids[0], ",")))
	}

	rows, err := db.sql.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var timestamp core.Timestamp
		err = rows.Scan(&timestamp.Id, &timestamp.Timestamp, &timestamp.Type)
		if err != nil {
			return nil, err
		}
		timestamps = append(timestamps, timestamp)
	}

	return timestamps, nil
}

func (db *DB) Delete(ids ...[]int) error {
	var query = deleteSql

	if len(ids) > 0 && ids[0] != nil {
		query = fmt.Sprintf("%s %s", deleteSql, fmt.Sprintf(whereIdInSql, core.ArrayToString(ids[0], ",")))
	}

	_, err := db.sql.Exec(query)

	return err
}

func SetTimestamp(timestamp core.Timestamp) error {
	db, err := Db()
	if err != nil {
		return err
	}

	defer db.Close()

	err = db.Add(timestamp)
	return err
}

func GetTimestamps() ([]core.Timestamp, error) {
	db, err := Db()
	if err != nil {
		return nil, err
	}

	timestamps, err := db.Select()
	if err != nil {
		return timestamps, err
	}

	return timestamps, nil
}

func DeleteTimestampsByIds(ids []int) ([]core.Timestamp, error) {
	db, err := Db()
	if err != nil {
		return nil, err
	}

	timestamps, err := db.Select(ids)
	if err != nil {
		return nil, err
	}

	err = db.Delete(ids)
	if err != nil {
		return nil, err
	}

	return timestamps, nil
}

func DeleteTimestamps() error {
	db, err := Db()
	if err != nil {
		return err
	}

	err = db.Delete()
	if err != nil {
		return err
	}

	return nil
}
