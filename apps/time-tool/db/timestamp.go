package db

import (
	"encoding/json"
	"finlab/apps/time-tool/core"
)

const (
	schemaTimestampSql = `CREATE TABLE IF NOT EXISTS timestamps (
        id INTEGER NOT NULL PRIMARY KEY,
        timestamp DATETIME NOT NULL,
        type VARCHAR(32));`
	insertTimestampSql    = `INSERT INTO timestamps (timestamp, type) VALUES (?, ?);`
	selectTimestampSql    = `SELECT * FROM timestamps`
	deleteTimestampSql    = `DELETE FROM timestamps`
	whereIdInTimestampSql = `WHERE id IN (%s)`
)

func GetTimestampsRepo() (repo *SQLiteRepository, err error) {
	repo, err = Repo(schemaTimestampSql)
	if err != nil {
		return nil, err
	}
	return
}

func SetTimestamp(timestamp core.Timestamp) error {
	repo, err := GetTimestampsRepo()
	if err != nil {
		return err
	}

	_, err = repo.Create(insertTimestampSql, timestamp.Timestamp, timestamp.Type)
	return err
}

func GetTimestamps(ids ...[]int) ([]core.Timestamp, error) {
	repo, err := GetTimestampsRepo()
	if err != nil {
		return nil, err
	}

	var timestampsJson []string

	if len(ids) > 0 && ids[0] != nil {
		timestampsJson, err = repo.Select(selectTimestampSql, ids[0])
	} else {
		timestampsJson, err = repo.Select(selectTimestampSql)
	}

	if err != nil {
		return nil, err
	}

	var timestamps []core.Timestamp

	for _, timestampJson := range timestampsJson {
		var timestamp core.Timestamp
		json.Unmarshal([]byte(timestampJson), &timestamp)
		timestamps = append(timestamps, timestamp)
	}

	return timestamps, nil
}

func DeleteTimestampsByIds(ids []int) ([]core.Timestamp, error) {
	repo, err := GetTimestampsRepo()
	if err != nil {
		return nil, err
	}

	timestamps, err := GetTimestamps(ids)
	if err != nil {
		return nil, err
	}

	err = repo.Delete(deleteTimestampSql, ids)
	if err != nil {
		return nil, err
	}

	return timestamps, nil
}

func DeleteTimestamps() error {
	repo, err := GetTimestampsRepo()
	if err != nil {
		return err
	}

	err = repo.Delete(deleteTimestampSql)
	if err != nil {
		return err
	}

	return nil
}
