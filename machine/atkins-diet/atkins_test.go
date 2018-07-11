// Copyright 2018 Aleksandr Demakin. All rights reserved.

package atkins

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
