package trigger

type Trigger interface {
	Suspend(reason string) error
	Hibernate(reason string) error
	SystemPreferred(fallback string) string
}

func New() Trigger { return newImpl() }
