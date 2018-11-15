package conf

type KiDConfig struct {
	CMAP_BLOCK_NUM 			uint8
	CMAP_BLOCK_INIT_SIZE 	uint32

	BUFFER_MAP_INIT_SIZE 	uint32
}

var DefaultConfig = &KiDConfig {
	CMAP_BLOCK_NUM:			16,
	CMAP_BLOCK_INIT_SIZE:	1024,

	BUFFER_MAP_INIT_SIZE:	1024,
}