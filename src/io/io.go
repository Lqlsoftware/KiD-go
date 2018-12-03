package io

import (
	"os"

	"github.com/Lqlsoftware/KiD/src/conf"
)

type KidIO struct {
	blockSize	uint32
	files 		[]*os.File
	fileSize	uint32
	fileBitmap	[][]uint64
}

func (io KidIO)Init(config conf.KiDConfig) {

}

func (io KidIO)Write(data []byte, offset int64) (err error) {
	io.files[0].WriteAt(data, offset)
}

func parseAddress(address Address) {
	fileIdx := (uint32(address) & (uint32(31) << 28)) >> 28
}