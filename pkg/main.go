package main

import (
	"fmt"
  "slices"
	"net/http"
	"text/template"

	"github.com/gorilla/websocket"
)

// const WORDS = "citrus,sir,sit,its,cut,suit,cuts,stir,tis,crust,rust,rut,curt,rustic,citrus"
var g Game

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]string)
var broadcast = make(chan Msg)


func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, "")
}

func handleConnect(w http.ResponseWriter, r *http.Request) {
  fmt.Println("upgrading")
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade problem:", err)
	}

	defer c.Close()

  for {
    var msg Cmd
    err := c.ReadJSON(&msg)
    if err != nil {
      fmt.Println(err)
      delete(clients, c)
      return
    }

    broadcast <-g.Process(msg)
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
  
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/connect", handleConnect)

	go handleMessage()

	fmt.Println("Starting server")
	http.ListenAndServe(":"+"8080", nil)
}
