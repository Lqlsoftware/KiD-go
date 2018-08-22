package core

type KiDInterface interface {
	Get(key string) (value string)
	Put(key string, value string)
	Delete(key string) (value string)
}