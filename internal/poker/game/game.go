package game

import (
	"fmt"
	"math/rand/v2"
)

type State string

const (
	StateWaiting  State = "waiting"
	StatePreFlop  State = "pre_flop"
	StateFlop     State = "flop"
	StateTurn     State = "turn"
	StateRiver    State = "river"
	StateShowdown State = "showdown"
)

type Player struct {
	ID     string
	Seat   int
	Chips  int64
	Cards  []Card
	Folded bool
	AllIn  bool
}

type Game struct {
	State State

	Players []Player

	CommunityCards []Card

	Pot int64

	Deck *Deck
}

func NewGame() *Game {
	return &Game{
		State:          StateWaiting,
		Players:        make([]Player, 0, 9),
		CommunityCards: make([]Card, 0, 5),
		Pot:            0,
	}
}

func (g *Game) AddPlayer(player Player) error {
	if g.State != StateWaiting {
		return fmt.Errorf("game is already running")
	}

	if len(g.Players) >= 9 {
		return fmt.Errorf("table is full")
	}

	for _, existing := range g.Players {
		if existing.ID == player.ID {
			return fmt.Errorf("player is already at the table")
		}

		if existing.Seat == player.Seat {
			return fmt.Errorf("seat is already occupied")
		}
	}

	if player.Chips <= 0 {
		return fmt.Errorf("player must have chips")
	}

	g.Players = append(g.Players, player)

	return nil
}

func (g *Game) StartHand(r *rand.Rand) error {
	if g.State != StateWaiting {
		return fmt.Errorf("game is not waiting")
	}

	if len(g.Players) < 2 {
		return fmt.Errorf("at least two players are required")
	}

	g.Deck = NewDeck()
	g.Deck.Shuffle(r)

	g.CommunityCards = g.CommunityCards[:0]
	g.Pot = 0

	for i := range g.Players {
		g.Players[i].Cards = nil
		g.Players[i].Folded = false
		g.Players[i].AllIn = false
	}

	for i := range g.Players {
		card1, ok := g.Deck.Draw()
		if !ok {
			return fmt.Errorf("failed to deal first card")
		}

		card2, ok := g.Deck.Draw()
		if !ok {
			return fmt.Errorf("failed to deal second card")
		}

		g.Players[i].Cards = []Card{
			card1,
			card2,
		}
	}

	g.State = StatePreFlop

	return nil
}
