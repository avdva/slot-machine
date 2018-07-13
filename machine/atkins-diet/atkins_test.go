// Copyright 2018 Aleksandr Demakin. All rights reserved.

package atkins

import (
	"testing"

	"github.com/avdva/slot-machine/machine"
	"github.com/stretchr/testify/assert"
)

type testLineSource struct {
	data Line
}

func (tls *testLineSource) Line() Line {
	return tls.data
}

func TestCheckPayLine(t *testing.T) {
	const wild = 13
	a := assert.New(t)
	pl := Line{0, 0, 0, 0, 0}
	lines := [3]Line{
		Line{1, 1, 1, 1, 1},
	}
	actual := checkPayLine(lines, pl, wild)
	expected := []strike{strike{n: 5, nsym: 5, symb: 1}}
	a.Equal(expected, actual)

	lines[0] = Line{1, 2, 2, 2, 2}
	actual = checkPayLine(lines, pl, wild)
	expected = []strike{strike{n: 1, nsym: 1, symb: 1}}
	a.Equal(expected, actual)

	lines[0] = Line{1, 1, 2, 2, 2}
	actual = checkPayLine(lines, pl, wild)
	expected = []strike{strike{n: 2, nsym: 2, symb: 1}}
	a.Equal(expected, actual)

	lines[0] = Line{1, 1, wild, 2, 2}
	actual = checkPayLine(lines, pl, wild)
	expected = []strike{strike{n: 3, nsym: 2, symb: 1}}
	a.Equal(expected, actual)

	lines[0] = Line{1, 1, wild, 1, 2}
	actual = checkPayLine(lines, pl, wild)
	expected = []strike{strike{n: 4, nsym: 3, symb: 1}}
	a.Equal(expected, actual)

	lines[0] = Line{wild, 1, wild, 1, 2}
	actual = checkPayLine(lines, pl, wild)
	expected = []strike{strike{n: 4, nsym: 2, symb: 1}, strike{n: 1, nsym: 1, symb: wild}}
	a.Equal(expected, actual)

	lines[0] = Line{wild, wild, wild, 1, 2}
	actual = checkPayLine(lines, pl, wild)
	expected = []strike{strike{n: 4, nsym: 1, symb: 1}, strike{n: 3, nsym: 3, symb: wild}}
	a.Equal(expected, actual)

	lines[0] = Line{wild, 1, 2, wild, 2}
	actual = checkPayLine(lines, pl, wild)
	expected = []strike{strike{n: 2, nsym: 1, symb: 1}, strike{n: 1, nsym: 1, symb: wild}}
	a.Equal(expected, actual)

	lines[0] = Line{wild, wild, wild, wild, wild}
	actual = checkPayLine(lines, pl, wild)
	expected = []strike{strike{n: 5, nsym: 5, symb: wild}}
	a.Equal(expected, actual)
}

func makeTestConfig() Config {
	config := Config{
		Wild:    10,
		Scatter: 11,
		Paylines: []Line{
			Line{1, 1, 1, 1, 1},
		},
		BonusBetMult:   3,
		BonusFreeSpins: 10,
		Pays:           make(map[int]Line),
	}
	for i := 0; i < 32; i++ {
		val := (i % 11) + 1
		l := Line{val, val, val, val, val}
		config.Reels = append(config.Reels, l)
	}
	for i := 1; i <= 11; i++ {
		l := Line{0, 0, 1, 2, 3}
		config.Pays[i] = l
	}
	return config
}

func TestDoMain(t *testing.T) {
	a := assert.New(t)
	config := makeTestConfig()
	tls := &testLineSource{data: Line{0, 0, 0, 0, 0}}
	m, _ := New(config, tls)
	res, free := m.doMain(1)
	a.Equal(machine.SpinResult{Type: "main", Total: 3, Stops: tls.data[:]}, res)
	a.False(free)

	tls.data = Line{8, 8, 8, 8, 8}
	_, free = m.doMain(2)
	a.False(free)

	tls.data = Line{9, 9, 9, 9, 9}
	_, free = m.doMain(2)
	a.True(free)

	tls.data = Line{10, 10, 10, 10, 10}
	res, free = m.doMain(2)
	a.Equal(machine.SpinResult{Type: "main", Total: 6, Stops: tls.data[:]}, res)
	a.True(free)

	tls.data = Line{11, 11, 11, 11, 11}
	_, free = m.doMain(2)
	a.True(free)

	tls.data = Line{12, 12, 12, 12, 12}
	_, free = m.doMain(2)
	a.False(free)

	tls.data = Line{3, 3, 3, 9, 10}
	res, free = m.doMain(2)
	a.Equal(machine.SpinResult{Type: "main", Total: 4, Stops: tls.data[:]}, res)
	a.False(free)

	tls.data = Line{9, 10, 11, 9, 10}
	res, free = m.doMain(2)
	a.Equal(machine.SpinResult{Type: "main", Total: 6, Stops: tls.data[:]}, res)
	a.True(free)
}

type testCountedLineSource struct {
	count int
	f     func(int) Line
}

func (ts *testCountedLineSource) Line() Line {
	result := ts.f(ts.count)
	ts.count++
	return result
}

func TestDoFree(t *testing.T) {
	a := assert.New(t)
	config := makeTestConfig()
	tls := &testCountedLineSource{
		f: func(count int) Line {
			if count == 5 {
				return Line{10, 10, 10, 10, 10}
			}
			return Line{1, 1, 1, 1, 1}
		},
	}
	m, _ := New(config, tls)

	totalFunc := func(res []machine.SpinResult) (sum int64) {
		for _, r := range res {
			sum += r.Total
		}
		return
	}

	res := m.doFree(1)
	a.Equal(20, len(res))
	a.Equal(int64(360), totalFunc(res))

	tls.count = 0
	tls.f = func(count int) Line {
		if count == 5 || count == 7 {
			return Line{11, 11, 11, 11, 11}
		}
		return Line{1, 1, 1, 1, 1}
	}
	res = m.doFree(1)
	a.Equal(30, len(res))
	a.Equal(int64(648), totalFunc(res))

	tls.count = 0
	tls.f = func(count int) Line {
		if count == 5 {
			return Line{11, 11, 11, 11, 11}
		}
		return Line{1, 1, 1, 1, 1}
	}
	res = m.doFree(1)
	a.Equal(20, len(res))
	a.Equal(int64(369), totalFunc(res))

	tls.count = 0
	tls.f = func(count int) Line {
		if count == 5 || count == 15 {
			return Line{10, 10, 10, 10, 10}
		}
		return Line{1, 1, 1, 1, 1}
	}
	res = m.doFree(1)
	a.Equal(30, len(res))
	a.Equal(int64(1170), totalFunc(res))

	tls.count = 0
	tls.f = func(count int) Line {
		if count == 5 || count == 15 || count == 25 {
			return Line{10, 10, 10, 10, 10}
		}
		return Line{1, 1, 1, 1, 1}
	}
	res = m.doFree(1)
	a.Equal(40, len(res))
	a.Equal(int64(3600), totalFunc(res))
}
