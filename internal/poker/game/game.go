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

type Game struct {
	Config GameConfig

	State State

	Players []Player

	CommunityCards []Card

	Pot int64

	Deck *Deck

	CurrentPlayer int
	CurrentBet    int64

	// Minimum amount by which the current bet can be increased.
	MinRaise int64

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
	g.MinRaise = g.Config.BigBlind

	g.Deck = NewDeck()
	g.Deck.Shuffle(r)

	g.CommunityCards = g.CommunityCards[:0]
	g.Pot = 0

	for i := range g.Players {
		g.Players[i].Cards = nil
		g.Players[i].Folded = false
		g.Players[i].AllIn = false
		g.Players[i].Bet = 0
		g.Players[i].Acted = false
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
	g.MinRaise = g.Config.BigBlind

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

		return g.setFirstPostFlopPlayer()

	case StateFlop:
		if err := g.dealTurn(); err != nil {
			return err
		}

		g.resetBettingRound()

		return g.setFirstPostFlopPlayer()

	case StateTurn:
		if err := g.dealRiver(); err != nil {
			return err
		}

		g.resetBettingRound()

		return g.setFirstPostFlopPlayer()

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

func (g *Game) FinishHand() error {
	if g.State != StateShowdown {
		return fmt.Errorf("hand is not at showdown")
	}

	g.rotateDealer()

	g.State = StateWaiting

	return nil
}

func (g *Game) EvaluateShowdown() (ShowdownResult, error) {
	if g.State != StateShowdown {
		return ShowdownResult{}, fmt.Errorf("game is not at showdown")
	}

	if len(g.CommunityCards) != 5 {
		return ShowdownResult{}, fmt.Errorf("expected 5 community cards, got %d", len(g.CommunityCards))
	}

	result := ShowdownResult{
		Winners: make([]int, 0),
		Hands:   make(map[int]HandValue),
	}

	var best HandValue
	first := true

	for i := range g.Players {
		player := &g.Players[i]

		if player.Folded {
			continue
		}

		if len(player.Cards) != 2 {
			return ShowdownResult{}, fmt.Errorf("player %s does not have 2 cards", player.ID)
		}

		cards := make([]Card, 0, 7)
		cards = append(cards, player.Cards...)
		cards = append(cards, g.CommunityCards...)

		hand, err := EvaluateSeven(cards)
		if err != nil {
			return ShowdownResult{}, fmt.Errorf("evaluating player %s: %w", player.ID, err)
		}

		result.Hands[i] = hand

		if first {
			best = hand
			result.Winners = []int{i}
			first = false
			continue
		}

		comparison := CompareHands(hand, best)

		if comparison > 0 {
			best = hand
			result.Winners = []int{i}
		} else if comparison == 0 {
			result.Winners = append(result.Winners, i)
		}
	}

	if len(result.Winners) == 0 {
		return ShowdownResult{}, fmt.Errorf("no eligible players for showdown")
	}

	return result, nil
}

func (g *Game) Payout(result ShowdownResult) error {
	if g.State != StateShowdown {
		return fmt.Errorf("game is not at showdown")
	}

	if len(result.Winners) == 0 {
		return fmt.Errorf("no winners")
	}

	if g.Pot <= 0 {
		return fmt.Errorf("pot is empty")
	}

	share := g.Pot / int64(len(result.Winners))
	remainder := g.Pot % int64(len(result.Winners))

	for i, winnerIndex := range result.Winners {
		if winnerIndex < 0 || winnerIndex >= len(g.Players) {
			return fmt.Errorf("invalid winner index %d", winnerIndex)
		}

		amount := share

		if int64(i) < remainder {
			amount++
		}

		g.Players[winnerIndex].Chips += amount
	}

	g.Pot = 0

	return nil
}
