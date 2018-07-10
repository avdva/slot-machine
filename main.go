// Copyright 2018 Aleksandr Demakin. All rights reserved.

package main

import (
	"flag"
)

func main() {
	flagAddr = flag.String("addr", ":1313", "serve addr")
	flag.Parse()
}
