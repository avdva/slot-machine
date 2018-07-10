// Copyright 2018 Aleksandr Demakin. All rights reserved.

package atkins

import (
	"math/rand"
	"sync"
	"time"

	"github.com/avdva/slot-machine/machine"
)

type Line [5]int

type Config struct {
	Wild     int
	Scatter  int
	Paylines []Line
	Reels    []Line
}

type Machine struct {
	config Config

	m sync.Mutex
	// use or own mutex-protected Rand object instead of the global one from math/rand.
	r *rand.Rand
}

func New(config Config) *Machine {
	return &Machine{
		config: config,
		r:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (m *Machine) Spin(bet float64) []machine.SpinResult {
	var workingReels [3]Line
	lines := m.roll()
	for i, l := range lines {
		workingReels[0][i] = m.config.Reels[(l-1+32)%32][i]
		workingReels[1][i] = m.config.Reels[l][i]
		workingReels[2][i] = m.config.Reels[(l+1)%32][i]
	}
	return nil
}

func (m *Machine) roll() (l Line) {
	m.m.Lock()
	for i := 0; i < 5; i++ {
		l[i] = m.r.Intn(32)
	}
	m.m.Unlock()
	return
}

func (m *Machine) checkPayLines(lines [3]Line, paylines []Line) {
	for _, l := range paylines {
		m.checkPayLine(lines, l)
	}
}

func (m *Machine) checkPayLine(lines [3]Line, payline Line) {
	prev := lines[payline[0]][0]
	count, wildCount := 1, 1
	if prev != m.config.Wild {
		wildCount = 0
	}
	for i := 1; i < 5; i++ {
		cur := lines[payline[i]][i]
		if prev == m.config.Wild || prev == cur {
			if cur == m.config.Wild && wildCount > 0 {
				wildCount++
			} else {
				wildCount = 0
			}
			count++
		} else {
			break
		}
	}
	if count < 3 {
		return
	}

}
