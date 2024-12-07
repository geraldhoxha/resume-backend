package migration

import (
	"github.com/geraldhoxha/resume-backend/config"
	"github.com/geraldhoxha/resume-backend/graph/model"
)

func MigrateTable() {
	db := config.GetDB()
	err := db.AutoMigrate(&model.User{})
	if err != nil {
		panic(err)
	}
}
