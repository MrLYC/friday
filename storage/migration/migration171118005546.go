package migration

import "friday/storage"
import "time"

//
const (
	ItemTagTypeString = "STRING"
	ItemTagTypeList   = "LIST"
	ItemTagTypeTable  = "TABLE"
	ItemTagTypeSet    = "SET"
)

// Migrate171118005546 :
func (c *Command) Migrate171118005546(migration *Migration, conn *storage.DatabaseConnection) error {
	for _, name := range []string{
		ItemTagTypeString, ItemTagTypeList, ItemTagTypeTable, ItemTagTypeSet,
	} {
		conn.Create(&ItemTag171112224339{
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
		ItemTagTypeString, ItemTagTypeList, ItemTagTypeTable, ItemTagTypeSet,
	} {
		conn.Delete(ItemTag171112224339{}, "name = ?", name)
	}
	return nil
}
