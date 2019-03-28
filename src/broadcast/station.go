package broadcast

import (
	"config"
	"fmt"
	"runtime"
)

var (
	maxWorker = runtime.NumCPU()
	maxQuene  = 1000
)

type _Action func()

type _Job struct {
	Action _Action
}

var _JobQueue chan _Job

type _Worker struct {
	WorkerPool chan chan _Job
	JobChannel chan _Job
	Quit       chan bool
}

func newWorker(workerPool chan chan _Job) _Worker {
	return _Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan _Job),
		Quit:       make(chan bool),
	}
}

func (w _Worker) start() {
	go func() {
		for {
			w.WorkerPool <- w.JobChannel
			select {
			case job := <-w.JobChannel:
				job.Action()
				fmt.Println("job doing..")

			case <-w.Quit:
				return
			}
		}
	}()
}

func (w _Worker) stop() {
	go func() {
		w.Quit <- true
	}()
}

type _Dispatcher struct {
	MaxWorkers int
	WorkerPool chan chan _Job
	Quit       chan bool
}

func newDispatcher(maxWorkers int) *_Dispatcher {
	pool := make(chan chan _Job, maxWorkers)
	return &_Dispatcher{MaxWorkers: maxWorkers, WorkerPool: pool, Quit: make(chan bool)}
}

func (d *_Dispatcher) run() {
	for i := 0; i < d.MaxWorkers; i++ {
		worker := newWorker(d.WorkerPool)
		worker.start()
	}

	go d.dispatch()

}

func (d *_Dispatcher) stop() {
	go func() {
		d.Quit <- true
	}()
}

func (d *_Dispatcher) dispatch() {
	for {
		select {
		case job := <-_JobQueue:
			go func(job _Job) {
				jobChannel := <-d.WorkerPool
				jobChannel <- job
			}(job)
		case <-d.Quit:
			return
		}
	}
}

func RunStation() {

	err := config.InitConfig()
	if err != nil {
		fmt.Println(err)
	}

	runtime.GOMAXPROCS(maxWorker)
	_JobQueue = make(chan _Job, maxQuene)
	dispatcher := newDispatcher(maxWorker)
	dispatcher.run()
}
