package main

type operationType string

const (
	set       operationType = "set"
	get       operationType = "get"
	deleteMap operationType = "deleteMap"
	lenMap    operationType = "lenMap"
)

type Request[K comparable, V any] struct {
	Operation operationType
	Key       K
	Value     V
	Result    chan any
	Succeed   bool
}

type SafeChannelMap[K comparable, V any] struct {
	m  map[K]V
	ch chan *Request[K, V]
}

func NewSafeChannelMap[K comparable, V any]() *SafeChannelMap[K, V] {
	safeMap := &SafeChannelMap[K, V]{
		m:  make(map[K]V),
		ch: make(chan *Request[K, V]),
	}
	go safeMap.Run()
	return safeMap
}

func (safeMap *SafeChannelMap[K, V]) Run() {
	for req := range safeMap.ch {
		switch req.Operation {
		case set:
			safeMap.m[req.Key] = req.Value
			req.Succeed = true
			req.Result <- nil
		case get:
			value, ok := safeMap.m[req.Key]
			req.Succeed = ok
			req.Result <- value
		case deleteMap:
			_, ok := safeMap.m[req.Key]
			delete(safeMap.m, req.Key)
			req.Succeed = ok
			req.Result <- nil
		case lenMap:
			req.Succeed = true
			req.Result <- len(safeMap.m)

		}
	}
}

func (safeMap *SafeChannelMap[K, V]) Set(k K, v V) {
	result := make(chan any)
	safeMap.ch <- &Request[K, V]{
		Operation: set,
		Key:       k,
		Value:     v,
		Result:    result,
	}
	<-result
}

func (safeMap *SafeChannelMap[K, V]) Get(k K) (V, bool) {
	result := make(chan any)
	req := &Request[K, V]{
		Operation: get,
		Key:       k,
		Result:    result,
	}
	safeMap.ch <- req
	res := <-result
	if v, ok := res.(V); ok {
		return v, req.Succeed
	}
	panic("type assertion failed")

}

func (safeMap *SafeChannelMap[K, V]) Delete(k K) {
	result := make(chan any)
	safeMap.ch <- &Request[K, V]{
		Operation: deleteMap,
		Key:       k,
		Result:    result,
	}
	<-result
}

func (safeMap *SafeChannelMap[K, V]) Len() int {
	result := make(chan any)
	safeMap.ch <- &Request[K, V]{
		Operation: lenMap,
		Result:    result,
	}
	l := <-result
	return l.(int)
}
