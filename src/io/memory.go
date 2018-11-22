package io

import (
	"container/list"
	"sync"

	"github.com/Lqlsoftware/KiD/src/conf"
)

type KidMemory struct {
	RWMem 		*map[Address][]byte
	RMem		*list.List
	RWMemlock	*sync.RWMutex
	RMemlock	*sync.RWMutex
	size		uint32
}

// Initial cache with a sized RW map and a RW locker
func (m *KidMemory)Init(config *conf.KiDConfig) {
	m.size = config.BUFFER_MAP_INIT_SIZE
	RWMem := make(map[Address][]byte,  config.BUFFER_MAP_INIT_SIZE)
	m.RWMem = &RWMem
	m.RMem = list.New()
	m.RWMemlock = new(sync.RWMutex)
	m.RMemlock = new(sync.RWMutex)
}

// Write data to cache
func (m *KidMemory)Write(data []byte) (address Address, size Size) {
	// TODO get space from space manager => address, size

	// Write to RWMem
	{
		m.RWMemlock.Lock()
		// get a new RWMem and Write old one to I/O
		if uint32(len(*m.RWMem)) == m.size {
			// TODO
		}
		(*m.RWMem)[address] = data
		m.RWMemlock.Unlock()
	}
	return 0,0
}

// Read data from cache
func (m *KidMemory)Read(address Address, size Size) (data []byte) {
	var value []byte
	var bufferHited, useRMem = false, false
	// Read from RWMem
	{
		m.RWMemlock.RLock()
		// read RWMem
		value, bufferHited = (*m.RWMem)[address]
		// RWMem missed, try reading from RMem
		useRMem = !bufferHited && m.RMem != nil
		m.RWMemlock.RUnlock()
	}
	if useRMem {
		{
			m.RMemlock.RLock()
			var RMap *map[Address][]byte
			for m := m.RMem.Back(); m != m.RMem.Front().Prev() && !bufferHited; m = m.Prev() {
				RMap = m.Value.(*map[Address][]byte)
				value, bufferHited = (*RMap)[address]
			}
			m.RMemlock.RUnlock()
		}
	}
	// all missed. try i/o
	if !bufferHited {
		// TODO
	}
	return value
}

func (m *KidMemory)Delete(address Address, size Size) {

}