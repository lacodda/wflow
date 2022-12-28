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

func SetTimestamp(timestamp core.Timestamp) {
	err := db().AutoMigrate(&core.Timestamp{})
	if err != nil {
		core.Danger("Migration failure!")
		return
	}
	db().Create(&timestamp)
}

func GetTimestamps() []core.Timestamp {
	var timestamps []core.Timestamp
	db().Find(&timestamps)

	return timestamps
}
