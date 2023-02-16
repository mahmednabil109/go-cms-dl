package utils

// Job encapsulate a Job function
type Job func() error

// Pool is a simple gorountine pool
type Pool struct {
	submit chan<- Job
	busy   chan<- bool
	done   chan bool
}

// New creates a new Job Pool
func NewPool(count int) *Pool {
	var (
		submit = make(chan Job)
		done   = make(chan bool)
		busy   = make(chan bool, count)
	)
	pool := &Pool{submit, busy, done}
	go run(submit, done, busy)
	return pool
}

// Submit adds a job to a pool
func (p *Pool) Submit(job Job) {
	go func() {
		select {
		case p.busy <- true:
			p.submit <- job
		case <-p.done:
			return
		}
	}()
}

// Close finishes the current job and ends the pool
func (p *Pool) Close() {
	go func() { p.done <- true }()
	close(p.submit)
	close(p.busy)
}

func run(submit <-chan Job, done <-chan bool, busy <-chan bool) {
	for {
		select {
		case <-done:
			return
		case job := <-submit:
			go func() {
				// TODO(mahmednabil109): need better error handling
				_ = (func() error)(job)()
				<-busy
			}()
		}
	}
}
