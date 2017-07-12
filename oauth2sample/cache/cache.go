package cache

import (
	"github.com/patrickmn/go-cache"
)

//declare global variable for cache
var C = cache.New(cache.NoExpiration, cache.NoExpiration)

/*
 * Method to add value to cache
 */
func AddToCache(key, value string) {

	C.Set(key, value, cache.NoExpiration)
}

/*
 * Method to retrieve value from cache
 */
func GetFromCache(key string) string {
	var value string
	result, found := C.Get(key)
	if found {
		value = result.(string)
	}
	return value
}
