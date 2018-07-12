// Copyright 2018 Aleksandr Demakin. All rights reserved.

package atkins

import (
	"math/rand"
	"sync"
	"time"
)

// RandLineSource prodoces random lines.
type RandLineSource struct {
	m sync.Mutex
	// use or own mutex-protected Rand object instead of the global one from math/rand.
	r *rand.Rand
}

// NewRandLineSource returns new RandLineSource.
func NewRandLineSource() *RandLineSource {
	return &RandLineSource{r: rand.New(rand.NewSource(time.Now().UnixNano()))}
}

// Line returns new line randomized by math/rand.
func (r *RandLineSource) Line() (l Line) {
	r.m.Lock()
	for i := 0; i < 5; i++ {
		l[i] = r.r.Intn(32)
	}
	r.m.Unlock()
	return
}
