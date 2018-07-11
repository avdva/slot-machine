// Copyright 2018 Aleksandr Demakin. All rights reserved.

package atkins

import (
	"math/rand"
	"sync"
	"time"

	"github.com/avdva/slot-machine/machine"
)

const (
	invalidSymbol = -1
)

type Line [5]int

type Config struct {
	Wild           int
	Scatter        int
	Paylines       []Line
	Reels          []Line
	Pays           map[int]Line
	BonusFreeSpins int
	BonusBetMult   int
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

func (m *Machine) Spin(bet int64) []machine.SpinResult {
	var result []machine.SpinResult
	sr, free := m.doMain(bet)
	result = append(result, sr)
	if free {
		result = append(result, m.doFree(bet)...)
	}
	return result
}

func (m *Machine) doSpin() (machine.SpinResult, bool) {
	var workingReels [3]Line
	lines := m.roll()
	for i, l := range lines {
		workingReels[0][i] = m.config.Reels[(l-1+32)%32][i]
		workingReels[1][i] = m.config.Reels[l][i]
		workingReels[2][i] = m.config.Reels[(l+1)%32][i]
	}
	strikes := m.checkPayLines(workingReels)
	pay := m.calcPay(strikes)
	bonusPay, freeSpins := m.checkBonus(workingReels, strikes)
	return machine.SpinResult{
		Total: int64(pay + bonusPay),
		Stops: lines[:],
	}, freeSpins
}

func (m *Machine) doMain(bet int64) (machine.SpinResult, bool) {
	sr, free := m.doSpin()
	sr.Type = "main"
	sr.Total *= bet
	return sr, free
}

func (m *Machine) doFree(bet int64) []machine.SpinResult {
	var result []machine.SpinResult
	mults := []int{m.config.BonusBetMult}
	for len(mults) > 0 {
		mult := mults[0]
		for i := 0; i < m.config.BonusFreeSpins; i++ {
			sr, free := m.doSpin()
			sr.Type = "free"
			sr.Total *= (int64(mult) * bet)
			result = append(result, sr)
			if free {
				mults = append(mults, mult*m.config.BonusBetMult)
			}
		}
		mults = mults[1:]
	}
	return result
}

func (m *Machine) roll() (l Line) {
	m.m.Lock()
	for i := 0; i < 5; i++ {
		l[i] = m.r.Intn(32)
	}
	m.m.Unlock()
	return
}

func (m *Machine) calcPay(strikes []strike) int {
	var result int
	for _, s := range strikes {
		result += m.strikePay(s)
	}
	return result
}

func (m *Machine) checkPayLines(lines [3]Line) []strike {
	var result []strike
	for _, l := range m.config.Paylines {
		if strikes := checkPayLine(lines, l, m.config.Wild); len(strikes) > 0 {
			maxIdx, max := 0, m.strikePay(strikes[0])
			for i := 1; i < len(strikes); i++ {
				if cur := m.strikePay(strikes[i]); cur > max {
					max = cur
					maxIdx = i
				}
			}
			result = append(result, strikes[maxIdx])
		}
	}
	return result
}

func (m *Machine) checkBonus(lines [3]Line, strikes []strike) (pay int, haveBonus bool) {
	var count int
	// scatters can be anywhere in the lines.
	for _, l := range lines {
		for _, symb := range l {
			if symb == m.config.Scatter {
				count++
			}
		}
	}
	pay = m.strikePay(strike{n: count, symb: m.config.Scatter})
	// we could've given it already as a part of a line bonus.
	// if so, pay should be 0.
	if haveBonus = pay > 0; haveBonus {
		for _, strike := range strikes {
			if strike.symb == m.config.Scatter && strike.nsym == count { // a line with exactly 'count' scatters.
				pay = 0
				break
			}
		}
	}
	return
}

func (m *Machine) strikePay(s strike) int {
	l, found := m.config.Pays[s.symb]
	if !found || s.n == 0 {
		return 0
	}
	return l[s.n-1]
}

// strike used to check if we should pay for a roll.
type strike struct {
	n    int // strike len.
	nsym int // how many 'symb's were in the strike, excluding wilds.
	symb int // what symbol makes a strike.
}

func checkPayLine(lines [3]Line, payline Line, wild int) []strike {
	var totalCount, symbCount, wildCount int
	symb := invalidSymbol
	for i, pi := range payline {
		cur := lines[pi][i]
		if cur == wild {
			if symb == -1 {
				wildCount++
			} else {
				totalCount++
			}
		} else if symb == invalidSymbol {
			totalCount = i + 1
			wildCount = i
			symbCount++
			symb = cur
		} else if symb == cur {
			symbCount++
			totalCount++
		} else {
			break
		}
	}
	var result []strike
	if totalCount > 0 {
		result = append(result, strike{n: totalCount, nsym: symbCount, symb: symb})
	}
	if wildCount > 0 {
		result = append(result, strike{n: wildCount, nsym: wildCount, symb: wild})
	}
	return result
}
