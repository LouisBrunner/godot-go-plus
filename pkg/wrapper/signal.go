package wrapper

type Signal interface {
	Emit(args ...any)
}

// TODO: need to intercept the creation to add our wrapper around the signal
