package cache

import (
	"encoding/json"
	"fmt"
	"runtime/debug"

	"../config"
	"../model"
	"github.com/coocood/freecache"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("cache_manager")

var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

//Manager manages cache initialization and cache related operations
type Manager struct {
	cacheconfig config.CacheConfiguration
	cache       *freecache.Cache
}

//Init intializes cache configuration and allocates memory for cache items
func (m *Manager) Init(cacheconfig config.CacheConfiguration) error {
	m.cacheconfig = cacheconfig
	m.cache = freecache.NewCache(m.cacheconfig.Memory * 1024 * 1024)
	debug.SetGCPercent(10)
	return nil
}

//AddItem adds the item into the cache
func (m *Manager) AddItem(employee model.Employee) error {
	b, err := json.Marshal(employee)
	if err != nil {
		log.Error("Error marshalling employee: ", err)
		return err
	}
	err = m.cache.Set([]byte(employee.ID), b, 0)
	if err != nil {
		log.Error("Error cache set: ", err)
		return err
	}
	return nil
}

//GetItem gets the item from the cache
func (m *Manager) GetItem(ID string) (*model.Employee, error) {
	var employee model.Employee
	b, err := m.cache.Get([]byte(ID))
	if err != nil {
		log.Error("Error cache get: ", err)
		return nil, err
	}
	err = json.Unmarshal(b, &employee)
	if err != nil {
		fmt.Println("Error unmarshalling employee:", err)
		return nil, err
	}
	return &employee, nil
}
