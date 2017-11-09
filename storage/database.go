package storage

import (
	"friday/config"
	"sync"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Errors
var (
	ErrRecordNotFound       = gorm.ErrRecordNotFound
	ErrCantStartTransaction = gorm.ErrCantStartTransaction
	ErrInvalidSQL           = gorm.ErrInvalidSQL
	ErrInvalidTransaction   = gorm.ErrInvalidTransaction
	ErrUnaddressable        = gorm.ErrUnaddressable
)

// DatabaseConnection :
type DatabaseConnection struct {
	*gorm.DB
}

var dbConectOnce sync.Once
var dbConnection *DatabaseConnection

// ConnectDatabase :
func ConnectDatabase() (*DatabaseConnection, error) {
	conf := config.Configuration.Database
	db, err := gorm.Open(conf.Type, conf.Name)
	if err != nil {
		return nil, err
	}
	return &DatabaseConnection{
		DB: db,
	}, nil
}

// GetDBConnection :
func GetDBConnection() *DatabaseConnection {
	if dbConnection == nil {
		dbConectOnce.Do(func() {
			db, err := ConnectDatabase()
			if err != nil {
				panic(err)
			}
			dbConnection = db
		})
	}
	return dbConnection
}
