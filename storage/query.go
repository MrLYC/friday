package storage

import (
	"time"

	"github.com/jinzhu/gorm"

	"friday/utils"
)

// QueryExpireConnection :
type QueryExpireConnection struct {
	DatabaseConnection
}

// WhereNotExpires :
func (c *DatabaseConnection) WhereNotExpires() *DatabaseConnection {
	return c.CopyWithDB(c.Where(
		"expire_at IS NULL OR expire_at >= ?", time.Now(),
	).Where("status = ?", ModelStatusNormal))
}

// WhereExpired :
func (c *DatabaseConnection) WhereExpired() *DatabaseConnection {
	return c.CopyWithDB(c.Where(
		"expire_at < ?", time.Now(),
	).Where("status = ?", ModelStatusNormal))
}

// NewExpireQueryConn :
func NewExpireQueryConn(db IBaseDB) *QueryExpireConnection {
	newDB := &QueryExpireConnection{}
	switch t := db.(type) {
	case *gorm.DB:
		newDB.DB = t
	case IDB:
		newDB.DB = t.GetDB()
	default:
		panic(utils.Errorf("unknown db type"))
	}
	return newDB
}
