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
	addr   = flag.String("addr", ":8080", "http service address")
	flagP  = flag.Bool("p", false, "dev: proxy to \"http://localhost:3000/\"")
	logdir = flag.String("l", "logs/", "log directory")
)

func main() {
	currentState.Status = StatusIdle
	flag.Parse()
	err := OpenLog(*logdir)
	defer CloseLog()
	if err != nil {
		log.Fatal(err)
	}
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

	// Start the flashrom monitoring job
	go JobMonitor()

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}
