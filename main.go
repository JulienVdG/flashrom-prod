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
)

var (
	addr  = flag.String("addr", ":8080", "http service address")
	flagP = flag.Bool("p", false, "dev: proxy to \"http://localhost:3000/\"")
)

func main() {
	flag.Parse()
	if *flagP {
		rpURL, err := url.Parse("http://localhost:3000/")
		if err != nil {
			log.Fatal(err)
		}
		http.Handle("/", httputil.NewSingleHostReverseProxy(rpURL))
	} // TODO: rice or file
	http.HandleFunc("/ws", WsHandler)

	// Testing
	currentState.Status = StatusIdle
	currentState.Config = []string{"One", "Two", "Three"}
	//currentState.ConfigId = 2

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}
