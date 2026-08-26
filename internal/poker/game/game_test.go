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
