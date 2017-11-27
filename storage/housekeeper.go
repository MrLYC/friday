package storage

// CleanNormalItems :
func CleanNormalItems(conn *DatabaseConnection) error {
	return conn.WhereExpired().Delete(Item{}).Error
}

// CleanNormalItemTags :
func CleanNormalItemTags(conn *DatabaseConnection) error {
	return conn.WhereExpired().Delete(ItemTag{}).Error
}
