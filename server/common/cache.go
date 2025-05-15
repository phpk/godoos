package common

import "time"

func GetCache(k string) (interface{}, error) {
	return Cache.Get(k)
}
func GetCacheVal(cacheDataByte interface{}) []byte {
	var dataBytes []byte
	switch v := cacheDataByte.(type) {
	case string:
		dataBytes = []byte(v)
	case []byte:
		dataBytes = v
	default:
		return dataBytes
	}
	return dataBytes
}
func SetCache(k string, v interface{}, t time.Duration) error {
	return Cache.Set(k, v, t)
}
func DelCache(k string) error {
	return Cache.Delete(k)
}
