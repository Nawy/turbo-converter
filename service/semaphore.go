package service

type empty struct{}
type Semaphore chan empty

func (s Semaphore) Acquire(n int) {
	e := empty{}

	for i := 0; i < n; i++ {
		s <- e
	}
}

func (s Semaphore) Release(n int) {
	for i := 0; i < n; i++ {
		<-s
	}
}

func (s Semaphore) Lock() {
	s.Acquire(1)
}

func (s Semaphore) Unlock() {
	s.Release(1)
}

func (s Semaphore) Signal() {
	s.Release(1)
}

func (s Semaphore) Wait(n int) {
	s.Acquire(n)
}
