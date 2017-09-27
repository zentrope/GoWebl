//
// Copyright (c) 2017 Keith Irwin
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published
// by the Free Software Foundation, either version 3 of the License,
// or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package internal

import (
	"net"
	"sync"
	"sync/atomic"
)

// alternatives:
//   https://github.com/orcaman/concurrent-map
//   https://github.com/patrickmn/go-cache
//
// but rolling my own naive versions just to keep
// see what it's like

type GenericMap map[string]interface{}

type CacheTransactor func() (interface{}, error)

type GenericCache struct {
	data  atomic.Value
	mutex sync.Mutex
}

func NewCache() *GenericCache {
	var data atomic.Value
	data.Store(make(GenericMap))

	return &GenericCache{data, sync.Mutex{}}
}

func (c *GenericCache) Get(key string) (interface{}, bool) {
	values := c.data.Load().(GenericMap)
	value, found := values[key]
	return value, found
}

func (c *GenericCache) Set(key string, val interface{}) interface{} {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	values := c.data.Load().(GenericMap)
	values[key] = val
	c.data.Store(values)
	return val
}

func (c *GenericCache) Transact(key string, tx CacheTransactor) (interface{}, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	values := c.data.Load().(GenericMap)

	if value, found := values[key]; found {
		return value, nil
	}

	val, err := tx()
	if err != nil {
		return nil, err
	}

	values[key] = val
	c.data.Store(values)
	return val, nil
}

func (c *GenericCache) GetOrSet(key string, tx CacheTransactor) (interface{}, error) {
	if val, found := c.Get(key); found {
		return val, nil
	}
	return c.Transact(key, tx)
}

var DNS_CACHE = NewCache()

func DNSLookup(address string) string {

	if address == "0.0.0.0" {
		return address
	}

	if name, found := DNS_CACHE.Get(address); found {
		return name.(string)
	}

	names, err := net.LookupAddr(address)
	if err != nil {
		return DNS_CACHE.Set(address, address).(string)
	}

	if len(names) == 0 {
		return DNS_CACHE.Set(address, address).(string)
	}

	return DNS_CACHE.Set(address, names[0]).(string)
}
