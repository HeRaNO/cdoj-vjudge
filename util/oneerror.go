package util

import "sync"

type OneError struct {
	errOnce sync.Once
	Err     error
}

func (e *OneError) Add(err error) {
	e.errOnce.Do(func() {
		e.Err = err
	})
}
