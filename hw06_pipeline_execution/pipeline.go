package hw06pipelineexecution

import "sync"

type keyValues struct {
	mu  sync.Mutex
	kvs []keyValue
}

type keyValue struct {
	key   interface{}
	value interface{}
}

func newKeyValues() *keyValues {
	return &keyValues{
		kvs: make([]keyValue, 0),
	}
}

func (kvs *keyValues) add(kv keyValue) {
	kvs.mu.Lock()
	defer kvs.mu.Unlock()
	kvs.kvs = append(kvs.kvs, kv)
}

func (kvs *keyValues) len() int {
	kvs.mu.Lock()
	defer kvs.mu.Unlock()
	return len(kvs.kvs)
}

func (kvs *keyValues) set(key, value interface{}) {
	kvs.mu.Lock()
	defer kvs.mu.Unlock()
	for i, kv := range kvs.kvs {
		if kv.key == key {
			kvs.kvs[i].value = value
			break
		}
	}
}

func (kvs *keyValues) get() []keyValue {
	kvs.mu.Lock()
	defer kvs.mu.Unlock()
	return kvs.kvs
}

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func (s Stage) exec(input interface{}) Out {
	in := make(Bi, 1)
	in <- input
	return s(in)
}

func executor(job interface{}, done In, stages ...Stage) interface{} {
	for _, s := range stages {
		select {
		case <-done:
			return nil
		case job = <-s.exec(job):
		}
	}
	return job
}

func ExecutePipeline(in, done In, stages ...Stage) Out {
	var out Bi
	wg := sync.WaitGroup{}

	// used keyvalues struct instead map to avoid ordering
	kvs := newKeyValues()

	for v := range in {
		// keep order of job
		kvs.add(keyValue{key: v})
		wg.Add(1)
		go func(v interface{}, done In) {
			defer wg.Done()
			result := executor(v, done, stages...)
			// saved the result by key(job)
			kvs.set(v, result)
		}(v, done)
	}

	wg.Wait()

	// if have result, sending ...
	// if kvs.len() > 0 {
	out = make(Bi, kvs.len())

	for _, kv := range kvs.get() {
		// send only non-null result
		if kv.value != nil {
			out <- kv.value
		}
	}
	//}
	// closed the channel after sending, to avoid waitings
	close(out)
	return out
}
