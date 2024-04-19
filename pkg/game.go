package main

import (
	"fmt"
	"slices"
	"strings"

	"golang.org/x/exp/maps"
)

const (
	// Game states
	Lobby   = "lobby"
	Playing = "playing"
	Ended   = "ended"
)

type Game struct {
	players map[string][]string
	words   []string
	base    string
	state   string
}

const (
	// Inputs
	Join    = "join"
	Start   = "start"
	TryWord = "tryWord"

	// Outputs
	Err       = "err"
	Joined    = "joined"
	Setup     = "setup"
	FoundWord = "foundWord"
)

type Cmd struct {
	Kind string      `json:"kind"`
	From string      `json:"from"`
	Data interface{} `json:"data"`
}

type Msg struct {
	Kind string      `json:"kind"`
	To   []string    `json:"to"`
	Data interface{} `json:"data"`
}

const WORDS = "citrus,sir,sit,its,cut,suit,cuts,stir,tis,crust,rust,rut,curt,rustic,citrus"

func NewGame() Game {
	words := strings.Split(WORDS, ",")
	return Game{
		players: make(map[string][]string),
		words:   words[1:],
		base:    words[0],
		state:   Lobby,
	}
}

func (g Game) process(cmd Cmd) Msg {
	everyone := maps.Keys(g.players)
	sender := []string{cmd.From}
	switch cmd.Kind {
	case Join:
		name, _ := cmd.Data.(string)
		if len(g.players) == 2 {
			return Msg{Err, []string{name}, "Game already full"}
		}
		g.players[name] = make([]string, 0)
		return Msg{Joined, everyone, maps.Keys(g.players)}
	case Start:
		// if len(g.players) != 2 {
		// 	return Msg{Err, sender, "Not enough players"}
		// }
		g.state = Playing
		wordCounts := make([]int, len(g.words))
		for i := range g.words {
			wordCounts[i] = len(g.words[i])
		}
		return Msg{Setup, everyone, wordCounts}
	case TryWord:
		word := cmd.Data.(string)
		if !slices.Contains(g.words, word) {
			return Msg{Err, sender, "Invalid word"}
		}
		if slices.Contains(g.players[cmd.From], word) {
			return Msg{Err, sender, "Already found word"}
		}
		return Msg{}
	default:
		return Msg{Err, sender, fmt.Sprintf("Unknown cmd %s", cmd.Kind)}
	}
}