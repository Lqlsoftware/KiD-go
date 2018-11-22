package core

import (
	"github.com/Lqlsoftware/KiD/src/cmap"
	"github.com/Lqlsoftware/KiD/src/conf"
	"github.com/Lqlsoftware/KiD/src/io"
)

type KiDCore struct {
	cMap 	*cmap.ConcurrentMap
	memory 	*io.KidMemory

}

func (KiD *KiDCore)Init(config *conf.KiDConfig) {
	if config == nil {
		config = conf.DefaultConfig
	}
	KiD.cMap.Init(config)
	KiD.memory.Init(config)
}

func (KiD *KiDCore)Get(key string) (value string, err error) {
	// fetch from cmap
	var keyHashed cmap.MapKey = BKDRHash(key)
	var mapValue *cmap.MapValue = KiD.cMap.Get(keyHashed)
	// Not found
	if mapValue == nil {
		// TODO add error ERR_NOT_FOUND
		return "", nil
	}
	// fetch from I/O
	// TODO add error ERR_READ_IO
	var rawValue []byte = KiD.memory.Read(mapValue.Address, mapValue.Length)
	return string(rawValue), nil
}

func (KiD *KiDCore)Put(key string, value string) (err error) {
	// fetch from cmap
	var keyHashed cmap.MapKey = BKDRHash(key)
	var oldValue *cmap.MapValue = KiD.cMap.Get(keyHashed)
	// write to I/O
	// TODO add error ERR_WRITE_IO
	var address, size = KiD.memory.Write([]byte(value))
	var mapValue = &cmap.MapValue{
		Address:address,
		Length:size,
	}
	// TODO add error ERR_WRITE_CMAP
	KiD.cMap.Put(keyHashed, mapValue)
	// overwrite
	if oldValue != nil {
		KiD.memory.Delete(oldValue.Address, oldValue.Length)
	}
	return nil
}

func (KiD *KiDCore)Delete(key string) (value string, err error) {
	// fetch from cmap
	var keyHashed cmap.MapKey = BKDRHash(key)
	var mapValue *cmap.MapValue = KiD.cMap.Get(keyHashed)
	// Not found
	if mapValue == nil {
		// TODO add error ERR_NOT_FOUND
		return "", nil
	}
	KiD.cMap.Delete(keyHashed)
	KiD.memory.Delete(mapValue.Address, mapValue.Length)
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