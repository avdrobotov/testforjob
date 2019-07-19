package cache

import (
	"database/sql"
	"errors"
	"fmt"
	"math"

	// driver for second lavel of cache
	_ "github.com/mattn/go-sqlite3"
)

type itemAllElements struct {
	Old   int
	Lavel int
}

// Cache - class of cache
type Cache struct {
	maxElementsFirstLavel  int
	maxElementsSecondLavel int
	mapElements            map[string]*interface{}
	mapAllElements         map[string]*itemAllElements
	connectionString       string
	connection             *sql.DB
}

const minOld = 10000
const strNotFound = "~~~Not found~~~"

// Init - class instance initialization
func (c *Cache) Init(MaxElementsFirstLavel int, MaxElementsSecondLavel int, ConnectionString string, CreateTable bool) (err error) {
	c.maxElementsFirstLavel = MaxElementsFirstLavel
	c.maxElementsSecondLavel = MaxElementsSecondLavel
	c.mapElements = make(map[string]*interface{}, MaxElementsFirstLavel)
	c.mapAllElements = make(map[string]*itemAllElements, MaxElementsFirstLavel+MaxElementsSecondLavel)
	c.connectionString = ConnectionString
	c.connection, err = sql.Open("sqlite3", c.connectionString)
	if err == nil {
		defer c.connection.Close()
		if CreateTable {
			_, err = c.connection.Exec("CREATE TABLE `Hache` ( `ID` INTEGER PRIMARY KEY AUTOINCREMENT, `Hache` TEXT NOT NULL UNIQUE, `Value` BLOB, `Old` INTEGER )")
			if err == nil {
				_, err = c.connection.Exec("CREATE INDEX `HacheByOld` ON `Hache` ( `Old` )")
				if err == nil {
					fmt.Printf("Таблица создана\n")
				}
			}
		}
	}

	if err != nil {
		fmt.Printf("%v", err)
	}

	return err
}

// Add - add to cache value with key
func (c *Cache) Add(key string, value interface{}) (err error) {
	for _, item := range c.mapAllElements {
		if math.MinInt32+1 <= item.Old {
			item.Old--
		}
	}
	// key already exists in hash
	if c.mapAllElements[key] != nil {
		if c.mapAllElements[key].Lavel == 1 {
			c.mapElements[key] = &value
			c.mapAllElements[key].Old = minOld
		} else {
			fmkey := strNotFound
			fmold := math.MaxInt32

			for key, item := range c.mapAllElements {
				// search item for down to second level
				if item.Old < fmold && item.Lavel == 1 {
					fmold = item.Old
					fmkey = key
				}
			}

			if fmkey != strNotFound {
				svalue := c.mapElements[fmkey]
				skey := fmkey
				sv := *svalue

				c.connection, err = sql.Open("sqlite3", c.connectionString)
				if err == nil {
					var result sql.Result
					defer c.connection.Close()
					result, err = c.connection.Exec("insert into `Hache` (`Hache`, `Value`) values ($1, $2)",
						skey, sv)
					fmt.Println(result.RowsAffected())
					fmt.Printf("Down %+v\n", skey, sv)
				}
			}
			c.mapElements[key] = &value
			c.mapAllElements[key].Old = minOld
			c.mapAllElements[key].Lavel = 1
		}
	} else {
		if len(c.mapElements) < c.maxElementsFirstLavel {
			c.mapElements[key] = &value
			item := new(itemAllElements)
			item.Old = minOld
			item.Lavel = 1
			c.mapAllElements[key] = item
			fmt.Printf("Add %+v with value %+v to lavel one\n", key, value)
		} else {
			fmkey := strNotFound
			fmold := math.MaxInt32
			smkey := strNotFound
			smold := math.MaxInt32
			for key, item := range c.mapAllElements {
				// search item for down to second level
				if item.Old < fmold && item.Lavel == 1 {
					fmold = item.Old
					fmkey = key
				}
				// search item for delete from second level
				if item.Old < smold && item.Lavel == 2 {
					smold = item.Old
					smkey = key
				}
			}

			if fmkey != strNotFound {
				value = c.mapElements[fmkey]
			}

		}
	}

	if err != nil {
		fmt.Printf("%v", err)
	}

	return err
}

// Get - get value from cache by key
func (c *Cache) Get(key string) (value *interface{}, err error) {
	v := c.mapAllElements[key]
	if v != nil {
		for _, item := range c.mapAllElements {
			if math.MinInt32+1 <= item.Old {
				item.Old--
			}
		}
		if v.Lavel == 1 {
			value = c.mapElements[key]
			if math.MaxInt32-10 > v.Old {
				v.Old += 10
			} else {
				v.Old = math.MaxInt32
			}
		} else {
			c.connection, err = sql.Open("sqlite3", c.connectionString)
			if err == nil {
				defer c.connection.Close()
				rows, err := c.connection.Query("select `Value` from Hache where `Hache`=$1", key)
				if err == nil {
					defer rows.Close()
					if rows.Next() {
						err = rows.Scan(value)
					}
				}
			}
		}
	} else {
		err = errors.New("Not found")
	}

	if err != nil {
		fmt.Printf("%v", err)
	}

	return value, err
}
