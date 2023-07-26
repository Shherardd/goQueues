package MQWP

type Job interface {
	Do()
}

type _Job struct {
	job       Job
	queue     string
	queueLock chan struct{}
}

type _FinishSignal struct{}

type _EmptyJob struct {
	finishSignReceiver chan _FinishSignal
}

func (quitJob *_EmptyJob) Do() {
	quitJob.finishSignReceiver <- _FinishSignal{}
}
