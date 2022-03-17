package db

import (
	"database/sql"
	"fmt"

	"github.com/parthoshuvo/authsvc/cfg"
	log "github.com/parthoshuvo/authsvc/log4u"

	_ "github.com/go-sql-driver/mysql" // mysql driver
)

// AuthDB database.
type AuthDB struct {
	db *sql.DB
}

// NewAuthDB creates a DB handler.
func NewAuthDB(dbDef *cfg.DBDef) *AuthDB {
	return &AuthDB{openDatabase(dbDef)}
}

func openDatabase(dbDef *cfg.DBDef) *sql.DB {
	dbSrc := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbDef.User, dbDef.Password, dbDef.Host, dbDef.Port, dbDef.Database)
	db, err := sql.Open("mysql", dbSrc)
	if err != nil {
		log.Fatalf("failed to open database %s:%d/%s: [%v]", dbDef.Host, dbDef.Port, dbDef.Database, err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("failed to open database %s:%d/%s: [%v]", dbDef.Host, dbDef.Port, dbDef.Database, err)
	}
	log.Info("Pong!! successfully database is connected!!")
	return db
}

// Close closes database connection.
func (ad *AuthDB) Close() {
	ad.db.Close()
}
