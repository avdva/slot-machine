// Copyright 2018 Aleksandr Demakin. All rights reserved.

package server

type response struct {
	Total float64 `json:"total"`
	Spins []spin  `json:"spins"`
}

type spin struct {
	Type  string  `json:"type"`
	Total float64 `json:"total"`
	Stops []int   `json:"stops"`
}
