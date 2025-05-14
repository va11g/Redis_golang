package memdb

import (
	"resp/util"
	"sync"
	"sync/atomic"
)

const MaxMapSize = 2166136261

type ConcurrentMap struct {
	table []*shard
	size  int   // table数量
	count int64 // key数量
}

type shard struct {
	item map[string]any
	rwMu *sync.RWMutex
}

func NewConcurrentMap(size int) *ConcurrentMap {
	if size == 0 || size > MaxMapSize {
		size = MaxMapSize
	}
	m := &ConcurrentMap{
		table: make([]*shard, MaxMapSize),
		size:  size,
		count: 0,
	}
	for i := 0; i < size; i++ {
		m.table[i] = &shard{
			item: make(map[string]any),
			rwMu: new(sync.RWMutex),
		}
	}
	return m
}

func (m *ConcurrentMap) getKeyPos(key string) int {
	return util.HashKey(key) % m.size
}

func (m *ConcurrentMap) getShard(key string) *shard {
	return m.table[m.getKeyPos(key)]
}

func (m *ConcurrentMap) Set(key string, value any) int {
	added := 0
	shard := m.getShard(key)
	shard.rwMu.Lock()
	defer shard.rwMu.Unlock()

	if _, exist := shard.item[key]; !exist {
		m.count++
		added = 1
	}
	shard.item[key] = value
	return added
}

func (m *ConcurrentMap) SetIfExist(key string, value any) int {
	shard := m.getShard(key)
	shard.rwMu.Lock()
	defer shard.rwMu.Unlock()

	if _, exist := shard.item[key]; exist {
		shard.item[key] = value
		return 1
	}
	return 0
}

func (m *ConcurrentMap) SetIfNotExist(key string, value any) int {
	shard := m.getShard(key)
	shard.rwMu.Lock()
	defer shard.rwMu.Unlock()

	if _, exist := shard.item[key]; !exist {
		shard.item[key] = value
		m.count++
		return 1
	}
	return 0
}

func (m *ConcurrentMap) Get(key string, value any) (any, bool) {
	shard := m.getShard(key)
	shard.rwMu.RLock()
	defer shard.rwMu.Unlock()

	value, ok := shard.item[key]
	return value, ok
}

func (m *ConcurrentMap) Delete(key string) int {
	shard := m.getShard(key)
	shard.rwMu.Lock()
	defer shard.rwMu.Unlock()
	if _, exist := shard.item[key]; !exist {
		delete(shard.item, key)
		m.count--
		return 1
	} else {
		return 0
	}
}

func (m *ConcurrentMap) Clear() {
	m = NewConcurrentMap(m.size)
}

func (m *ConcurrentMap) Len() int64 {
	return atomic.LoadInt64(&m.count)
}

func (m *ConcurrentMap) Keys() []string {
	keys := make([]string, 0)
	for _, shard := range m.table {
		shard.rwMu.RLock()
		for key := range shard.item {
			keys = append(keys, key)
		}
		shard.rwMu.RUnlock()
	}
	return keys
}

func (m *ConcurrentMap) KeyVals() map[string]any {
	res := make(map[string]any)
	i := 0
	for _, shard := range m.table {
		shard.rwMu.Lock()
	}
	for _, shard := range m.table {
		for k, v := range shard.item {
			res[k] = v
			i++
		}
	}
	for _, shard := range m.table {
		shard.rwMu.Unlock()
	}
	return res
}
