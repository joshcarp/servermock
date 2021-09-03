package dmt

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
)

func storeData(sm *sync.Map, trace string, data Request) {
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

func loadData(sm *sync.Map, key string) (Request, bool) {
	val, ok := sm.Load(key)
	if !ok {
		return Request{}, false
	}
	valBytes := val.([]Request)
	if len(valBytes) == 0 {
		return Request{}, false
	}
	req := valBytes[len(valBytes)-1]
	if req.IsStack {
		sm.Store(key, valBytes[:len(valBytes)-1])
	}
	return req, true
}

func loadAllData(sm *sync.Map, key string) ([]Request, bool) {
	val, ok := sm.Load(key)
	if !ok {
		return nil, false
	}
	valBytes := val.([]Request)
	if len(valBytes) == 0 {
		return nil, false
	}
	return valBytes, true
}

func resetData(sm *sync.Map) {
	sm.Range(func(key, _ interface{}) bool {
		sm.Delete(key)
		return true
	})
}

func setData(wr http.ResponseWriter, r *http.Request, log Logger, sm *sync.Map, Endpoint string) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log("Error setting Data for request: %s\n", Endpoint)
	}
	entry := Request{}
	if err = json.Unmarshal(b, &entry); err != nil {
		wr.WriteHeader(500)
		return
	}
	storeData(sm, Endpoint, entry)
	log("Setting Data for request: %s Length: %d\n", Endpoint, len(entry.Body))
}
