package main

import (
	"fmt"
	"math/rand"
	"slices"
	"sort"
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

	foundWordBuffer []Msg
}

const (
	// Inputs
	Join    = "join"    // data: name
	Start   = "start"   // data: nothing
	TryWord = "tryWord" // data: word to try

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

type FoundWordInfo struct {
	Word   string `json:"word"`
	Player string `json:"player"`
	Pos    int    `json:"pos"`
}

type SetupInfo struct {
	Letters string `json:"letters"`
	Counts  []int  `json:"counts"`
}

const WORDS = "citrus,sir,sit,its,cut,suit,cuts,stir,tis,crust,rust,rut,curt,rustic,citrus"

func NewGame() Game {
	words := strings.Split(WORDS, ",")
	base := strings.Split(words[0], "")
	rand.Shuffle(len(base), func(i, j int) {
		base[i], base[j] = base[j], base[i]
	})
	words = words[1:]
	sort.Slice(words, func(i, j int) bool {
		if len(words[i]) == len(words[j]) {
			return words[i] < words[j]
		}
		return len(words[i]) < len(words[j])
	})
	return Game{
		players: make(map[string][]string),
		words:   words,
		base:    strings.Join(base, ""),
		state:   Lobby,
	}
}

func (g *Game) getSetupInfo() SetupInfo {
	wordCounts := make([]int, len(g.words))
	for i := range g.words {
		wordCounts[i] = len(g.words[i])
	}
	return SetupInfo{Letters: g.base, Counts: wordCounts}
}

func (g *Game) Process(cmd Cmd) []Msg {
	everyone := maps.Keys(g.players)
	sender := []string{cmd.From}
	switch cmd.Kind {
	case Join:
		name, _ := cmd.Data.(string)
		// if len(g.players) == 2 {
		// 	return Msg{Err, []string{name}, "Game already full"}
		// }
		g.players[name] = make([]string, 0)
		ret := []Msg{{Joined, maps.Keys(g.players), name}}
		if g.state == Playing {
			ret = append(ret, Msg{Setup, sender, g.getSetupInfo()})
			ret = append(ret, g.foundWordBuffer...)
		}
		return ret
	case Start:
		// if len(g.players) != 2 {
		// 	return Msg{Err, sender, "Not enough players"}
		// }
		g.state = Playing
		return []Msg{{Setup, everyone, g.getSetupInfo()}}
	case TryWord:
		if g.state != Playing {
			return []Msg{{Err, sender, "Need to start"}}
		}
		word := cmd.Data.(string)
		if !slices.Contains(g.words, word) {
			return []Msg{{Err, sender, fmt.Sprintf("Invalid word: %s", word)}}
		}
		idx := slices.Index(g.words, word)
		if idx == -1 {
			return []Msg{{Err, sender, "Already found word"}}
		}
		// update player words
		g.players[cmd.From] = append(g.players[cmd.From], word)
		ret := []Msg{{
			FoundWord,
			everyone,
			FoundWordInfo{
				word,
				cmd.From,
				idx,
			},
		}}
		g.foundWordBuffer = append(g.foundWordBuffer, ret[0])
		return ret
	default:
		return []Msg{{Err, sender, fmt.Sprintf("Unknown cmd %s", cmd.Kind)}}
	}
}
