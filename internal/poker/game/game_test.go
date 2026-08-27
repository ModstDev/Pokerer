package game

import (
	"math/rand/v2"
	"testing"
)

func TestStartHand(t *testing.T) {
	game := NewGame()

	err := game.AddPlayer(Player{
		ID:    "player-1",
		Seat:  0,
		Chips: 1000,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = game.AddPlayer(Player{
		ID:    "player-2",
		Seat:  1,
		Chips: 1000,
	})
	if err != nil {
		t.Fatal(err)
	}

	rng := rand.New(rand.NewPCG(42, 0))

	err = game.StartHand(rng)
	if err != nil {
		t.Fatal(err)
	}

	if game.State != StatePreFlop {
		t.Fatalf(
			"expected state %q, got %q",
			StatePreFlop,
			game.State,
		)
	}

	for _, player := range game.Players {
		if len(player.Cards) != 2 {
			t.Fatalf(
				"expected 2 cards, got %d",
				len(player.Cards),
			)
		}
	}
}

func TestDeckContains52Cards(t *testing.T) {
	deck := NewDeck()

	seen := make(map[Card]bool)

	for i := 0; i < 52; i++ {
		card, ok := deck.Draw()

		if !ok {
			t.Fatalf("deck ran out after %d cards", i)
		}

		if seen[card] {
			t.Fatalf("duplicate card: %+v", card)
		}

		seen[card] = true
	}

	_, ok := deck.Draw()

	if ok {
		t.Fatal("expected deck to be empty")
	}
}

func TestAdvanceRounds(t *testing.T) {
	game := NewGame()

	if err := game.AddPlayer(Player{
		ID:    "player-1",
		Seat:  0,
		Chips: 1000,
	}); err != nil {
		t.Fatal(err)
	}

	if err := game.AddPlayer(Player{
		ID:    "player-2",
		Seat:  1,
		Chips: 1000,
	}); err != nil {
		t.Fatal(err)
	}

	rng := rand.New(rand.NewPCG(42, 0))

	if err := game.StartHand(rng); err != nil {
		t.Fatal(err)
	}

	if game.State != StatePreFlop {
		t.Fatalf("expected pre-flop, got %s", game.State)
	}

	if len(game.CommunityCards) != 0 {
		t.Fatalf("expected 0 community cards, got %d",
			len(game.CommunityCards))
	}

	// Pre-flop -> Flop.
	if err := game.AdvanceRound(); err != nil {
		t.Fatal(err)
	}

	if game.State != StateFlop {
		t.Fatalf("expected flop, got %s", game.State)
	}

	if len(game.CommunityCards) != 3 {
		t.Fatalf(
			"expected 3 community cards, got %d",
			len(game.CommunityCards),
		)
	}

	// Flop -> Turn.
	if err := game.AdvanceRound(); err != nil {
		t.Fatal(err)
	}

	if game.State != StateTurn {
		t.Fatalf("expected turn, got %s", game.State)
	}

	if len(game.CommunityCards) != 4 {
		t.Fatalf(
			"expected 4 community cards, got %d",
			len(game.CommunityCards),
		)
	}

	// Turn -> River.
	if err := game.AdvanceRound(); err != nil {
		t.Fatal(err)
	}

	if game.State != StateRiver {
		t.Fatalf("expected river, got %s", game.State)
	}

	if len(game.CommunityCards) != 5 {
		t.Fatalf(
			"expected 5 community cards, got %d",
			len(game.CommunityCards),
		)
	}

	// River -> Showdown.
	if err := game.AdvanceRound(); err != nil {
		t.Fatal(err)
	}

	if game.State != StateShowdown {
		t.Fatalf(
			"expected showdown, got %s",
			game.State,
		)
	}
}

func TestCannotAdvanceWaitingGame(t *testing.T) {
	game := NewGame()

	if err := game.AdvanceRound(); err == nil {
		t.Fatal("expected error when advancing waiting game")
	}
}

func TestCannotAdvanceShowdown(t *testing.T) {
	game := NewGame()

	game.State = StateShowdown

	if err := game.AdvanceRound(); err == nil {
		t.Fatal("expected error when advancing showdown")
	}
}

func TestBurn(t *testing.T) {
	deck := NewDeck()

	for i := 0; i < 5; i++ {
		if !deck.Burn() {
			t.Fatalf("expected burn %d to succeed", i)
		}
	}

	for i := 0; i < 47; i++ {
		if _, ok := deck.Draw(); !ok {
			t.Fatalf("expected draw %d to succeed", i)
		}
	}

	if deck.Burn() {
		t.Fatal("expected burn to fail on empty deck")
	}
}

func TestCheck(t *testing.T) {
	game := NewGame()

	game.Players = []Player{
		{
			ID:    "player-1",
			Seat:  0,
			Chips: 1000,
		},
		{
			ID:    "player-2",
			Seat:  1,
			Chips: 1000,
		},
	}

	game.State = StatePreFlop
	game.CurrentPlayer = 0
	game.CurrentBet = 0

	err := game.ApplyAction(Action{
		Type: ActionCheck,
	})

	if err != nil {
		t.Fatal(err)
	}

	if game.CurrentPlayer != 1 {
		t.Fatalf(
			"expected player 1, got %d",
			game.CurrentPlayer,
		)
	}
}

func TestCannotCheckAgainstBet(t *testing.T) {
	game := NewGame()

	game.Players = []Player{
		{
			ID:    "player-1",
			Seat:  0,
			Chips: 1000,
			Bet:   0,
		},
	}

	game.State = StatePreFlop
	game.CurrentPlayer = 0
	game.CurrentBet = 100

	err := game.ApplyAction(Action{
		Type: ActionCheck,
	})

	if err == nil {
		t.Fatal("expected check to fail")
	}
}

func TestCall(t *testing.T) {
	game := NewGame()

	game.Players = []Player{
		{
			ID:    "player-1",
			Seat:  0,
			Chips: 1000,
			Bet:   0,
		},
		{
			ID:    "player-2",
			Seat:  1,
			Chips: 1000,
			Bet:   100,
		},
	}

	game.State = StatePreFlop
	game.CurrentPlayer = 0
	game.CurrentBet = 100

	err := game.ApplyAction(Action{
		Type: ActionCall,
	})

	if err != nil {
		t.Fatal(err)
	}

	if game.Players[0].Chips != 900 {
		t.Fatalf(
			"expected 900 chips, got %d",
			game.Players[0].Chips,
		)
	}

	if game.Players[0].Bet != 100 {
		t.Fatalf(
			"expected bet 100, got %d",
			game.Players[0].Bet,
		)
	}

	if game.Pot != 100 {
		t.Fatalf(
			"expected pot 100, got %d",
			game.Pot,
		)
	}
}
