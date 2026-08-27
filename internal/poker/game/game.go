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
	Bet    int64
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
		g.Players[i].Bet = 0
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

func (g *Game) AdvanceRound() error {
	switch g.State {
	case StatePreFlop:
		return g.dealFlop()

	case StateFlop:
		return g.dealTurn()

	case StateTurn:
		return g.dealRiver()

	case StateRiver:
		g.State = StateShowdown
		return nil

	default:
		return fmt.Errorf("cannot advance game from state %q", g.State)
	}
}

func (g *Game) dealFlop() error {
	if g.Deck == nil {
		return fmt.Errorf("deck is not initialized")
	}

	if !g.Deck.Burn() {
		return fmt.Errorf("failed to burn card before flop")
	}

	for i := 0; i < 3; i++ {
		card, ok := g.Deck.Draw()
		if !ok {
			return fmt.Errorf("failed to deal flop")
		}

		g.CommunityCards = append(g.CommunityCards, card)
	}

	g.State = StateFlop

	return nil
}

func (g *Game) dealTurn() error {
	if g.Deck == nil {
		return fmt.Errorf("deck is not initialized")
	}

	if !g.Deck.Burn() {
		return fmt.Errorf("failed to burn card before flop")
	}

	card, ok := g.Deck.Draw()
	if !ok {
		return fmt.Errorf("failed to deal turn")
	}

	g.CommunityCards = append(g.CommunityCards, card)
	g.State = StateTurn

	return nil
}

func (g *Game) dealRiver() error {
	if g.Deck == nil {
		return fmt.Errorf("deck is not initialized")
	}

	if !g.Deck.Burn() {
		return fmt.Errorf("failed to burn card before flop")
	}

	card, ok := g.Deck.Draw()
	if !ok {
		return fmt.Errorf("failed to deal river")
	}

	g.CommunityCards = append(g.CommunityCards, card)
	g.State = StateRiver

	return nil
}
