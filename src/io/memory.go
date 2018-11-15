package io

import (
	"container/list"
	"sync"

	"github.com/Lqlsoftware/KiD/src/conf"
)

type Memory struct {
	RWMem 		*map[Address][]byte
	RMem		*list.List
	RWMemlock	*sync.RWMutex
	RMemlock	*sync.RWMutex
	size		uint32
}

// Initial cache with a sized RW map and a RW locker
func (memory *Memory)Init(config *conf.KiDConfig) {
	memory.size = config.BUFFER_MAP_INIT_SIZE
	RWMem := make(map[Address][]byte,  config.BUFFER_MAP_INIT_SIZE)
	memory.RWMem = &RWMem
	memory.RMem = list.New()
	memory.RWMemlock = new(sync.RWMutex)
	memory.RMemlock = new(sync.RWMutex)
}

// Write data to cache
func (memory *Memory)Write(data []byte) (address Address, size Size) {
	// TODO get space from space manager => address, size

	// Write to RWMem
	{
		memory.RWMemlock.Lock()
		// get a new RWMem and Write old one to I/O
		if uint32(len(*memory.RWMem)) == memory.size {
			// TODO
		}
		(*memory.RWMem)[address] = data
		memory.RWMemlock.Unlock()
	}
	return 0,0
}

// Read data from cache
func (memory *Memory)Read(address Address, size Size) (data []byte) {
	var value []byte
	var bufferHited, useRMem = false, false
	// Read from RWMem
	{
		memory.RWMemlock.RLock()
		// read RWMem
		value, bufferHited = (*memory.RWMem)[address]
		// RWMem missed, try reading from RMem
		useRMem = !bufferHited && memory.RMem != nil
		memory.RWMemlock.RUnlock()
	}
	if useRMem {
		memory.RMemlock.RLock()
		var RMap *map[Address][]byte
		for m := memory.RMem.Back();m != memory.RMem.Front().Prev() && !bufferHited;m = m.Prev() {
			RMap = m.Value.(*map[Address][]byte)
			value, bufferHited = (*RMap)[address]
		}
		memory.RMemlock.RUnlock()
	}
	// all missed. try i/o
	if !bufferHited {
		// TODO
	}
	return value
}