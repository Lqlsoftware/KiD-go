package io

import (
	"container/list"
	"sync"
)

type Buffer struct {
	RWMem 		*map[Address][]byte
	RMem		*list.List
	RWMemlock	*sync.RWMutex
	RMemlock	*sync.RWMutex
	size		uint32
}

// Initial cache with a sized RW map and a RW locker
func (buffer *Buffer)Init(size uint32) {
	buffer.size = size
	RWMem := make(map[Address][]byte, size)
	buffer.RWMem = &RWMem
	buffer.RMem = list.New()
	buffer.RWMemlock = new(sync.RWMutex)
	buffer.RMemlock = new(sync.RWMutex)
}

// Write data to cache
func (buffer *Buffer)Write(data []byte) (address Address, size Size) {
	// TODO get space from space manager => address, size

	// Write to RWMem
	{
		buffer.RWMemlock.Lock()
		// get a new RWMem and Write old one to I/O
		if uint32(len(*buffer.RWMem)) == buffer.size {

		}
		(*buffer.RWMem)[address] = data
		buffer.RWMemlock.Unlock()
	}
	return 0,0
}

// Read data from cache
func (buffer *Buffer)Read(address Address, size Size) (data []byte) {
	var value []byte
	var bufferHited, useRMem = false, false
	// Read from RWMem
	{
		buffer.RWMemlock.RLock()
		// read RWMem
		value, bufferHited = (*buffer.RWMem)[address]
		// RWMem missed, try reading from RMem
		useRMem = !bufferHited && buffer.RMem != nil
		buffer.RWMemlock.RUnlock()
	}
	if useRMem {
		buffer.RMemlock.RLock()
		var RMap *map[Address][]byte
		for m := buffer.RMem.Back();m != buffer.RMem.Front().Prev() && !bufferHited;m = m.Prev() {
			RMap = m.Value.(*map[Address][]byte)
			value, bufferHited = (*RMap)[address]
		}
		buffer.RMemlock.RUnlock()
	}
	// all missed. try i/o
	if !bufferHited {
		// TODO
	}
	return value
}