package io

import (
	"sync"

	"github.com/Lqlsoftware/KiD/src/conf"
)

type KidRecycler struct {
	GCLock		*sync.RWMutex
	SleepTime	uint32
	BlockSize	uint32
}

func (r KidRecycler)Init(config *conf.KiDConfig) {
	panic("implement me")
}

func (r KidRecycler)Add(address Address, size Size) {
	// TODO add unused space
	panic("implement me")
}

func (r KidRecycler)Get(need Size) (address Address, size Size) {
	// TODO return need space
	panic("implement me")
}

func (r KidRecycler)RecycleThread() {
	// TODO hdd space gc
}