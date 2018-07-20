// Copyright 2018 Aleksandr Demakin. All rights reserved.

package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"

	"github.com/avdva/slot-machine/machine"
	"github.com/avdva/slot-machine/machine/atkins-diet"
	"github.com/avdva/slot-machine/server"
)

func main() {
	flagAddr := flag.String("addr", ":1313", "serve addr")
	flagMachine := flag.String("machines", "atkins-diet", "comma-separated list of machines")
	flag.Parse()
	parts := strings.Split(*flagMachine, ",")
	config := server.Config{
		Addr:     *flagAddr,
		Machines: make(map[string]machine.Interface),
	}
	for _, name := range parts {
		config.Machines[name] = makeMachine(name)
	}
	if len(config.Machines) == 0 {
		log.Fatal("no machines")
	}
	s := server.New(config)
	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, os.Interrupt)
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	go func() {
		<-sigchan
		s.Stop()
	}()
	if err := s.Serve(); err != nil {
		log.Fatal(err)
	}
}

func makeMachine(machine string) machine.Interface {
	if machine != "atkins-diet" {
		log.Fatal("bad machine name " + machine)
	}
	var config atkins.Config
	f, err := os.Open("machine/atkins-diet/data/atkins.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&config); err != nil {
		log.Fatal(err)
	}
	m, err := atkins.New(config, atkins.NewRandLineSource())
	if err != nil {
		log.Fatal(err)
	}
	return m
}
