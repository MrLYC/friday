package config

// Database : Database meta configuration
type Database struct {
	Type string `yaml:"type" validate:"regexp=^(sqlite3|mysql|mssql|postgres)$"`
	Name string `yaml:"name"`
}

// Init : init Database
func (d *Database) Init() {
	d.Type = "sqlite3"
	d.Name = "friday.db"
}
