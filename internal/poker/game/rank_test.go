package game

import "testing"

func TestHighCard(t *testing.T) {
	hand, err := EvaluateFive([]Card{
		card(Ace, Spades),
		card(King, Hearts),
		card(Nine, Clubs),
		card(Five, Diamonds),
		card(Two, Spades),
	})

	if err != nil {
		t.Fatal(err)
	}

	if hand.Rank != HighCard {
		t.Fatalf(
			"expected high card, got %v",
			hand.Rank,
		)
	}
}

func TestPair(t *testing.T) {
	hand, err := EvaluateFive([]Card{
		card(Ace, Spades),
		card(Ace, Hearts),
		card(Nine, Clubs),
		card(Five, Diamonds),
		card(Two, Spades),
	})

	if err != nil {
		t.Fatal(err)
	}

	if hand.Rank != Pair {
		t.Fatalf(
			"expected pair, got %v",
			hand.Rank,
		)
	}
}

func TestTwoPair(t *testing.T) {
	hand, err := EvaluateFive([]Card{
		card(Ace, Spades),
		card(Ace, Hearts),
		card(King, Clubs),
		card(King, Diamonds),
		card(Two, Spades),
	})

	if err != nil {
		t.Fatal(err)
	}

	if hand.Rank != TwoPair {
		t.Fatalf(
			"expected two pair, got %v",
			hand.Rank,
		)
	}
}

func TestThreeOfAKind(t *testing.T) {
	hand, err := EvaluateFive([]Card{
		card(Ace, Spades),
		card(Ace, Hearts),
		card(Ace, Clubs),
		card(King, Diamonds),
		card(Two, Spades),
	})

	if err != nil {
		t.Fatal(err)
	}

	if hand.Rank != ThreeOfAKind {
		t.Fatalf(
			"expected three of a kind, got %v",
			hand.Rank,
		)
	}
}

func TestStraight(t *testing.T) {
	hand, err := EvaluateFive([]Card{
		card(Six, Spades),
		card(Five, Hearts),
		card(Four, Clubs),
		card(Three, Diamonds),
		card(Two, Spades),
	})

	if err != nil {
		t.Fatal(err)
	}

	if hand.Rank != Straight {
		t.Fatalf(
			"expected straight, got %v",
			hand.Rank,
		)
	}

	if hand.Tiebreaks[0] != Six {
		t.Fatalf(
			"expected straight high card 6, got %d",
			hand.Tiebreaks[0],
		)
	}
}

func TestAceLowStraight(t *testing.T) {
	hand, err := EvaluateFive([]Card{
		card(Ace, Spades),
		card(Five, Hearts),
		card(Four, Clubs),
		card(Three, Diamonds),
		card(Two, Spades),
	})

	if err != nil {
		t.Fatal(err)
	}

	if hand.Rank != Straight {
		t.Fatalf(
			"expected straight, got %v",
			hand.Rank,
		)
	}

	if hand.Tiebreaks[0] != Five {
		t.Fatalf(
			"expected high card 5, got %d",
			hand.Tiebreaks[0],
		)
	}
}

func TestFlush(t *testing.T) {
	hand, err := EvaluateFive([]Card{
		card(Ace, Spades),
		card(Jack, Spades),
		card(Nine, Spades),
		card(Five, Spades),
		card(Two, Spades),
	})

	if err != nil {
		t.Fatal(err)
	}

	if hand.Rank != Flush {
		t.Fatalf(
			"expected flush, got %v",
			hand.Rank,
		)
	}
}

func TestFullHouse(t *testing.T) {
	hand, err := EvaluateFive([]Card{
		card(Ace, Spades),
		card(Ace, Hearts),
		card(Ace, Clubs),
		card(King, Diamonds),
		card(King, Spades),
	})

	if err != nil {
		t.Fatal(err)
	}

	if hand.Rank != FullHouse {
		t.Fatalf(
			"expected full house, got %v",
			hand.Rank,
		)
	}
}

func TestFourOfAKind(t *testing.T) {
	hand, err := EvaluateFive([]Card{
		card(Ace, Spades),
		card(Ace, Hearts),
		card(Ace, Clubs),
		card(Ace, Diamonds),
		card(King, Spades),
	})

	if err != nil {
		t.Fatal(err)
	}

	if hand.Rank != FourOfAKind {
		t.Fatalf(
			"expected four of a kind, got %v",
			hand.Rank,
		)
	}
}

func TestStraightFlush(t *testing.T) {
	hand, err := EvaluateFive([]Card{
		card(Nine, Spades),
		card(Eight, Spades),
		card(Seven, Spades),
		card(Six, Spades),
		card(Five, Spades),
	})

	if err != nil {
		t.Fatal(err)
	}

	if hand.Rank != StraightFlush {
		t.Fatalf(
			"expected straight flush, got %v",
			hand.Rank,
		)
	}
}

func TestCompareHands(t *testing.T) {
	aces := HandValue{
		Rank: Pair,
		Tiebreaks: []Rank{
			Ace,
			King,
			Nine,
			Five,
		},
	}

	kings := HandValue{
		Rank: Pair,
		Tiebreaks: []Rank{
			King,
			Ace,
			Nine,
			Five,
		},
	}

	if CompareHands(aces, kings) <= 0 {
		t.Fatal("expected pair of aces to win")
	}
}

func TestKickerComparison(t *testing.T) {
	aceKicker := HandValue{
		Rank: Pair,
		Tiebreaks: []Rank{
			King,
			Ace,
			Nine,
			Five,
		},
	}

	queenKicker := HandValue{
		Rank: Pair,
		Tiebreaks: []Rank{
			King,
			Queen,
			Nine,
			Five,
		},
	}

	if CompareHands(aceKicker, queenKicker) <= 0 {
		t.Fatal("expected ace kicker to win")
	}
}

func TestEvaluateSeven(t *testing.T) {
	hand, err := EvaluateSeven([]Card{
		card(Ace, Spades),
		card(King, Spades),

		card(Queen, Spades),
		card(Jack, Spades),
		card(Ten, Spades),
		card(Two, Diamonds),
		card(Three, Clubs),
	})

	if err != nil {
		t.Fatal(err)
	}

	if hand.Rank != StraightFlush {
		t.Fatalf(
			"expected straight flush, got %v",
			hand.Rank,
		)
	}

	if hand.Tiebreaks[0] != Ace {
		t.Fatalf(
			"expected ace-high straight flush",
		)
	}
}

func TestEvaluateFiveInvalidCardCount(t *testing.T) {
	_, err := EvaluateFive([]Card{
		card(Ace, Spades),
		card(King, Spades),
	})

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestEvaluateSevenInvalidCardCount(t *testing.T) {
	_, err := EvaluateSeven([]Card{
		card(Ace, Spades),
		card(King, Spades),
	})

	if err == nil {
		t.Fatal("expected error")
	}
}
