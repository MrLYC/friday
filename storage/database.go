package storage

import (
	"friday/config"
	"friday/logging"
	"sync"
	"time"

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

// DatabaseConnection :
type DatabaseConnection struct {
	*gorm.DB
}

// CopyWithDB :
func (c *DatabaseConnection) CopyWithDB(db *gorm.DB) *DatabaseConnection {
	return &DatabaseConnection{
		DB: db,
	}
}

// WhereNotExpires :
func (c *DatabaseConnection) WhereNotExpires() *DatabaseConnection {
	return c.CopyWithDB(c.Where(
		"expire_at IS NULL OR expire_at <= ?", time.Now(),
	))
}

// WhereExpired :
func (c *DatabaseConnection) WhereExpired() *DatabaseConnection {
	return c.CopyWithDB(c.Where(
		"expire_at > ?", time.Now(),
	))
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
