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

// MigrationStatus
const (
	MigrationStatusNone       = ""
	MigrationStatusMigrating  = "migrating"
	MigrationStatusMigrated   = "migrated"
	MigrationStatusError      = "error"
	MigrationStatusRollbacked = "rollbacked"
)

// Migration :
type Migration struct {
	ID         uint `gorm:"primary_key"`
	CreatedAt  time.Time
	MigrateAt  time.Time
	RollbackAt *time.Time
	Status     string `gorm:"size:32"`
	Name       string `gorm:"size:32"`

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
	if m.Name == "" {
		return false
	}
	db := storage.GetDBConnection()
	err := db.Last(m, "name = ?", m.Name).Error
	if err == storage.ErrRecordNotFound {
		return false
	} else if err != nil {
		panic(err)
	}
	return true
}

// SaveToDB :
func (m *Migration) SaveToDB() bool {
	var (
		db      = storage.GetDBConnection()
		err     error
		created = false
	)
	if db.NewRecord(m) {
		err = db.Create(m).Error
		created = true
	} else {
		err = db.Save(m).Error
	}
	if err != nil {
		panic(err)
	}
	return created
}

// TableName :
func (m *Migration) TableName() string {
	return "migration"
}

// ToString :
func (m *Migration) ToString() string {
	flag := ""
	remark := ""
	switch m.Status {
	case MigrationStatusMigrated:
		flag = "+"
		remark = fmt.Sprintf("(%s)", m.MigrateAt.Format("2006-01-02 15:04:05"))
	case MigrationStatusRollbacked:
		flag = "-"
		remark = fmt.Sprintf("(%s)", m.RollbackAt.Format("2006-01-02 15:04:05"))
	case MigrationStatusMigrating:
		flag = "*"
	case MigrationStatusError:
		flag = "x"
	case MigrationStatusNone:
		flag = " "
	default:
		flag = "?"
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

// ActionList :
func (c *Command) ActionList() error {
	migrations := c.GetMigrationsByMethod()
	for _, migration := range migrations {
		migration.FetchFromDB()
		fmt.Printf("%s\n", migration.ToString())
	}
	return nil
}

// ActionRun :
func (c *Command) ActionRun() error {
	var (
		migrations = c.GetMigrationsByMethod()
		db         = storage.GetDBConnection()
		err        error
		fun        MigrateFunc
		fetched    bool
	)
	for _, migration := range migrations {
		if err != nil {
			break
		}

		migration.Status = MigrationStatusMigrating
		fetched = migration.FetchFromDB()
		fmt.Printf("%s\n", migration.ToString())

		if fetched && migration.RollbackAt == nil {
			continue
		}

		db.LogMode(true)
		fun = migration.GetMigrateFunc()
		if fun != nil {
			err = fun(migration, db)
			migration.Status = MigrationStatusMigrated
		}
		if fun == nil || err != nil {
			migration.Status = MigrationStatusError
		}
		db.LogMode(false)
		migration.MigrateAt = time.Now()
		migration.RollbackAt = nil
		migration.SaveToDB()
	}
	return err
}

// ActionRollback :
func (c *Command) ActionRollback() error {
	var (
		migrations = c.GetMigrationsByMethod()
		db         = storage.GetDBConnection()
		now        = time.Now()
		fun        MigrateFunc
		err        error
	)
	for index := len(migrations) - 1; index >= 0; index-- {
		migration := migrations[index]
		migration.FetchFromDB()
		if migration.Status == MigrationStatusMigrated || migration.Status == MigrationStatusError {
			db.LogMode(true)
			fun = migration.GetRollbackFunc()
			if fun != nil {
				err = fun(migration, db)
				migration.Status = MigrationStatusRollbacked
			}
			if fun == nil || err != nil {
				migration.Status = MigrationStatusError
			}
			db.LogMode(false)
			migration.RollbackAt = &now
			migration.SaveToDB()
			break
		}
	}
	return err
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
		return c.ActionList()
	case "run":
		return c.ActionRun()
	case "rollback":
		return c.ActionRollback()
	}
	return nil
}
