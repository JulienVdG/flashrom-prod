// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	rice "github.com/GeertJohan/go.rice"
)

var (
	addr  = flag.String("addr", ":8080", "http service address")
	flagP = flag.Bool("p", false, "dev: proxy to \"http://localhost:3000/\"")
)

func main() {
	currentState.Status = StatusIdle
	flag.Parse()
	ReadConfig()
	if *flagP {
		rpURL, err := url.Parse("http://localhost:3000/")
		if err != nil {
			log.Fatal(err)
		}
		http.Handle("/", httputil.NewSingleHostReverseProxy(rpURL))
	} else {
		box, err := rice.FindBox("dist")
		if err != nil {
			log.Fatal(err)
		}
		http.Handle("/", http.FileServer(box.HTTPBox()))

	}
	http.HandleFunc("/ws", WsHandler)

	// Testing
	AddLog(StatusIdle, "Starting", "This is to display raw command for more info,\n\nText is:\n - preformated,\n - multiline.\n")

	// Start the flashrom monitoring job
	go JobMonitor()

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}
