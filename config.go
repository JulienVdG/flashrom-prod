// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Configs []struct {
		Name        string
		Commandline string
		FlashChip   string
	}
}

var Cfg Config

// ReadConfig reads config.json, store it in Cfg and update currentState
// Warning this must be done early as currentState is not protected here
func ReadConfig() {
	f, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}

	decoder := json.NewDecoder(f)
	err = decoder.Decode(&Cfg)
	if err != nil {
		log.Fatal(err)
	}

	currentState.Config = []string{}
	for _, c := range Cfg.Configs {
		currentState.Config = append(currentState.Config, c.Name)
	}
}
