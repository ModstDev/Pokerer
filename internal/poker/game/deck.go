package game

import "math/rand/v2"

type Deck struct {
	cards []Card
}

func NewDeck() *Deck {
	cards := make([]Card, 0, 52)

	for suit := Clubs; suit <= Spades; suit++ {
		for rank := Two; rank <= Ace; rank++ {
			cards = append(cards, Card{
				Suit: suit,
				Rank: rank,
			})
		}
	}

	return &Deck{
		cards: cards,
	}
}

func (d *Deck) Shuffle(r *rand.Rand) {
	r.Shuffle(len(d.cards), func(i, j int) {
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	})
}

func (d *Deck) Draw() (Card, bool) {
	if len(d.cards) == 0 {
		return Card{}, false
	}

	card := d.cards[len(d.cards)-1]
	d.cards = d.cards[:len(d.cards)-1]

	return card, true
}
