package migration

import "friday/storage"

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
		conn.Create(&storage.ItemTag{
			Name: name,
		})
	}
	return nil
}

// Rollback171118005546 :
func (c *Command) Rollback171118005546(migration *Migration, conn *storage.DatabaseConnection) error {
	for _, name := range []string{
		ItemTagTypeString, ItemTagTypeList, ItemTagTypeTable, ItemTagTypeSet,
	} {
		conn.Delete(storage.ItemTag{}, "name = ?", name)
	}
	return nil
}
