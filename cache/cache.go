/**------------------------------------------------------------**
 * @file     jinycoo.com/cache.go
 * @version  1.0.0
 * @author   jinycoo - caojingyin@jinycoo.com
 * @date     2020/8/12 9:48
 * @desc     jinycoo.com - main - summary
 **------------------------------------------------------------**/

package cache

import (
	"errors"
)

type Cache interface {
	Set(key, value string) bool
	Get(key string) interface{}
}

type RedisCache struct {
	data map[string]string
}

type MemCache struct {
	data map[string]string
}

func (redis *RedisCache) Set(key, value string) bool {
	redis.data[key] = value
	return true
}

func (redis *RedisCache) Get(key string) interface{} {
	return "redis:" + redis.data[key]
}

func (mem *MemCache) Set(key, value string) bool {
	mem.data[key] = value
	return true
}

func (mem *MemCache) Get(key string) interface{} {
	return "mem:" + mem.data[key]
}

type CacheType int

const (
	redis CacheType = iota
	mem
)

//type DataStoreFactory func(conf map[string]string) (Cache, error)
//
//func NewPostgreSQLDataStore(conf map[string]string) (Cache, error) {
//	dsn, ok := conf.Get("DATASTORE_POSTGRES_DSN", "")
//	if !ok {
//		return nil, errors.New(fmt.Sprintf("%s is required for the postgres datastore", "DATASTORE_POSTGRES_DSN"))
//	}
//
//	db, err := sqlx.Connect("postgres", dsn)
//	if err != nil {
//		log.Panicf("Failed to connect to datastore: %s", err.Error())
//		return nil, datastore.FailedToConnect
//	}
//
//	return &PostgresDataStore{
//		DSN: dsn,
//		DB:  db,
//	}, nil
//}
//
//func NewMemoryDataStore(conf map[string]string) (DataStore, error) {
//	return &MemoryDataStore{
//		Users: &map[int64]string{
//			1: "mnbbrown",
//			0: "root",
//		},
//		RWMutex: &sync.RWMutex{},
//	}, nil
//}

type CacheFactory struct {
}

func (factory *CacheFactory) Create(cacheType CacheType) (Cache, error) {
	if cacheType == redis {
		var dt = make(map[string]string, 0)
		return &RedisCache{
			data: dt,
		}, nil
	}
	if cacheType == mem {
		var dt = make(map[string]string, 0)
		return &MemCache{
			data: dt,
		}, nil
	}

	return nil, errors.New("error cache type")
}
