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
			Chips:             800,
			TotalContribution: 200,
		},
		{
			ID: "bob",
			Cards: []Card{
				card(King, Spades),
				card(King, Hearts),
			},
			Chips:             800,
			TotalContribution: 200,
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

	payout, err := game.BuildPayout()
	if err != nil {
		t.Fatal(err)
	}

	if len(payout.Pots) != 1 {
		t.Fatalf(
			"expected 1 pot, got %d",
			len(payout.Pots),
		)
	}

	if len(payout.Pots[0].Winners) != 1 {
		t.Fatalf("expected one winner")
	}

	if payout.Pots[0].Winners[0] != 0 {
		t.Fatalf("expected Alice to win")
	}

	if err := game.ApplyPayout(payout); err != nil {
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
			Chips:             800,
			TotalContribution: 200,
		},
		{
			ID: "bob",
			Cards: []Card{
				card(Four, Clubs),
				card(Five, Diamonds),
			},
			Chips:             800,
			TotalContribution: 200,
		},
	}

	game.CommunityCards = []Card{
		card(Ace, Spades),
		card(King, Hearts),
		card(Queen, Clubs),
		card(Jack, Diamonds),
		card(Ten, Spades),
	}

	game.Pot = 400
	game.State = StateShowdown

	payout, err := game.BuildPayout()
	if err != nil {
		t.Fatal(err)
	}

	if len(payout.Pots) != 1 {
		t.Fatalf(
			"expected 1 pot, got %d",
			len(payout.Pots),
		)
	}

	if len(payout.Pots[0].Winners) != 2 {
		t.Fatalf(
			"expected two winners, got %d",
			len(payout.Pots[0].Winners),
		)
	}

	if err := game.ApplyPayout(payout); err != nil {
		t.Fatal(err)
	}

	if game.Players[0].Chips != 1000 {
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

func TestBuildPots(t *testing.T) {
	game := NewGame(GameConfig{
		SmallBlind: 10,
		BigBlind:   20,
	})

	game.Players = []Player{
		{
			ID:                "alice",
			TotalContribution: 100,
		},
		{
			ID:                "bob",
			TotalContribution: 300,
		},
		{
			ID:                "charlie",
			TotalContribution: 500,
		},
	}

	pots := game.BuildPots()

	if len(pots) != 2 {
		t.Fatalf(
			"expected 2 pots, got %d",
			len(pots),
		)
	}

	if pots[0].Amount != 300 {
		t.Fatalf(
			"expected main pot 300, got %d",
			pots[0].Amount,
		)
	}

	if len(pots[0].EligiblePlayers) != 3 {
		t.Fatalf(
			"expected 3 eligible players, got %d",
			len(pots[0].EligiblePlayers),
		)
	}

	if pots[1].Amount != 400 {
		t.Fatalf(
			"expected side pot 400, got %d",
			pots[1].Amount,
		)
	}

	if len(pots[1].EligiblePlayers) != 2 {
		t.Fatalf(
			"expected 2 eligible players, got %d",
			len(pots[1].EligiblePlayers),
		)
	}
}

func TestSidePotPayout(t *testing.T) {
	game := NewGame(GameConfig{
		SmallBlind: 10,
		BigBlind:   20,
	})

	game.Players = []Player{
		{
			ID:                "alice",
			Chips:             0,
			TotalContribution: 100,
			Cards: []Card{
				card(Ace, Spades),
				card(Ace, Hearts),
			},
		},
		{
			ID:                "bob",
			Chips:             0,
			TotalContribution: 300,
			Cards: []Card{
				card(King, Spades),
				card(King, Hearts),
			},
		},
		{
			ID:                "charlie",
			Chips:             0,
			TotalContribution: 500,
			Cards: []Card{
				card(Queen, Spades),
				card(Queen, Hearts),
			},
		},
	}

	game.CommunityCards = []Card{
		card(Two, Clubs),
		card(Seven, Diamonds),
		card(Nine, Spades),
		card(Five, Hearts),
		card(Three, Clubs),
	}

	game.Pot = 900
	game.State = StateShowdown

	result, err := game.BuildPayout()
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Pots) != 2 {
		t.Fatalf(
			"expected 2 pots, got %d",
			len(result.Pots),
		)
	}

	if result.Pots[0].Winners[0] != 0 {
		t.Fatalf(
			"expected Alice to win main pot",
		)
	}

	if result.Pots[1].Winners[0] != 1 {
		t.Fatalf(
			"expected Bob to win side pot",
		)
	}

	if err := game.ApplyPayout(result); err != nil {
		t.Fatal(err)
	}

	if game.Players[0].Chips != 300 {
		t.Fatalf(
			"expected Alice to receive 300, got %d",
			game.Players[0].Chips,
		)
	}

	if game.Players[1].Chips != 400 {
		t.Fatalf(
			"expected Bob to receive 400, got %d",
			game.Players[1].Chips,
		)
	}

	if game.Players[2].Chips != 0 {
		t.Fatalf(
			"expected Charlie to receive 0, got %d",
			game.Players[2].Chips,
		)
	}
}
