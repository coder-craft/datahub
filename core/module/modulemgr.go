package module

import (
	"../common"
	"../common/zlog"
	"time"
)

type ModuleManager struct {
	pool map[string]ModuleEntity
}

var moduleMgr = &ModuleManager{
	pool: make(map[string]ModuleEntity),
}

func RegisterModule(m ModuleInterface, interval time.Duration) {
	moduleMgr.pool[m.Name()] = ModuleEntity{
		signal:   make(chan bool),
		interval: interval.Nanoseconds(),
		module:   m,
	}
}
func ModuleStart() {
	for _, value := range moduleMgr.pool {
		value.module.Init()
		go moduleLoop(value)
	}
}
func ModuleStop() {
	for _, value := range moduleMgr.pool {
		value.signal <- true
	}
}
func moduleLoop(m ModuleEntity) {
	for {
		select {
		case <-m.signal:
			m.module.End()
			zlog.Info("Module Stop.", zlog.String("Name", m.module.Name()))
			return
		default:
			now := time.Now().UnixNano()
			if now-m.lastTs >= m.interval {
				_ = common.SafeCall(func(args ...interface{}) (err error) {
					if len(args) > 0 {
						module, ok := args[0].(ModuleInterface)
						if ok {
							module.Update()
						}
					}
					return nil
				}, m.module)
				m.lastTs = now
			} else {
				if (m.interval / 5) > time.Second.Nanoseconds()*5 {
					time.Sleep(time.Second * 5)
				} else {
					time.Sleep(time.Duration(m.interval / 5))
				}
			}
		}
	}
}
