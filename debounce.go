// Copyright © 2019 Bjørn Erik Pedersen <bjorn.erik.pedersen@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package debounce provides a debouncer func. The most typical use case would be
// the user typing a text into a form; the UI needs an update, but let's wait for
// a break.
package debounce

import (
	"sync"
	"time"
)

// New returns a debounced function that takes another functions as its argument.
// This function will be called when the debounced function stops being called
// for the given duration, or has been called the given reps number of times.
// Pass a negative integer to consider the duration only.
// The debounced function can be invoked with different functions, if needed,
// the last one will win.
func New(after time.Duration, reps uint64) func(f func()) {
	d := &debouncer{after: after, reps: reps}

	return func(f func()) {
		d.add(f)
	}
}

type debouncer struct {
	mu    sync.Mutex
	after time.Duration
	timer *time.Timer
	reps  uint64
	count uint64
}

func (d *debouncer) add(f func()) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
	}

	if d.count == d.reps {
		f()
		d.count = 0
	} else {
		d.count++
		d.timer = time.AfterFunc(d.after, f)
	}
}
