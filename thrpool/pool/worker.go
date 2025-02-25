package pool

type Worker struct {
	jobChannel  chan interface{}
	workerPool  chan chan interface{}
	closeHandle chan bool
}

//------------------------------------------------------------------------------
// NewWorker : Create the new worker
//------------------------------------------------------------------------------
func NewWorker(workerPool chan chan interface{}, closeHandle chan bool) *Worker {
	return &Worker{workerPool: workerPool, jobChannel: make(chan interface{}), closeHandle: closeHandle}
}

//------------------------------------------------------------------------------
// NewWorker : Create the new worker
//------------------------------------------------------------------------------
func (w Worker) Start() {
	go func() {
		for {
			// Put the worker to the worker threadpool
			w.workerPool <- w.jobChannel

			select {
			// Wait for the job
			case job := <-w.jobChannel:
				w.ExecuteJob(job)
			case <-w.closeHandle:
				return
			}
		}
	}()
}

//------------------------------------------------------------------------------
// NewWorker : Create the new worker
//------------------------------------------------------------------------------
func (w Worker) ExecuteJob(job interface{}) {
	// Execute the job based on the task type
	switch task := job.(type) {
	case Runnable:
		task.Run()
		break
	default:
		break
	}
}
