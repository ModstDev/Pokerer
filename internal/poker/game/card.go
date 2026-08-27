package game

type Suit uint8

const (
	Clubs Suit = iota
	Diamonds
	Hearts
	Spades
)

type Rank uint8

const (
	Two Rank = iota + 2
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
)

type Card struct {
	Suit Suit
	Rank Rank
}

func card(rank Rank, suit Suit) Card {
	return Card{
		Rank: rank,
		Suit: suit,
	}
}
