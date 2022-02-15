package yunjing

import (
	"github.com/gavinlhchen/logconvert/errors"
	"github.com/gavinlhchen/logconvert/log"
	"sync"
)

var (
	providersMu sync.RWMutex
	providers   = make(map[string]YjToSocMsgHandler)
)

func Register(name string, c YjToSocMsgHandler) {
	providersMu.Lock()
	defer providersMu.Unlock()

	if c == nil {
		log.Panicf("yjconsumer:Register provider is nil")
	}

	if _, dup := providers[name]; dup {
		log.Panicf("yjconsumer:Register called twice for provider: %s", name)
	}
	providers[name] = c
}

func New(providerName string) (YjToSocMsgHandler, error) {
	providersMu.RLock()
	p, ok := providers[providerName]
	providersMu.RUnlock()

	if !ok {
		return nil, errors.Errorf("yjconsumer:unknown provider %s", providerName)
	}
	return p, nil
}
