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

type Position uint8

const (
	PositionDealer Position = iota
	PositionSmallBlind
	PositionBigBlind
)

type GameConfig struct {
	SmallBlind int64
	BigBlind   int64
}

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
	Config GameConfig

	State State

	Players []Player

	CommunityCards []Card

	Pot int64

	Deck *Deck

	CurrentPlayer int
	CurrentBet    int64

	DealerPosition int
}

func NewGame(config GameConfig) *Game {
	return &Game{
		Config:         config,
		State:          StateWaiting,
		Players:        make([]Player, 0, 9),
		CommunityCards: make([]Card, 0, 5),
		DealerPosition: 0,
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

	g.CurrentPlayer = 0
	g.CurrentBet = 0

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

	smallBlind, bigBlind := g.blindPositions()

	g.postBlind(smallBlind, g.Config.SmallBlind)

	g.postBlind(bigBlind, g.Config.BigBlind)

	g.CurrentBet = g.Players[bigBlind].Bet

	if len(g.Players) == 2 {
		g.CurrentPlayer = smallBlind
	} else {
		g.CurrentPlayer = nextPosition(bigBlind, 1, len(g.Players))
	}

	g.State = StatePreFlop

	return nil
}

func (g *Game) AdvanceRound() error {
	switch g.State {
	case StatePreFlop:
		if err := g.dealFlop(); err != nil {
			return err
		}

		g.resetBettingRound()
		return nil

	case StateFlop:
		if err := g.dealTurn(); err != nil {
			return err
		}

		g.resetBettingRound()
		return nil

	case StateTurn:
		if err := g.dealRiver(); err != nil {
			return err
		}

		g.resetBettingRound()
		return nil

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

func (g *Game) CurrentPlayerID() string {
	if g.CurrentPlayer < 0 || g.CurrentPlayer >= len(g.Players) {
		return ""
	}

	return g.Players[g.CurrentPlayer].ID
}

func (g *Game) advancePlayer() error {
	if len(g.Players) == 0 {
		return fmt.Errorf("no players")
	}

	for i := 1; i <= len(g.Players); i++ {
		index := (g.CurrentPlayer + i) % len(g.Players)

		player := &g.Players[index]

		if !player.Folded && !player.AllIn {
			g.CurrentPlayer = index
			return nil
		}
	}

	return fmt.Errorf("no active players")
}

func (g *Game) resetBettingRound() {
	for i := range g.Players {
		g.Players[i].Bet = 0
	}

	g.CurrentBet = 0
}

func (g *Game) blindPositions() (int, int) {
	playerCount := len(g.Players)

	if playerCount == 2 {
		// Dealer is small blind.
		return g.DealerPosition, nextPosition(g.DealerPosition, 1, playerCount)
	}

	smallBlind := nextPosition(g.DealerPosition, 1, playerCount)

	bigBlind := nextPosition(g.DealerPosition, 2, playerCount)

	return smallBlind, bigBlind
}

func nextPosition(position, offset, playerCount int) int {
	return (position + offset) % playerCount
}

func (g *Game) postBlind(playerIndex int, amount int64) int64 {
	player := &g.Players[playerIndex]

	if amount >= player.Chips {
		amount = player.Chips
		player.Chips = 0
		player.Bet += amount
		player.AllIn = true
	} else {
		player.Chips -= amount
		player.Bet += amount
	}

	g.Pot += amount

	return amount
}
