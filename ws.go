// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var (
	verbDebug = fmt.Printf
	//verbDebug = func(string, ...interface{}) {}
	logDebug = fmt.Printf
	//logDebug = func(string, ...interface{}) {}
)

// wsCmd describe the command send in the jsom message over ws
type WsCmd string

// WsCmd possible values
const (
	WsCmdSet    WsCmd = "set"
	WsCmdUpdate WsCmd = "update"
	WsCmdStart  WsCmd = "start"
)

type WsMessage struct {
	Cmd   WsCmd
	Field string      `json:",omitempty"`
	Value interface{} `json:",omitempty"`
}

var (
	clients = make(map[*websocket.Conn]chan State)
)

func register(ws *websocket.Conn) chan State {
	logDebug("register %v\n", ws.RemoteAddr())
	ch := make(chan State)
	go func() {
		ch <- GetState()
	}()
	clients[ws] = ch
	return ch
}

func unregister(ws *websocket.Conn) {
	logDebug("unregister %v\n", ws.RemoteAddr())
	delete(clients, ws)
}

func SendCurrentState() {
	for _, ch := range clients {
		ch <- GetState()
	}
}

const (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

var upgrader = websocket.Upgrader{}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logDebug("%v\n", err)
		return
	}
	logDebug("Client subscribed %v\n", ws.RemoteAddr())
	go writer(ws)
	reader(ws)
	logDebug("Client unsubscribed %v\n", ws.RemoteAddr())
}

func reader(ws *websocket.Conn) {
	defer ws.Close()
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		messageType, p, err := ws.ReadMessage()
		if err != nil {
			verbDebug("ReadMessage: %v\n", err)
			break
		}
		if messageType == websocket.TextMessage {
			verbDebug("Recived %q\n", p)
			var msg WsMessage
			err := json.Unmarshal(p, &msg)
			if err != nil {
				verbDebug("ReadMessage Unmarshal: %v\n", err)
			}
			logDebug("Recived %v\n", msg)
			switch msg.Cmd {
			case WsCmdStart:
				confId, ok := msg.getConfigId()
				if !ok || confId == 0 {
					UpdateConfigId(0)
					SetErrorState("Please select a configuration!")
					break
				}
				UpdateConfigId(confId)
				StartJob()
			default:
				verbDebug("ReadMessage Unsupported Cmd: %v\n", msg.Cmd)
			}
		}
	}
}

func writer(ws *websocket.Conn) {
	pingTicker := time.NewTicker(pingPeriod)
	defer func() {
		pingTicker.Stop()
		ws.Close()
	}()
	//ws.WriteMessage(websocket.TextMessage, []byte("Connected !"))
	ch := register(ws)
	defer unregister(ws)
	//var lastState State
	for {
		select {
		case newState := <-ch:
			msg := WsMessage{Cmd: WsCmdSet, Value: newState}
			b, err := json.Marshal(&msg)
			if err != nil {
				verbDebug("MarshalMessage: %v\n", err)
				return
			}

			verbDebug("sending %q\n", b)
			err = ws.WriteMessage(websocket.TextMessage, b)
			if err != nil {
				verbDebug("WriteMessage: %v\n", err)
				return
			}

		case <-pingTicker.C:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				verbDebug("Write PingMessage: %v\n", err)
				return
			}
		}
	}
}

func (msg WsMessage) getConfigId() (int, bool) {
	//fmt.Println(reflect.TypeOf(msg.Value).String())
	value, ok := msg.Value.(map[string]interface{})
	if !ok {
		verbDebug("ReadMessage start: missing value map\n")
		return 0, false
	}
	confId, ok := value["ConfigId"]
	if !ok {
		verbDebug("ReadMessage start: missing ConfigId value\n")
		return 0, false
	}
	//fmt.Println(reflect.TypeOf(confId).String())
	val, ok := confId.(float64)
	if !ok {
		verbDebug("ReadMessage start: ConfigId is not aa number\n")
		return 0, false
	}
	return int(val), ok
}
