package api

import (
	"finlab/apps/time-tool/config"
	"finlab/apps/time-tool/core"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func db() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(config.DbPath()), &gorm.Config{})
	if err != nil {
		panic("Failed to open the SQLite database!")
	}

	return db
}

func SetTimestamp(timestamp core.Timestamp) error {
	err := db().AutoMigrate(&core.Timestamp{})
	if err != nil {
		return err
	}
	db().Create(&timestamp)
	return nil
}

func GetTimestamps() ([]core.Timestamp, error) {
	var timestamps []core.Timestamp
	err := db().AutoMigrate(&core.Timestamp{})
	if err != nil {
		return timestamps, err
	}
	db().Find(&timestamps)

	return timestamps, nil
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
