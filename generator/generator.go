package generator

type Generator interface {
	Start() error
	Stop()
}
