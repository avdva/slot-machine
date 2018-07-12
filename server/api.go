// Copyright 2018 Aleksandr Demakin. All rights reserved.

package server

type request struct {
	Machine string
	UID     string
	Chips   int64
	Bet     int64
}

type response struct {
	Total int64  `json:"total"`
	Spins []spin `json:"spins"`
	JWT   string `json:"jwt"`
}

type spin struct {
	Type  string `json:"type"`
	Total int64  `json:"total"`
	Stops []int  `json:"stops"`
}
