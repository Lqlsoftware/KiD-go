package core

type KiDInterface struct {}

func (KiD *KiDInterface)Get(key string) {}

func (KiD *KiDInterface)Put(key string, value string) {}

func (KiD *KiDInterface)Delete(key string) {}