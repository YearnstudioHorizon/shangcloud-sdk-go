package shangcloudsdkgo

import (
	"fmt"
	"sync"
)

// TempVatStorage接口 需要内部实现线程安全
type TempVarStorage interface {
	SetTempVarible(string, string)
	GetTempVarible(string) (string, error)
	DeleteTempVarible(string)
}

// 由sync.Map实现的线程安全的内存KV存储
type ramKv struct {
	storage *sync.Map
}

func newRamKv() *ramKv {
	return &ramKv{
		storage: &sync.Map{},
	}
}

func (self *ramKv) SetTempVarible(key string, value string) {
	self.storage.Store(key, value)
}

func (self *ramKv) GetTempVarible(key string) (string, error) {
	raw, ok := self.storage.Load(key)
	if !ok {
		return "", fmt.Errorf("Get value of %v failed", key)
	}
	v, ok := raw.(string)
	if !ok {
		return "", fmt.Errorf("Convert value of %v to string failed", key)
	}
	return v, nil
}

func (self *ramKv) DeleteTempVarible(key string) {
	self.storage.Delete(key)
}
