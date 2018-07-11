// Copyright 2018 Aleksandr Demakin. All rights reserved.

package main

import (
	"flag"

	"github.com/avdva/slot-machine/machine"
	"github.com/avdva/slot-machine/server"
)

func main() {
	flagAddr := flag.String("addr", ":1313", "serve addr")
	flag.Parse()
	config := server.Config{Addr: *flagAddr,
		Machines: map[string]machine.Interface{
			"atkins-diet": nil,
		}}
	server.New(config).Serve()
}
