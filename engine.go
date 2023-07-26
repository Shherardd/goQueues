package MQWP

import (
	"sync"
)

type MqwpEngine interface {
	Dispatch(job Job, queueName ...string)
	Stop()
}

type mqwp struct {
	sync.RWMutex
	wg          *sync.WaitGroup
	mqListener  chan _Job
	workerPools map[string]chan *_Worker
}

func NewDispatcher(numWorkers int) (MqwpEngine, error) {
	if numWorkers <= 0 {
		return nil, newError("Invalid number of workers given to create a new mqwp")
	}

	wg := &sync.WaitGroup{}
	mqwp := &mqwp{wg: wg}
	mqwp.mqListener = make(chan _Job)
	mqwp.workerPools = make(map[string]chan *_Worker)

	for i := 0; i < numWorkers; i++ {
		mqwp.workerPools["default"] = make(chan *_Worker, numWorkers)
		worker := &_Worker{wg: wg}
		mqwp.workerPools["default"] <- worker
	}

	go func() {
		for delayedJob := range mqwp.mqListener {
			queueName := delayedJob.queue
			if queueName == "" {
				queueName = "default"
			}

			workerPool, ok := mqwp.workerPools[queueName]
			if !ok {
				workerPool = make(chan *_Worker, numWorkers)
				mqwp.workerPools[queueName] = workerPool

				for i := 0; i < numWorkers; i++ {
					worker := &_Worker{wg: wg}
					workerPool <- worker
				}
			}

			worker := <-workerPool
			worker.wg.Add(1)
			go func(job Job, worker *_Worker, wpool chan *_Worker) {
				job.Do()
				wpool <- worker
				worker.wg.Done()
			}(delayedJob.job, worker, workerPool)
		}
	}()

	return mqwp, nil
}

func (dispatcher *mqwp) Dispatch(j Job, queueName ...string) {
	dispatcher.RLock()
	defer dispatcher.RUnlock()

	if len(queueName) == 0 {
		queueName = append(queueName, "default")
	}

	if len(queueName) > 1 {
		panic("Too many arguments given to Dispatch")
	}

	dispatcher.mqListener <- _Job{job: j, queue: queueName[0]}
}

func (dispatcher *mqwp) Stop() {
	dispatcher.Lock()
	defer dispatcher.Unlock()

	finishSignReceiver := make(chan _FinishSignal)
	dispatcher.mqListener <- _Job{
		job: &_EmptyJob{
			finishSignReceiver: finishSignReceiver},
		queue: "stop"}

	<-finishSignReceiver

	dispatcher.wg.Wait()
	close(dispatcher.mqListener)
}
