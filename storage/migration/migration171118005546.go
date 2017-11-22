package migration

import "friday/storage"
import "time"

//
const (
	ItemTagTypeString = "STRING"
	ItemTagTypeList   = "LIST"
	ItemTagTypeTable  = "TABLE"
)

// Migrate171118005546 :
func (c *Command) Migrate171118005546(migration *Migration, conn *storage.DatabaseConnection) error {
	for _, name := range []string{
		ItemTagTypeString, ItemTagTypeList, ItemTagTypeTable,
	} {
		conn.Create(&ItemTag{
			Name:      name,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}
	return nil
}

// Rollback171118005546 :
func (c *Command) Rollback171118005546(migration *Migration, conn *storage.DatabaseConnection) error {
	for _, name := range []string{
		ItemTagTypeString, ItemTagTypeList, ItemTagTypeTable,
	} {
		conn.Delete(ItemTag{}, "name = ?", name)
	}
	return nil
}
