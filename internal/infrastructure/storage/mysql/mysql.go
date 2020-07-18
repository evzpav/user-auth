package mysql

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jinzhu/gorm"
)

// New creates new database connection to a mysql database
func New(url string) (*gorm.DB, error) {
	db, err := gorm.Open("mysql", url)
	if err != nil {
		return nil, err
	}

	db.LogMode(true)
	return db, nil
}
