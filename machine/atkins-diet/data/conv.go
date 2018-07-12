// Copyright 2018 Aleksandr Demakin. All rights reserved.

package main

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"strconv"

	atkins "github.com/avdva/slot-machine/machine/atkins-diet"
)

func main() {
	config := atkins.Config{
		Wild:           10,
		Scatter:        11,
		BonusFreeSpins: 10,
		BonusBetMult:   3,
		Pays:           make(map[int]atkins.Line),
	}
	payLines, err := os.Open("paylines.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer payLines.Close()
	r := csv.NewReader(payLines)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	for _, record := range records {
		var l atkins.Line
		if len(record) != 5 {
			log.Fatal("bad paylines")
		}
		for i, r := range record {
			if l[i], err = strconv.Atoi(r); err != nil {
				log.Fatal(err)
			}
		}
		config.Paylines = append(config.Paylines, l)
	}

	payTable, err := os.Open("payTable.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer payTable.Close()
	r = csv.NewReader(payTable)
	records, err = r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	for _, record := range records {
		if len(record) != 6 {
			log.Fatal("bad paytable")
		}
		if symb, err := strconv.Atoi(record[0]); err == nil {
			var l atkins.Line
			for i := 1; i < len(record); i++ {
				if l[i-1], err = strconv.Atoi(record[i]); err != nil {
					log.Fatal(err)
				}
			}
			config.Pays[symb] = l
		} else {
			log.Fatal(err)
		}
	}

	reels, err := os.Open("reels.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer reels.Close()
	r = csv.NewReader(reels)
	records, err = r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for _, record := range records {
		var l atkins.Line
		if len(record) != 5 {
			log.Fatal("bad reels")
		}
		for i, r := range record {
			if l[i], err = strconv.Atoi(r); err != nil {
				log.Fatal(err)
			}
		}
		config.Reels = append(config.Reels, l)
	}

	out, err := os.Create("atkins.json")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	enc := json.NewEncoder(out)
	if err := enc.Encode(config); err != nil {
		log.Fatal(err)
	}
}
