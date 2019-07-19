package cache

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// Cache - class of cache
type Cache struct {
	maxElementsFirstLavel  int
	maxElementsSecondLavel int
	mapElements            *map[string]*interface{}
	connectionString       string
	connection             *sql.DB
}

// Init - class instance initialization
func (c *Cache) Init(MaxElementsFirstLavel int, MaxElementsSecondLavel int, ConnectionString string, CreateTable bool) (err error) {
	c.maxElementsFirstLavel = MaxElementsFirstLavel
	c.maxElementsSecondLavel = MaxElementsSecondLavel
	c.connectionString = ConnectionString
	c.connection, err = sql.Open("sqlite3", c.connectionString)
	if err == nil {
		defer c.connection.Close()
		if CreateTable {
			_, err = c.connection.Exec("CREATE TABLE `Hache` ( `ID` INTEGER PRIMARY KEY AUTOINCREMENT, `Hache` TEXT NOT NULL UNIQUE, `Value` BLOB, `Old` INTEGER )")
			_, err = c.connection.Exec("CREATE INDEX `HacheByOld` ON `Hache` ( `Old` )")
			if err == nil {
				fmt.Printf("Таблица создана\n")
			}
		}
	}

	return err
}
