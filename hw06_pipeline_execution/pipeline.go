package hw06pipelineexecution

import (
	"sort"
	"sync"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type resWithIndex struct {
	index int
	value interface{}
}

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(Bi)

	var mu sync.Mutex
	var unsortedResults []resWithIndex
	var wg sync.WaitGroup

	sortIndex := 0
loop:
	for {
		select {
		case <-done:
			break loop
		case v, ok := <-in:
			if !ok {
				break loop
			}
			wg.Add(1)

			// создаем отдельную горутину для каждого значения из in, и записываем результат в слайс
			go func(sortIndex int, value interface{}) {
				defer wg.Done()
				vChannel := channelWrap(done, value)
				for _, stage := range stages {
					if stage != nil {
						select {
						case <-done:
							return
						default:
							vChannel = stage(wrapWithDone(done, vChannel))
						}
					}
				}
				resValue, ok := <-vChannel
				if !ok {
					return
				}
				mu.Lock()
				unsortedResults = append(unsortedResults, resWithIndex{sortIndex, resValue})
				mu.Unlock()
			}(sortIndex, v)
			sortIndex++
		}
	}

	// дожидаемся обработки канала, сортируем и отдаем.
	go func() {
		wg.Wait()
		sort.Slice(unsortedResults, func(i, j int) bool {
			return unsortedResults[i].index < unsortedResults[j].index
		})

		for _, sortedRes := range unsortedResults {
			out <- sortedRes.value
		}
		close(out)
	}()

	return out
}

func channelWrap(done In, value interface{}) Out {
	stream := make(Bi)
	go func() {
		defer close(stream)
		select {
		case <-done:
			return
		case stream <- value:
		}
	}()
	return stream
}

func wrapWithDone(done In, in In) Out {
	out := make(Bi)
	go func() {
		defer func() {
			close(out)
			for range in {
			}
		}()
		select {
		case <-done:
			return
		default:
		}

		for {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				out <- v
			}
		}
	}()
	return out
}
