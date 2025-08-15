package inhibit

type Inhibitor interface {
	Acquire(reason string) error
	Release() error
}

func New() Inhibitor { return newImpl() }
