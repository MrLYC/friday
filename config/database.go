package config

import (
	"fmt"
	"strings"
)

// Database : Database meta configuration
type Database struct {
	Type string `yaml:"type" validate:"regexp=^(sqlite3|mysql|mssql|postgresql)$"`
	Name string `yaml:"name"`

	Connection *string `yaml:"connection,omitempty"`

	Host     *string `yaml:"host,omitempty"`
	Port     *uint   `yaml:"port,omitempty"`
	User     *string `yaml:"user,omitempty"`
	Password *string `yaml:"password,omitempty"`
}

// GetConnectionString :
func (d *Database) GetConnectionString() string {
	if d.Connection != nil {
		return *(d.Connection)
	}

	switch d.Type {
	case "sqlite3":
		return d.Name
	case "mysql":
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			*(d.User), *(d.Password), *(d.Host), *(d.Port), d.Name,
		)
	case "postgresql":
		return fmt.Sprintf(
			"host=%s port=%d user=%s dbname=%s sslmode=disable password=%s",
			*(d.Host), *(d.Port), *(d.User), d.Name, *(d.Password),
		)
	case "mssql":
		return fmt.Sprintf(
			"sqlserver://%s:%s@%s:%d?database=%s",
			*(d.User), *(d.Password), *(d.Host), *(d.Port), d.Name,
		)
	}
	panic(fmt.Errorf("Can not make connection string for %v", d))
}

// Init : init Database
func (d *Database) Init() {
	for _, t := range strings.Split(BuildTag, ",") {
		switch t {
		case "dball", "dbsqlite":
			d.initSQLite()
			return
		case "dbmysql":
			d.initMySQL()
			return
		case "dbpostgresql":
			d.initPostgreSQL()
			return
		case "dbmssql":
			d.initSQLServer()
			return
		default:
			d.initSQLite()
		}
	}
}

func (d *Database) initSQLite() {
	d.Type = "sqlite3"
	d.Name = "friday.db"
}

func (d *Database) initMySQL() {
	d.Type = "mysql"
	d.Name = "friday"

	d.Host = new(string)
	*(d.Host) = "localhost"

	d.Port = new(uint)
	*(d.Port) = 3306

	d.User = new(string)
	*(d.User) = "root"

	d.Password = new(string)
	*(d.Password) = "mrlyc"
}

func (d *Database) initPostgreSQL() {
	d.Type = "postgresql"
	d.Name = "friday"

	d.Host = new(string)
	*(d.Host) = "localhost"

	d.Port = new(uint)
	*(d.Port) = 5432

	d.User = new(string)
	*(d.User) = "root"

	d.Password = new(string)
	*(d.Password) = "mrlyc"
}

func (d *Database) initSQLServer() {
	d.Type = "mssql"
	d.Name = "friday"

	d.Host = new(string)
	*(d.Host) = "localhost"

	d.Port = new(uint)
	*(d.Port) = 1433

	d.User = new(string)
	*(d.User) = "root"

	d.Password = new(string)
	*(d.Password) = "mrlyc"
}
