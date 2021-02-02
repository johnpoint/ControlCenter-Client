// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package apis

import (
	"log"
	"os/exec"
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

	c := websocket.Conn{}
	defer c.Close()
	for {
		log.Println("[Websocket] Retry...")
		c = websocket.Conn{}
		log.Printf("[Websocket] connecting to %s", url)
		c, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			log.Println("[Websocket]", err)
			time.Sleep(time.Duration(5) * time.Second)
			continue
		}

		err = c.WriteMessage(websocket.TextMessage, []byte(token))
		if err != nil {
			log.Println("[Websocket] write:", err)
			continue
		}

		go func() {
			for {
				_, message, err := c.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					break
				}
				log.Printf("recv: %s", message)
			}
		}()

		status := make(chan string)
		go func() {
			for {
				status <- "pushStatus#" + infoMiniJSON()
				time.Sleep(time.Duration(3) * time.Second)
			}
		}()

		go func() {
			for {
				status <- "psStatus#" + runCommand("ps axc -o pid,user,stat,pcpu,pmem,command --sort -pcpu --no-header | head -n 100 | sed 's/\\ \\+/\\ /g'")
				time.Sleep(time.Duration(30) * time.Second)
			}
		}()

		for {
			t := <-status
			err := c.WriteMessage(websocket.TextMessage, []byte(t))
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	}
}

func runCommand(command string) string {
	cmd := exec.Command("bash", "-c", command)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	return string(out)
}
