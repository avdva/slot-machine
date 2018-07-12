// Copyright 2018 Aleksandr Demakin. All rights reserved.

package machine

// SpinResult is what we get after a spin.
type SpinResult struct {
	Type  string
	Total int64
	Stops []int
}

// Result contains spins and total pay.
type Result struct {
	Spins []SpinResult
	Total int64
}

// Interface must be satifsied by any slot-machine.
type Interface interface {
	Wager(bet int64) int64
	Spin(bet int64) Result
}
