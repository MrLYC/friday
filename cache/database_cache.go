package cache

import (
	"time"

	"friday/logging"
	"friday/storage"
	"friday/utils"
)

// DBCache :
type DBCache struct {
	Conn *storage.DatabaseConnection

	TagTypeString *storage.ItemTag
	TagTypeList   *storage.ItemTag
	TagTypeTable  *storage.ItemTag
}

// Init :
func (c *DBCache) Init() {
	var (
		err error
	)
	c.Conn = storage.GetDBConnection()

	c.TagTypeString = &storage.ItemTag{}
	err = c.Conn.Find(c.TagTypeString, "name = ?", storage.ItemTagTypeString).Error
	if err != nil {
		panic(utils.ErrorWrap(err))
	}

	c.TagTypeList = &storage.ItemTag{}
	err = c.Conn.Find(c.TagTypeList, "name = ?", storage.ItemTagTypeList).Error
	if err != nil {
		panic(utils.ErrorWrap(err))
	}

	c.TagTypeTable = &storage.ItemTag{}
	err = c.Conn.Find(c.TagTypeTable, "name = ?", storage.ItemTagTypeTable).Error
	if err != nil {
		panic(utils.ErrorWrap(err))
	}
}

// Close :
func (c *DBCache) Close() error {
	return c.Conn.Close()
}

// GetItemByKeyAndTag :
func (c *DBCache) GetItemByKeyAndTag(key string, tag string) (*storage.Item, error) {
	item := &storage.Item{}

	result := c.Conn.WhereNotExpires().Preload(
		"Tags", "name = ?",
		c.TagTypeString.Name,
	).Where(
		"key = ?", key,
	).Last(item)

	if result.RecordNotFound() {
		return nil, nil
	} else if result.Error != nil {
		return nil, result.Error
	}
	return item, nil
}

// MakeItemExpired :
func (c *DBCache) MakeItemExpired(key string, expireTime time.Time, args ...interface{}) error {
	db := c.Conn.Model(&storage.Item{}).Where(
		"key = ?", key,
	)
	if len(args) > 1 {
		db = db.Where(args[0], args[1:])
	}
	return db.Update(
		"ExpireAt", expireTime,
	).Error
}

// MakeItemTagExpired :
func (c *DBCache) MakeItemTagExpired(key string, expireTime time.Time, args ...interface{}) error {
	db := c.Conn.Model(&storage.ItemTag{}).Where(
		"name = ?", key,
	)
	if len(args) > 1 {
		db = db.Where(args[0], args[1:])
	}
	return db.Update(
		"ExpireAt", expireTime,
	).Error
}

// GetKey :
func (c *DBCache) GetKey(key string) (string, error) {
	item, err := c.GetItemByKeyAndTag(key, c.TagTypeString.Name)
	if item == nil {
		return "", err
	}
	return item.Value, nil
}

// SetKey :
func (c *DBCache) SetKey(key string, val string) error {
	item := &storage.Item{
		Key: key, Value: val,
		Tags: []*storage.ItemTag{
			c.TagTypeString,
		},
	}
	err := c.Conn.Create(item).Error
	if err != nil {
		return err
	}

	err = c.Conn.WhereNotExpires().Model(&storage.Item{}).Where(
		"key = ? AND (id < ? OR created_at < ?)",
		key, item.ID, item.CreatedAt,
	).Update("expire_at", item.CreatedAt).Error
	if err != nil {
		logging.Errorf("MakeKeyExpired error: %v", err)
	}
	err = c.MakeItemTagExpired(key, item.CreatedAt)
	if err != nil {
		logging.Errorf("MakeKeyExpired error: %v", err)
	}

	return nil
}
