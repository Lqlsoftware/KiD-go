package io

// Address contains: file idx, file Offset
type Address 	uint32
// Size is the value length in file
type Size		uint32

type IO interface {
	Write(data []byte) (address Address, size Size)
	Read(address Address, size Size) (data []byte)
}