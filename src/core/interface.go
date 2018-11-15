package core

type KiDInterface interface {
	Get(key string) (value string)
	Put(key string, value string) (err error)
	Delete(key string) (value string)
}