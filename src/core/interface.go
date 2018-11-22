package core

type Core interface {
	Get(key string) (value string, err error)
	Put(key string, value string) (err error)
	Delete(key string) (value string, err error)
}