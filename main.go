// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"

	"github.com/avdrobotov/testforjob/cache"
	_ "github.com/avdrobotov/testforjob/cache"
)

func main() {
	var err error
	fmt.Println("Начнём!!!")
	var c *cache.Cache
	c = new(cache.Cache)
	err = c.Init(10, 10, "data.db", false)
	if err == nil {
		c.Add("1", 1)
		c.Add("2", 2)
		c.Add("3", 3)
		c.Add("4", 4)
		c.Add("5", 5)
		c.Add("6", 6)
		c.Add("7", 7)
		c.Add("8", 8)
		c.Add("9", 9)
		c.Add("0", 0)
		c.Add("1", 1)
		c.Add("1", 1)
		c.Add("1", 1)
		c.Add("1", 1)
		c.Add("1", 1)
	}
}
