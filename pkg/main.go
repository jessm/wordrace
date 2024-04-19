package main

import (
	"fmt"
	"net/http"
	"slices"
	"sort"
	"strings"
	"text/template"

	"github.com/gorilla/websocket"
)

// const WORDS = "citrus,sir,sit,its,cut,suit,cuts,stir,tis,crust,rust,rut,curt,rustic,citrus"

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn][]string)
var broadcast = make(chan Message)

type Message struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, "")
}

func handleConnect(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade problem:", err)
	}

	defer c.Close()

	clients[c] = make([]string, 0)

	words := strings.Split(WORDS, ",")
	letters := strings.Split(words[0], "")
	sort.Slice(letters, func(i, j int) bool {
		return letters[i] < letters[j]
	})

	for client := range clients {
		err := client.WriteJSON(Message{
			Name:    "server",
			Message: strings.Join(letters, ","),
		})
		if err != nil {
			fmt.Println(err)
			client.Close()
			delete(clients, client)
			return
		}
	}

	for {
		var msg Message
		err := c.ReadJSON(&msg)
		if err != nil {
			fmt.Println(err)
			delete(clients, c)
			return
		}

		fmt.Println(msg)

		if slices.Contains(words, msg.Message) && !slices.Contains(clients[c], msg.Message) {
			// if slices.Contains(words, msg.Message) {
			broadcast <- msg
			clients[c] = append(clients[c], msg.Message)
		}

		if len(words)-1 == len(clients[c]) {
			broadcast <- Message{"server", "You Win!"}
		}

	}
}

func handleMessage() {
	for {
		msg := <-broadcast
		fmt.Println("Sending", msg)
		for client := range clients {
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
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/connect", handleConnect)

	go handleMessage()

	fmt.Println("Starting server")
	http.ListenAndServe(":"+"8080", nil)
}
