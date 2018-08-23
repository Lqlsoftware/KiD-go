package cmap

import (
	"sync"

	"github.com/Lqlsoftware/KiD/src/conf"
	)

type MapKey 	uint32
// TODO
type MapValue 	string
type MapData struct {
	Value 		MapValue
	Address 	uint32
	Length		uint32
}

// ConcurrentMap
// Divide Map to CMAP_BLOCK_NUM block
// each block is RW independent
type ConcurrentMap struct {
	// Private
	base		[]map[MapKey]*RBTree
	lock 		[]*sync.RWMutex
	shiftLength uint8
	// Public
	Size		uint32
}

// Init cMap
// cMap have CMAP_BLOCK_NUM block_maps
// each block_map have a rw_locker and initial space of CMAP_BLOCK_INIT_SIZE
func (cMap *ConcurrentMap)Init(conf conf.KiDConfig) {
	cMap.base = make([]map[MapKey]*RBTree, conf.CMAP_BLOCK_NUM)
	cMap.lock = make([]*sync.RWMutex, conf.CMAP_BLOCK_NUM)
	// allocated map & rw_locker
	for i := uint8(0);i < conf.CMAP_BLOCK_NUM;i++ {
		cMap.base[i] = make(map[MapKey]*RBTree, conf.CMAP_BLOCK_INIT_SIZE)
		cMap.lock[i] = new(sync.RWMutex)
	}
	// calculate shift length
	var i uint8 = 0
	for t := conf.CMAP_BLOCK_NUM;t > 0;t >>= 1 {
		i++
	}
	cMap.shiftLength = 32 - i
}

// Put a uint32-key and node to cMap
// The key should be already hashed to uint	32
// base map index is the highest 4 (16 = 2^4) digits (key >> 28)
// To reduce conflict, make treeNode to memory key and value
// write to cache after operate in treeNode
func (cMap *ConcurrentMap)Put(key MapKey, value MapValue) {
	// put
	var idx = uint8(key >> cMap.shiftLength)
	var _map = cMap.base[idx]
	var data *MapData = nil
	// -----------------------------------------------
	// Write Lock
	//
	cMap.lock[idx].Lock()
	// bucket not have been access
	if _, ok := _map[key];ok {
		_map[key] = NewTree(idx)
	}
	// already have value in I/O
	// delete it first
	data = _map[key].Get(key)
	if data != nil {
		// TODO I/O write and gc

	}
	// write to I/O
	// TODO I/O write

	_map[key].Put(key, &MapData{Value:value})
	cMap.lock[idx].Unlock()
	//
	// Write Unlock
	// -----------------------------------------------
}

// Get a uint32-key's value
// The key should be already hashed to uint32
// base map index is the highest 4 (16 = 2^4) digits (key >> 28)
// search from treeNode in map[key]
func (cMap *ConcurrentMap)Get(key MapKey) *MapData {
	var idx = uint8(key >> cMap.shiftLength)
	var data *MapData = nil
	// -----------------------------------------------
	// Read Lock
	//
	cMap.lock[idx].RLock()
	if list, ok := cMap.base[idx][key];ok {
		data = list.Get(key)
		// TODO I/O read

	}
	cMap.lock[idx].RUnlock()
	//
	// Read Unlock
	// -----------------------------------------------
	return data
}

// Put a uint32-key and node to cMap
// The key should be already hashed to uint32
// base map index is the highest 4 (16 = 2^4) digits (key >> 28)
// delete from treeNode in map[key]
// write to cache after operate in treeNode
func (cMap *ConcurrentMap)Delete(key MapKey) *MapData {
	var idx = uint8(key >> cMap.shiftLength)
	var data *MapData = nil
	// -----------------------------------------------
	// Write Lock
	//
	cMap.lock[idx].Lock()
	if list, ok := cMap.base[idx][key];ok {
		data = list.Delete(key)
		// TODO I/O write and gc
	}
	cMap.lock[idx].Unlock()
	//
	// Write Unlock
	// -----------------------------------------------
	return data
}