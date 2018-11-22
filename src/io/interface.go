package io

// Address contains: file idx, file Offset
type Address 	uint32
// Size is the value length in file
type Size		uint32

type Mem interface {
	Write(data []byte) (address Address, size Size)
	Read(address Address, size Size) (data []byte)
	Delete(address Address, size Size)
}

type Recycle interface {
	Add(address Address, size Size)
	Get(need Size) (address Address, size Size)
}
