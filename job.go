package nbjobqueue

type Job interface {
	Run()
}

type JobFunc func()

func (f JobFunc) Run() {}
