package idle

type Reader interface{ UserIdleSeconds() (uint64, error) }

func New() Reader { return newImpl() }
