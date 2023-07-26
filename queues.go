package MQWP

type JobQueueCollection struct {
	queues map[string][]Job
}

func (qc *JobQueueCollection) Enqueue(queueName string, job Job) {
	if qc.queues == nil {
		qc.queues = make(map[string][]Job)
	}

	if _, exists := qc.queues[queueName]; !exists {
		qc.queues[queueName] = make([]Job, 0)
	}

	qc.queues[queueName] = append(qc.queues[queueName], job)
}

func (qc *JobQueueCollection) Dequeue(queueName string) Job {
	if queue, exists := qc.queues[queueName]; exists {
		if len(queue) > 0 {
			job := queue[0]
			qc.queues[queueName] = queue[1:]

			return job
		}
	}
	return nil
}
