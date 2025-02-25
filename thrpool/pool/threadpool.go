package pool

import "fmt"

var (
	ErrQueueFull = fmt.Errorf("queue is full, not able add the task")
)

type ThreadPool struct {
	queueSize   int64
	noOfWorkers int

	jobQueue    chan interface{}
	workerPool  chan chan interface{}
	closeHandle chan bool // Channel used to stop all the workers

	//-------------------------------------------------------------
	jobCnt         int
	jobCloseHandle chan chan bool
}

//------------------------------------------------------------------------------
// NewThreadPool : create new threadpool
//------------------------------------------------------------------------------
func NewThreadPool(noOfWorkers int, queueSize int64) *ThreadPool {
	threadpool := &ThreadPool{queueSize: queueSize, noOfWorkers: noOfWorkers}
	threadpool.jobQueue = make(chan interface{}, queueSize)
	threadpool.workerPool = make(chan chan interface{}, noOfWorkers)
	threadpool.closeHandle = make(chan bool)
	threadpool.createPool()

	return threadpool
}

//------------------------------------------------------------------------------
// submitTask : Add the task to the jobQueue
//------------------------------------------------------------------------------
func (t *ThreadPool) submitTask(task interface{}) error {
	// Add the task to the job queue
	if len(t.jobQueue) == int(t.queueSize) {
		return ErrQueueFull
	}

	// set closeHandle
	
	//t.jobCloseHandle <- a.CloseHandle
	// jonCnt increase
	t.jobCnt++

	t.jobQueue <- task
	return nil
}

//------------------------------------------------------------------------------
// Execute : Submit the job to available worker
//------------------------------------------------------------------------------
func (t *ThreadPool) Execute(task Runnable) error {
	return t.submitTask(task)
}

//------------------------------------------------------------------------------
// Close : Close the threadpool
// TODO : need to check the existing /running task before closeing the theadpool
//------------------------------------------------------------------------------
func (t *ThreadPool) Close() {
	close(t.closeHandle) // stop all the routine
	close(t.workerPool)  // close the job threadpool
	close(t.jobQueue)    // close the job queue
}

//------------------------------------------------------------------------------
// createPool : creates the workers and start listening on the jobQueue
//------------------------------------------------------------------------------
func (t *ThreadPool) createPool() {
	for i := 0; i < t.noOfWorkers; i++ {
		worker := NewWorker(t.workerPool, t.closeHandle)
		worker.Start()
	}

	go t.dispatch()
}

//------------------------------------------------------------------------------
// createPool : creates the workers and start listening on the jobQueue
//------------------------------------------------------------------------------
func (t *ThreadPool) dispatch() {
	for {
		select {
		case job := <-t.jobQueue:
			// Get job
			func(job interface{}) {
				// Find a worker for the job
				jobChannel := <-t.workerPool
				// submit job to the worker
				jobChannel <- job
			}(job)
		case <-t.closeHandle:
			return
		}
	}
}
