package core

import (
	"github.com/Lqlsoftware/KiD/src/cmap"
	"github.com/Lqlsoftware/KiD/src/conf"
	"github.com/Lqlsoftware/KiD/src/io"
)

type KiDCore struct {
	cMap 	*cmap.ConcurrentMap
	memory 	*io.Memory
}

func (KiD *KiDCore)Init(config *conf.KiDConfig) {
	if config == nil {
		config = conf.DefaultConfig
	}
	KiD.cMap.Init(config)
	KiD.memory.Init(config)
}

func (KiD *KiDCore)Get(key string) (value string) {
	// fetch from cmap
	var keyHashed cmap.MapKey = BKDRHash(key)
	var mapValue *cmap.MapValue = KiD.cMap.Get(keyHashed)
	if mapValue == nil {
		// TODO
	}
	// fetch from I/O
	var rawValue []byte = KiD.memory.Read(mapValue.Address, mapValue.Length)
	return string(rawValue)
}

func (KiD *KiDCore)Put(key string, value string) (err error) {
	// fetch from cmap
	var keyHashed cmap.MapKey = BKDRHash(key)
	var mapValue *cmap.MapValue = KiD.cMap.Get(keyHashed)
	if mapValue != nil {
		// TODO
	}
	// write to I/O
	var address, size = KiD.memory.Write()

}

func (KiD *KiDCore)Delete(key string) (value string) {

}

func BKDRHash(str string) (key cmap.MapKey) {
	var seed uint32 = 131
	var hash uint32 = 0
	for _, v := range str {
		hash = hash * seed + uint32(v)
	}
	// highest digit => 0
	return cmap.MapKey(hash & 0x7FFFFFFF)
}