package database

import (
	"gorm.io/gorm"
)

var AutoMaintainRange = []any{}

func RunMigration(source *gorm.DB) error {
	if err := source.AutoMigrate(
		AutoMaintainRange...,
	); err != nil {
		return err
	}

	return nil
}
