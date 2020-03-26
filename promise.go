package pinky

import (
	"errors"
	"time"
)

// NewPromise Creates a new promise with a starting value
func NewPromise(value interface{}) *Promise {
	return &Promise{
		wasCaught:  false,
		isResolved: false,
		isRejected: false,

		value: value,
		err:   nil,

		chn: nil,
	}
}

// GetChan returns a channel with PromiseResult. The result is only sent when `.Finally` is called.
func (p *Promise) GetChan() chan PromiseResult {
	if p.chn == nil {
		p.chn = make(chan PromiseResult, 1)
	}

	return p.chn
}

// Resolve stores value in promise
func (p *Promise) Resolve(res interface{}) {
	p.isResolved = true
	p.value = res
}

// Reject rejects promise and stores error
func (p *Promise) Reject(err error) {
	p.isRejected = true
	p.err = err
}

// Then allows for synchronous execution of functions chained to promise
//	// Param f has multi types types:
//	// func(value interface{}, resolve ResolveFunc, reject RejectFunc)
//	// func(resolve ResolveFunc, reject RejectFunc)
//	// func(value interface{}) (interface{}, error)
//	// func() (interface{}, error)
//	// func(value interface{}) error
//	// func() error
//	// func(value interface{})
//	// func()
func (p *Promise) Then(f interface{}) *Promise {
	if p.isRejected {
		return p
	}

	switch function := f.(type) {
	case func(value interface{}, resolve ResolveFunc, reject RejectFunc):
		function(p.value, p.Resolve, p.Reject)
	case func(resolve ResolveFunc, reject RejectFunc):
		function(p.Resolve, p.Reject)
	case func(interface{}) (interface{}, error):
		res, err := function(p.value)
		if err != nil {
			p.Reject(err)
		} else {
			p.Resolve(res)
		}
	case func() (interface{}, error):
		res, err := function()
		if err != nil {
			p.Reject(err)
		} else {
			p.Resolve(res)
		}
	case func(interface{}) error:
		err := function(p.value)
		if err != nil {
			p.Reject(err)
		}
	case func() error:
		err := function()
		if err != nil {
			p.Reject(err)
		}
	case func(interface{}):
		function(p.value)
	case func():
		function()
	default:
		panic(errors.New("bad promise Then function type"))
	}

	return p
}

// CatchCase catches error from promise chain matches it to target error
func (p *Promise) CatchCase(targetErr error, f func(err error)) *Promise {
	if p.wasCaught || targetErr == nil {
		return p
	} else if errors.Is(p.err, targetErr) {
		f(p.err)
	}

	p.wasCaught = true
	return p
}

// Catch catches any error from promise chain
func (p *Promise) Catch(f func(err error)) *Promise {
	return p.CatchDefault(f)
}

// CatchDefault catches any error from promise chain that wasnt caught from CatchCase
func (p *Promise) CatchDefault(f func(err error)) *Promise {
	if p.wasCaught {
		return p
	} else if p.err != nil {
		f(p.err)
	}

	p.wasCaught = true
	return p
}

// Tap allows for synchronous execution of functions chained to promise without changing promise value
func (p *Promise) Tap(f func(interface{})) *Promise {
	f(p.value)

	return p
}

// Finally allows the promise to finish. This chain can be used to signal go channels for async execution.
func (p *Promise) Finally(f func(interface{}, error)) (interface{}, error) {
	f(p.value, p.err)

	if p.chn != nil {
		p.chn <- PromiseResult{
			Value: p.value,
			Err:   p.err,
		}
	}

	return p.value, p.err
}

// Delay stops promise chain
func (p *Promise) Delay(d time.Duration) *Promise {
	time.Sleep(d)

	return p
}
