package cmap

import (
	"sync"

	"github.com/Lqlsoftware/KiD/src/conf"
)

type MapKey 	uint32
// TODO
type MapValue 	interface {}
type MapData struct {
	Value MapValue
}

// ConcurrentMap
// Divide Map to CMAP_BLOCK_NUM block
// each block is RW independent
type ConcurrentMap struct {
	base []map[MapKey]*Node
	lock []*sync.RWMutex
}

// Init cMap
// cMap have CMAP_BLOCK_NUM block_maps
// each block_map have a rw_locker and initial space of CMAP_BLOCK_INIT_SIZE
func (cMap *ConcurrentMap)Init(conf conf.KiDConfig) {
	cMap.base = make([]map[MapKey]*Node, conf.CMAP_BLOCK_NUM)
	cMap.lock = make([]*sync.RWMutex, conf.CMAP_BLOCK_NUM)
	// allocated map & rw_locker
	for i := uint8(0);i < conf.CMAP_BLOCK_NUM;i++ {
		cMap.base[i] = make(map[MapKey]*Node, conf.CMAP_BLOCK_INIT_SIZE)
		cMap.lock[i] = new(sync.RWMutex)
	}
}

// Put a uint32-key and node to cMap
// The key should be already hashed to uint32
// base map index is the highest 4 digits (key >> 28)
// To reduce conflict, make list node to memory key and value
func (cMap *ConcurrentMap)Put(key MapKey, value MapValue) {
	// put
	var idx = key >> 28
	cMap.lock[idx].Lock()
	// TODO I/O
	cMap.lock[idx].Unlock()
}

// Get a uint32-key's value
// The key should be already hashed to uint32
// base map index is the highest 4 digits (key >> 28)
// search from
func (cMap *ConcurrentMap)Get(key MapKey) interface{} {
	idx := key >> 28
	cMap.lock[idx].RLock()
	cMap.lock[idx].RUnlock()
	// TODO I/O
	return nil
}

// Put a uint32-key and node to cMap
// The key should be already hashed to uint32
// base map index is the highest 4 digits (key >> 28)
// To reduce conflict, make list node to memory key and value
func (cMap *ConcurrentMap)Delete(key MapKey) interface{} {
	idx := key >> 28
	cMap.lock[idx].RLock()
	cMap.lock[idx].RUnlock()
	// TODO I/O
	return nil
}