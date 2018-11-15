package cmap

import (
	"sync"

	"github.com/Lqlsoftware/KiD/src/conf"
	"github.com/Lqlsoftware/KiD/src/io"
)

// Key && Value define
type MapKey 	uint32
type MapValue struct {
	Address 	io.Address
	Length		io.Size
}

// ConcurrentMap
// Divide Map to CMAP_BLOCK_NUM block
// each block is RW independent
type ConcurrentMap struct {
	base []map[MapKey]*RBTree
	lock []*sync.RWMutex
}

// Init cMap
// cMap have CMAP_BLOCK_NUM block_maps
// each block_map have a rw_locker and initial space of CMAP_BLOCK_INIT_SIZE
func (cMap *ConcurrentMap)Init(config *conf.KiDConfig) {
	cMap.base = make([]map[MapKey]*RBTree, config.CMAP_BLOCK_NUM)
	cMap.lock = make([]*sync.RWMutex, config.CMAP_BLOCK_NUM)
	// allocated map & rw_locker
	for i := uint8(0);i < config.CMAP_BLOCK_NUM;i++ {
		cMap.base[i] = make(map[MapKey]*RBTree, config.CMAP_BLOCK_INIT_SIZE)
		cMap.lock[i] = new(sync.RWMutex)
	}
}

// Put a uint32-key and node to cMap
// The key should be already hashed to uint32
// base map index is the highest 4 (16 = 2^4) digits (key >> 28)
// To reduce conflict, make treeNode to memory key and value
// write to cache after operate in treeNode
func (cMap *ConcurrentMap)Put(key MapKey, value *MapValue) (err error) {
	// put
	var idx = key >> 28
	cMap.lock[idx].Lock()
	cMap.base[idx][key].Put(key, value)
	cMap.lock[idx].Unlock()
}

// Get a uint32-key's value
// The key should be already hashed to uint32
// base map index is the highest 4 (16 = 2^4) digits (key >> 28)
// search from treeNode in map[key]
func (cMap *ConcurrentMap)Get(key MapKey) *MapValue {
	var idx = key >> 28
	cMap.lock[idx].RLock()
	var value = cMap.base[idx][key].Get(key)
	cMap.lock[idx].RUnlock()
	return value
}

// Put a uint32-key and node to cMap
// The key should be already hashed to uint32
// base map index is the highest 4 (16 = 2^4) digits (key >> 28)
// delete from treeNode in map[key]
// write to cache after operate in treeNode
func (cMap *ConcurrentMap)Delete(key MapKey) *MapValue {
	var idx = key >> 28
	cMap.lock[idx].Lock()
	var value = cMap.base[idx][key].Delete(key)
	cMap.lock[idx].Unlock()
	return value
}