package storage

import (
	"friday/config"
	"friday/logging"
	"sync"

	"github.com/jinzhu/gorm"
)

// Errors
var (
	ErrRecordNotFound       = gorm.ErrRecordNotFound
	ErrCantStartTransaction = gorm.ErrCantStartTransaction
	ErrInvalidSQL           = gorm.ErrInvalidSQL
	ErrInvalidTransaction   = gorm.ErrInvalidTransaction
	ErrUnaddressable        = gorm.ErrUnaddressable
)

// IBaseDB :
type IBaseDB interface {
	Close() error
}

// IDB :
type IDB interface {
	IBaseDB
	GetDB() *gorm.DB
	SetDB(*gorm.DB)
}

// DatabaseConnection :
type DatabaseConnection struct {
	*gorm.DB
}

// GetDB :
func (c *DatabaseConnection) GetDB() *gorm.DB {
	return c.DB
}

// SetDB :
func (c *DatabaseConnection) SetDB(db *gorm.DB) {
	c.DB = db
}

// CopyWithDB :
func (c *DatabaseConnection) CopyWithDB(db *gorm.DB) *DatabaseConnection {
	return &DatabaseConnection{
		DB: db,
	}
}

// Close :
func (c *DatabaseConnection) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}

var dbConectOnce sync.Once
var dbConnection *DatabaseConnection

// ConnectDatabase :
func ConnectDatabase() (*DatabaseConnection, error) {
	conf := config.Configuration.Database
	connectStr := conf.GetConnectionString()
	if config.Configuration.Debug {
		logging.Debugf("connect string: %s", connectStr)
	}
	db, err := gorm.Open(conf.Type, connectStr)
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
