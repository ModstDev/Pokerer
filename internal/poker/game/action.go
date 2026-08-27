package game

import "fmt"

type ActionType string

const (
	ActionFold  ActionType = "fold"
	ActionCheck ActionType = "check"
	ActionCall  ActionType = "call"
	ActionBet   ActionType = "bet"
	ActionRaise ActionType = "raise"
	ActionAllIn ActionType = "all_in"
)

type Action struct {
	Type   ActionType
	Amount int64
}

func (g *Game) ApplyAction(action Action) error {
	if g.State == StateWaiting {
		return fmt.Errorf("hand has not started")
	}

	if g.State == StateShowdown {
		return fmt.Errorf("hand is already over")
	}

	if g.CurrentPlayer < 0 || g.CurrentPlayer >= len(g.Players) {
		return fmt.Errorf("invalid current player")
	}

	player := &g.Players[g.CurrentPlayer]

	if player.Folded {
		return fmt.Errorf("player has folded")
	}

	if player.AllIn {
		return fmt.Errorf("player is all-in")
	}

	switch action.Type {
	case ActionFold:
		return g.fold(player)

	case ActionCheck:
		return g.check(player)

	case ActionCall:
		return g.call(player)

	case ActionBet:
		return g.bet(player, action.Amount)

	case ActionRaise:
		return g.raise(player, action.Amount)

	case ActionAllIn:
		return g.allIn(player)

	default:
		return fmt.Errorf("unknown action: %q", action.Type)
	}
}

func (g *Game) fold(player *Player) error {
	player.Folded = true

	return g.advanceAfterAction()
}

func (g *Game) check(player *Player) error {
	if player.Bet != g.CurrentBet {
		return fmt.Errorf("cannot check, outstanding bet exists")
	}

	return g.advanceAfterAction()
}

func (g *Game) call(player *Player) error {
	amount := g.CurrentBet - player.Bet

	if amount <= 0 {
		return fmt.Errorf("nothing to call")
	}

	if amount >= player.Chips {
		return g.allIn(player)
	}

	player.Chips -= amount
	player.Bet += amount
	g.Pot += amount

	return g.advanceAfterAction()
}

func (g *Game) bet(player *Player, amount int64) error {
	if g.CurrentBet != 0 {
		return fmt.Errorf("cannot bet, an existing bet must be raised")
	}

	if amount <= 0 {
		return fmt.Errorf("bet amount must be positive")
	}

	if amount > player.Chips {
		return fmt.Errorf("bet exceeds available chips")
	}

	player.Chips -= amount
	player.Bet += amount

	g.CurrentBet = player.Bet
	g.Pot += amount

	if player.Chips == 0 {
		player.AllIn = true
	}

	return g.advanceAfterAction()
}

func (g *Game) raise(player *Player, amount int64) error {
	if g.CurrentBet <= 0 {
		return fmt.Errorf("cannot raise without an existing bet")
	}

	if amount <= g.CurrentBet {
		return fmt.Errorf("raise must exceed current bet")
	}

	additional := amount - player.Bet

	if additional > player.Chips {
		return fmt.Errorf("raise exceeds available chips")
	}

	player.Chips -= additional
	player.Bet = amount

	g.CurrentBet = amount
	g.Pot += additional

	if player.Chips == 0 {
		player.AllIn = true
	}

	return g.advanceAfterAction()
}

func (g *Game) allIn(player *Player) error {
	if player.Chips <= 0 {
		return fmt.Errorf("player has no chips")
	}

	amount := player.Chips

	player.Chips = 0
	player.Bet += amount
	player.AllIn = true

	if player.Bet > g.CurrentBet {
		g.CurrentBet = player.Bet
	}

	g.Pot += amount

	return g.advanceAfterAction()
}

func (g *Game) advanceAfterAction() error {
	activePlayers := 0

	for _, player := range g.Players {
		if !player.Folded {
			activePlayers++
		}
	}

	if activePlayers <= 1 {
		g.State = StateShowdown
		return nil
	}

	return g.advancePlayer()
}
