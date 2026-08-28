package game

import "testing"

func TestShowdown(t *testing.T) {
	game := NewGame(GameConfig{
		SmallBlind: 10,
		BigBlind:   20,
	})

	game.Players = []Player{
		{
			ID: "alice",
			Cards: []Card{
				card(Ace, Spades),
				card(Ace, Hearts),
			},
			Chips: 800,
		},
		{
			ID: "bob",
			Cards: []Card{
				card(King, Spades),
				card(King, Hearts),
			},
			Chips: 800,
		},
	}

	game.CommunityCards = []Card{
		card(Two, Clubs),
		card(Seven, Diamonds),
		card(Nine, Spades),
		card(Five, Hearts),
		card(Three, Clubs),
	}

	game.Pot = 400
	game.State = StateShowdown

	result, err := game.EvaluateShowdown()
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Winners) != 1 {
		t.Fatalf(
			"expected one winner, got %d",
			len(result.Winners),
		)
	}

	if result.Winners[0] != 0 {
		t.Fatalf(
			"expected Alice to win, got player %d",
			result.Winners[0],
		)
	}

	if err := game.Payout(result); err != nil {
		t.Fatal(err)
	}

	if game.Players[0].Chips != 1200 {
		t.Fatalf(
			"expected Alice to have 1200 chips, got %d",
			game.Players[0].Chips,
		)
	}

	if game.Pot != 0 {
		t.Fatalf(
			"expected pot to be empty, got %d",
			game.Pot,
		)
	}
}

func TestShowdownTie(t *testing.T) {
	game := NewGame(GameConfig{
		SmallBlind: 10,
		BigBlind:   20,
	})

	game.Players = []Player{
		{
			ID: "alice",
			Cards: []Card{
				card(Two, Clubs),
				card(Three, Diamonds),
			},
			Chips: 800,
		},
		{
			ID: "bob",
			Cards: []Card{
				card(Four, Clubs),
				card(Five, Diamonds),
			},
			Chips: 800,
		},
	}

	game.CommunityCards = []Card{
		card(Ace, Spades),
		card(King, Hearts),
		card(Queen, Clubs),
		card(Jack, Diamonds),
		card(Ten, Spades),
	}

	game.Pot = 401
	game.State = StateShowdown

	result, err := game.EvaluateShowdown()
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Winners) != 2 {
		t.Fatalf(
			"expected two winners, got %d",
			len(result.Winners),
		)
	}

	if err := game.Payout(result); err != nil {
		t.Fatal(err)
	}

	if game.Players[0].Chips != 1001 {
		t.Fatalf(
			"expected Alice to have 1001 chips, got %d",
			game.Players[0].Chips,
		)
	}

	if game.Players[1].Chips != 1000 {
		t.Fatalf(
			"expected Bob to have 1000 chips, got %d",
			game.Players[1].Chips,
		)
	}
}
