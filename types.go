package pinky

// Promise ...
type Promise struct {
	wasCaught  bool
	isResolved bool
	isRejected bool

	value interface{}
	err   error

	chn chan PromiseResult
}

// PromiseResult returns the value and or error of a promise
type PromiseResult struct {
	Value interface{}
	Err   error
}

// ResolveFunc resolve function
type ResolveFunc func(res interface{})

// RejectFunc reject function
type RejectFunc func(err error)
