package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"slices"
	"strconv"

	"github.com/gorilla/websocket"
)

//go:embed static/*
var content embed.FS

// const WORDS = "citrus,sir,sit,its,cut,suit,cuts,stir,tis,crust,rust,rut,curt,rustic,citrus"
var g Game

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]string)
var broadcast = make(chan Msg)

func handleConnect(w http.ResponseWriter, r *http.Request) {
	fmt.Println("upgrading")
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade problem:", err)
	}

	defer c.Close()

	clients[c] = fmt.Sprintf("user%v", strconv.FormatInt(int64(len(clients)+1), 10))
	for _, msg := range g.Process(Cmd{Kind: Join, From: clients[c], Data: clients[c]}) {
		broadcast <- msg
	}
	for _, msg := range g.Process(Cmd{Kind: Start, From: clients[c]}) {
		broadcast <- msg
	}
	for {
		var msg Cmd
		err := c.ReadJSON(&msg)
		if err != nil {
			fmt.Println(err)
			delete(clients, c)
			return
		}

		if msg.Kind == "reset" {
			for client := range clients {
				client.Close()
				delete(clients, client)
			}
			g = NewGame()
			return
		}

		for _, msg := range g.Process(msg) {
			broadcast <- msg
		}
	}

}

func handleMessage() {
	for {
		msg := <-broadcast
		fmt.Println("Sending", msg)
		for client := range clients {
			if !slices.Contains(msg.To, clients[client]) {
				continue
			}
			err := client.WriteJSON(msg)
			if err != nil {
				fmt.Println(err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	g = NewGame()

	if _, ok := os.LookupEnv("REPL_SLUG"); ok {
		fs, _ := fs.Sub(content, "static")
		http.Handle("/", http.FileServer(http.FS(fs)))
	} else {
		http.Handle("/", http.FileServer(http.Dir("static")))
	}
	http.HandleFunc("/connect", handleConnect)

	go handleMessage()

	fmt.Println("Starting server")
	http.ListenAndServe(":"+"8080", nil)
}
