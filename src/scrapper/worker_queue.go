package scrapper

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

// ScrapResult is the success case with the metric value
type ScrapResult struct {
	Val               interface{}
	ConstMetricIdxOut int
	JobName           string
	InstanceName      string
	DataSource        string
}

// ScrapErrResult is the fail case with the error happened
type ScrapErrResult struct {
	Err               error
	ConstMetricIdxOut int
	JobName           string
	InstanceName      string
	DataSource        string
}

// ScrapRequest have the scrapper to do an scrap, the channels to put the result, and the metric index to return
type ScrapRequest struct {
	Scrap            Scrapper
	Res              chan ScrapResult
	ConstMetricIdxIn int
	JobName          string
	InstanceName     string
	DataSource       string
	Err              chan ScrapErrResult
	MetricsCollector chan<- prometheus.Metric
}

type scrapWork struct {
	scrapper         Scrapper
	res              chan ScrapResult
	constMetricIdxIn int
	jobName          string
	instanceName     string
	dataSource       string
	err              chan ScrapErrResult
	metricsCollector chan<- prometheus.Metric
}

type workQueue chan scrapWork
type workerQueue chan workQueue

// Pool of workers to run scrap works
type Pool struct {
	workers  workerQueue
	works    workQueue
	nWorkers uint
	wg       sync.WaitGroup
}

// NewPool create a pool of workers to run scrap tasks
func NewPool(workersNum uint) *Pool {
	return &Pool{
		workers:  make(chan workQueue, workersNum),
		works:    make(chan scrapWork, workersNum*2),
		nWorkers: workersNum,
		wg:       sync.WaitGroup{},
	}
}

// Wait for all goroutines running inside the wait group
func (p *Pool) Wait() {
	p.wg.Wait()
}

// Apply push a scrapper task to the workers pool to be executed
func (p *Pool) Apply(ri ScrapRequest) {
	work := scrapWork{
		scrapper:         ri.Scrap,
		res:              ri.Res,
		constMetricIdxIn: ri.ConstMetricIdxIn,
		jobName:          ri.JobName,
		instanceName:     ri.InstanceName,
		dataSource:       ri.DataSource,
		err:              ri.Err,
		metricsCollector: ri.MetricsCollector,
	}
	p.works <- work
}

type workerT struct {
	works    chan scrapWork
	workers  chan workQueue
	quitChan chan bool
	wg       *sync.WaitGroup
}

func (p *Pool) newWorker() workerT {
	return workerT{
		works:    make(chan scrapWork),
		workers:  p.workers,
		quitChan: make(chan bool),
		wg:       &p.wg,
	}
}

func (w *workerT) start() {
	w.wg.Add(1)
	go func() {
		for {
			// Add ourselves as available into the worker queue.
			w.workers <- w.works
			select {
			// Wait for a work request.
			case work := <-w.works:
				val, err := work.scrapper.GetMetric(work.metricsCollector)
				if err == nil {
					work.res <- ScrapResult{Val: val, ConstMetricIdxOut: work.constMetricIdxIn, JobName: work.jobName, InstanceName: work.instanceName, DataSource: work.dataSource}
				} else {
					work.err <- ScrapErrResult{Err: err, ConstMetricIdxOut: work.constMetricIdxIn, JobName: work.jobName, InstanceName: work.instanceName, DataSource: work.dataSource}
				}
				// wait for a quit msg
			case <-w.quitChan:
				// Receive a close worker request.
				w.wg.Done()
				return
			}
		}
	}()
}

// StartDispatcher make the workers pool ready to run scraps
func (p *Pool) StartDispatcher() {
	for i := uint(0); i < p.nWorkers; i++ {
		worker := p.newWorker()
		worker.start()
	}
	go func() {
		// wait for an incoming work
		for work := range p.works {
			// wait for an available worker
			worker := <-p.workers
			// dispatch work into worker
			worker <- work
		}
	}()
}
