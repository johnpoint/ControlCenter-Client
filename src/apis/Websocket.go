// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package apis

import (
	"log"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

func websocketPush(url string, token string) {
	if strings.Contains(url, "https://") {
		url = "wss://" + strings.Split(url, "//")[1] + "/server/v2"
	} else {
		url = "ws://" + strings.Split(url, "//")[1] + "/server/v2"
	}
	for {
		log.Printf("connecting to %s", url)
		c, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			log.Println("[Websocket]", err)
			time.Sleep(time.Duration(10) * time.Second)
			continue
		}
		defer c.Close()
		err = c.WriteMessage(websocket.TextMessage, []byte(token))
		if err != nil {
			log.Println("[Websocket] write:", err)
			time.Sleep(time.Duration(10) * time.Second)
			continue
		}

		done := make(chan struct{})

		go func() {
			defer close(done)
			for {
				_, message, err := c.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					return
				}
				log.Printf("recv: %s", message)
			}
		}()

		status := make(chan string)
		go func() {
			defer close(done)
			for {
				status <- "pushStatus#" + infoMiniJSON()
				time.Sleep(time.Duration(3) * time.Second)
			}
		}()

		for {
			select {
			case <-done:
				return
			case t := <-status:
				err := c.WriteMessage(websocket.TextMessage, []byte(t))
				if err != nil {
					log.Println("write:", err)
					return
				}
			}
		}
	}
}
