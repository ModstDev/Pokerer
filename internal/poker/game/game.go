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
	StateFinished State = "finished"
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

	if g.dealerIndex() == -1 {
		g.DealerPosition = g.Players[0].Seat
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
		g.Players[i].TotalContribution = 0
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
	dealer := g.dealerIndex()

	if dealer == -1 {
		return -1, -1
	}

	playerCount := len(g.Players)

	if playerCount == 2 {
		return dealer,
			nextPosition(dealer, 1, playerCount)
	}

	smallBlind := nextPosition(dealer, 1, playerCount)

	bigBlind := nextPosition(dealer, 2, playerCount)

	return smallBlind, bigBlind
}

func nextPosition(position, offset, playerCount int) int {
	return (position + offset) % playerCount
}

func (g *Game) postBlind(playerIndex int, amount int64) int64 {
	player := &g.Players[playerIndex]

	if amount > player.Chips {
		amount = player.Chips
	}

	g.contribute(player, amount)

	if player.Chips == 0 {
		player.AllIn = true
	}

	return amount
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

func (g *Game) resetHand() {
	for i := range g.Players {
		g.Players[i].Cards = nil
		g.Players[i].Bet = 0
		g.Players[i].Acted = false
		g.Players[i].TotalContribution = 0
		g.Players[i].Folded = false
		g.Players[i].AllIn = false
	}

	g.CommunityCards = g.CommunityCards[:0]
	g.Pot = 0
	g.CurrentBet = 0
	g.MinRaise = g.Config.BigBlind
}

func (g *Game) removeBustedPlayers() {
	players := g.Players

	remaining := make([]Player, 0, len(players))

	for _, player := range players {
		if player.Chips > 0 {
			remaining = append(remaining, player)
		}
	}

	g.Players = remaining
}

func (g *Game) dealerIndex() int {
	for i, player := range g.Players {
		if player.Seat == g.DealerPosition {
			return i
		}
	}

	return -1
}

func (g *Game) FinishHand() error {
	if g.State != StateShowdown {
		return fmt.Errorf("game is not at showdown")
	}

	payout, err := g.BuildPayout()
	if err != nil {
		return fmt.Errorf("building payout: %w", err)
	}

	if err := g.ApplyPayout(payout); err != nil {
		return fmt.Errorf("applying payout: %w", err)
	}

	g.resetHand()

	g.removeBustedPlayers()

	if len(g.Players) >= 2 {
		g.rotateDealer()
		g.State = StateWaiting
	} else {
		g.State = StateFinished
	}

	return nil
}
