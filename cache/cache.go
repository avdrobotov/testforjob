package cache

// Cache - class of cache
type Cache struct {
	maxElementsFirstLavel  int
	maxElementsSecondLavel int
	mapElements            *map[string]*interface{}
}
