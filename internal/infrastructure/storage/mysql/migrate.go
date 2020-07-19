package mysql

import (
	"log"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	// _ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"

	_ "github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
)

type Migration struct {
	m *migrate.Migrate
}

func NewMigration(mySQLURL string) *Migration {
	m, err := migrate.New("file://cmd/migration", mySQLURL)
	if err != nil {
		log.Fatal(err)
	}

	return &Migration{m: m}
}

func (mi *Migration) Up() error {
	return mi.m.Up()
}

func (mi *Migration) Down() error {
	return mi.m.Down()
}

func (mi *Migration) Version() (version uint, dirty bool, err error) {
	return mi.m.Version()
}

func (mi *Migration) Drop() error {
	return mi.m.Drop()
}

func (mi *Migration) Force(version int) error {
	return mi.m.Force(version)
}

func (mi *Migration) Steps(steps int) error {
	return mi.m.Steps(steps)
}
