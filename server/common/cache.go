package common

import "time"

func GetCache(k string) (interface{}, error) {
	return Cache.Get(k)
}
func SetCache(k string, v interface{}, t time.Duration) error {
	return Cache.Set(k, v, t)
}
func DelCache(k string) error {
	return Cache.Delete(k)
}
