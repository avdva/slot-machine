// Copyright 2018 Aleksandr Demakin. All rights reserved.

package machine

// SpinResult is what we get after a spin.
type SpinResult struct {
	Type  string
	Total float64
	Stops []int
}

// Machine must be satifsied by any slot-machine.
type Interface interface {
	Spin(bet float64) []SpinResult
}
