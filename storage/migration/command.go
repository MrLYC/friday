package migration

import (
	"flag"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	"friday/storage"
)

// MigrateFunc :
type MigrateFunc func(*Migration, *storage.DatabaseConnection) error

// Migration :
type Migration struct {
	ID         uint `gorm:"primary_key"`
	CreatedAt  time.Time
	MigrateAt  *time.Time
	RollbackAt *time.Time
	Name       string `gorm:"size:32"`

	fetched      bool
	migrateFunc  func(*Migration, *storage.DatabaseConnection) error
	rollbackFunc func(*Migration, *storage.DatabaseConnection) error
}

// SetMigrateFunc :
func (m *Migration) SetMigrateFunc(fun MigrateFunc) {
	m.migrateFunc = fun
}

// SetRollbackFunc :
func (m *Migration) SetRollbackFunc(fun MigrateFunc) {
	m.rollbackFunc = fun
}

// GetMigrateFunc :
func (m *Migration) GetMigrateFunc() MigrateFunc {
	return m.migrateFunc
}

// GetRollbackFunc :
func (m *Migration) GetRollbackFunc() MigrateFunc {
	return m.rollbackFunc
}

// FetchFromDB :
func (m *Migration) FetchFromDB() bool {
	m.fetched = false
	if m.Name == "" {
		return false
	}
	db := storage.GetDBConnection()
	err := db.Last(m, "name = ?", m.Name).Error
	if err == storage.ErrRecordNotFound {
		m.fetched = true
		return false
	} else if err != nil {
		panic(err)
	}
	m.fetched = true
	return true
}

// TableName :
func (m *Migration) TableName() string {
	return "migration"
}

// ToString :
func (m *Migration) ToString() string {
	flag := "*"
	remark := ""
	if m.RollbackAt != nil {
		flag = "+"
		remark = fmt.Sprintf("(%s)", m.RollbackAt.Format("2006-01-02 15:04:05"))
	} else if m.MigrateAt != nil {
		flag = "-"
		remark = fmt.Sprintf("(%s)", m.MigrateAt.Format("2006-01-02 15:04:05"))
	} else if m.fetched {
		flag = " "
	}
	return fmt.Sprintf("[%s]%s%s", flag, m.Name, remark)
}

// MigrationSortedArr :
type MigrationSortedArr []*Migration

// Len :
func (m MigrationSortedArr) Len() int {
	return len(m)
}

// Swap :
func (m MigrationSortedArr) Swap(i int, j int) {
	m[i], m[j] = m[j], m[i]
}

// Less :
func (m MigrationSortedArr) Less(i int, j int) bool {
	return m[i].CreatedAt.Before(m[j].CreatedAt)
}

// Sort :
func (m MigrationSortedArr) Sort() {
	sort.Sort(m)
}

// Command : Migrate
type Command struct {
	Action       string
	TableOptions string
}

// GetDescription :
func (c *Command) GetDescription() string {
	return "Migrate database"
}

// SetFlags : set parsing flags
func (c *Command) SetFlags() {
	flag.StringVar(&c.Action, "action", "run", "Command action")
	flag.StringVar(&c.TableOptions, "table_options", "", "Table option")
}

// CreateMigrationTableIfNotExists :
func (c *Command) CreateMigrationTableIfNotExists() error {
	db := storage.GetDBConnection()
	var m Migration
	if !db.HasTable(m.TableName()) {
		return db.AutoMigrate(m).Error
	}
	return nil
}

func reflectMigrateFunc(method reflect.Value) MigrateFunc {
	return func(m *Migration, db *storage.DatabaseConnection) error {
		result := method.Call([]reflect.Value{
			reflect.ValueOf(m), reflect.ValueOf(db),
		})
		value := result[0].Interface()
		if value == nil {
			return nil
		}
		return value.(error)
	}
}

// GetMigrationsByMethod :
func (c *Command) GetMigrationsByMethod() MigrationSortedArr {
	var (
		err                  error
		ok                   bool
		rtype                = reflect.TypeOf(c)
		rvalue               = reflect.ValueOf(c)
		number               = rvalue.NumMethod()
		migrateMethodPrefix  = "Migrate"
		rollbackMethodPrefix = "Rollback"
		migrations           = make(MigrationSortedArr, 0)
		migration            *Migration
		migrateMethod        reflect.Value
		rollbackMethod       reflect.Value
		migrationTime        time.Time
	)

	for i := 0; i < number; i++ {
		method := rtype.Method(i)
		if !strings.HasPrefix(method.Name, migrateMethodPrefix) {
			continue
		}

		migration = &Migration{
			Name: strings.TrimPrefix(method.Name, migrateMethodPrefix),
		}
		migrationTime, err = time.Parse("060102150405", migration.Name)
		if err != nil {
			panic(err)
		}
		migration.CreatedAt = migrationTime

		migrateMethod = rvalue.Method(method.Index)
		migration.SetMigrateFunc(reflectMigrateFunc(migrateMethod))

		method, ok = rtype.MethodByName(fmt.Sprintf(
			"%s%s", rollbackMethodPrefix, migration.Name,
		))
		if ok {
			rollbackMethod = rvalue.Method(method.Index)
			migration.SetRollbackFunc(reflectMigrateFunc(rollbackMethod))
		}
		migrations = append(migrations, migration)
	}
	return migrations
}

// ListAction :
func (c *Command) ListAction() error {
	migrations := c.GetMigrationsByMethod()
	for _, migration := range migrations {
		migration.FetchFromDB()
		fmt.Printf("%s\n", migration.ToString())
	}
	return nil
}

// RunAction :
func (c *Command) RunAction() error {
	return nil
}

// Run : run command
func (c *Command) Run() error {
	db := storage.GetDBConnection()
	defer db.Close()

	err := c.CreateMigrationTableIfNotExists()
	if err != nil {
		return err
	}
	switch c.Action {
	case "list":
		return c.ListAction()
	case "run":
		return c.RunAction()
	}
	return nil
}
