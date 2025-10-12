package interfaces

type Job interface {
	Name() string
	Schedule() string
	Run()
}
