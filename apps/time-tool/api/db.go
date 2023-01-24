package api

import (
	"database/sql"
	"errors"
	"finlab/apps/time-tool/config"
	"finlab/apps/time-tool/core"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	insertSQL     = `INSERT INTO timestamps (timestamp, type) VALUES (?, ?);`
	selectAllSQL  = `SELECT * FROM timestamps;`
	deleteAllSQL  = `DELETE FROM timestamps;`
	deleteByIdSQL = `DELETE FROM timestamps WHERE id = ?;`
	schemaSQL     = `CREATE TABLE IF NOT EXISTS timestamps (
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

	if _, err = sqlDB.Exec(schemaSQL); err != nil {
		return nil, err
	}

	stmt, err := sqlDB.Prepare(insertSQL)
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

func (db *DB) Select() ([]core.Timestamp, error) {
	var timestamps []core.Timestamp
	rows, err := db.sql.Query(selectAllSQL)

	if err != nil {
		return nil, err
	}

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
	defer db.Close()
	timestamps, err := db.Select()
	if err != nil {
		return timestamps, err
	}
	return timestamps, nil
}

// ++++++++++++++++++++++++++++++++++++++++++

func db() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(config.DbPath()), &gorm.Config{})
	if err != nil {
		panic("Failed to open the SQLite database!")
	}

	return db
}

func DeleteTimestamps() {
	db().Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&core.Timestamp{})
}

func DeleteTimestampsByIds(ids []int) []core.Timestamp {
	var timestamps []core.Timestamp
	db().Where(ids).Find(&timestamps)
	db().Delete(&core.Timestamp{}, ids)
	return timestamps
}
