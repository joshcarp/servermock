package data

import "sync"

func StoreData(sm *sync.Map, trace string, data Request) {
	val, ok := sm.Load(trace)
	if !ok {
		sm.Store(trace, []Request{data})
		return
	}
	valBytes := val.([]Request)
	if len(valBytes) == 0 {
		sm.Store(trace, []Request{data})
		return
	}
	sm.Store(trace, append([]Request{data}, valBytes...)) //[n-1] is always the element to be read (and deleted) first
}

func LoadData(sm *sync.Map, trace string) (Request, bool) {
	val, ok := sm.Load(trace)
	if !ok {
		return Request{}, false
	}
	valBytes := val.([]Request)
	if len(valBytes) == 0 {
		return Request{}, false
	}
	sm.Store(trace, valBytes[:len(valBytes)-1])
	return valBytes[len(valBytes)-1], true
}
