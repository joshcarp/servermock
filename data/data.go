package data

import "sync"

func StoreData(sm *sync.Map, trace string, data []byte) {
	val, ok := sm.Load(trace)
	if !ok {
		sm.Store(trace, [][]byte{data})
		return
	}
	valBytes := val.([][]byte)
	if len(valBytes) == 0 {
		sm.Store(trace, [][]byte{data})
		return
	}
	sm.Store(trace, append([][]byte{data}, valBytes...)) //[n-1] is always the element to be read (and deleted) first
}

func LoadData(sm *sync.Map, trace string) []byte {
	val, ok := sm.Load(trace)
	if !ok {
		return nil
	}
	valBytes := val.([][]byte)
	if len(valBytes) == 0 {
		return nil
	}
	sm.Store(trace, valBytes[:len(valBytes)-1])
	return valBytes[len(valBytes)-1]
}
