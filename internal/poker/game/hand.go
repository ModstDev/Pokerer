package game

import (
	"fmt"
	"sort"
)

type HandRank uint8

const (
	HighCard HandRank = iota
	Pair
	TwoPair
	ThreeOfAKind
	Straight
	Flush
	FullHouse
	FourOfAKind
	StraightFlush
)

type HandValue struct {
	Rank      HandRank
	Tiebreaks []Rank
}

type rankGroup struct {
	rank  Rank
	count int
}

func CompareHands(a, b HandValue) int {
	if a.Rank > b.Rank {
		return 1
	}

	if a.Rank < b.Rank {
		return -1
	}

	length := len(a.Tiebreaks)

	if len(b.Tiebreaks) < length {
		length = len(b.Tiebreaks)
	}

	for i := 0; i < length; i++ {
		if a.Tiebreaks[i] > b.Tiebreaks[i] {
			return 1
		}

		if a.Tiebreaks[i] < b.Tiebreaks[i] {
			return -1
		}
	}

	return 0
}

func rankCounts(cards []Card) map[Rank]int {
	counts := make(map[Rank]int)

	for _, card := range cards {
		counts[card.Rank]++
	}

	return counts
}

func sortedRanksDescending(cards []Card) []Rank {
	ranks := make([]Rank, 0, len(cards))

	for _, card := range cards {
		ranks = append(ranks, card.Rank)
	}

	sort.Slice(ranks, func(i, j int) bool {
		return ranks[i] > ranks[j]
	})

	return ranks
}

func isFlush(cards []Card) bool {
	for i := 1; i < len(cards); i++ {
		if cards[i].Suit != cards[0].Suit {
			return false
		}
	}

	return true
}

func straitghtHighCard(cards []Card) (Rank, bool) {
	ranks := sortedRanksDescending(cards)

	unique := make([]Rank, 0, len(ranks))

	for _, rank := range ranks {
		if len(unique) == 0 || unique[len(unique)-1] != rank {
			unique = append(unique, rank)
		}
	}

	if len(unique) != 5 {
		return 0, false
	}

	if unique[0] == Ace &&
		unique[1] == Five &&
		unique[2] == Four &&
		unique[3] == Three &&
		unique[4] == Two {
		return Five, true
	}

	for i := 1; i < len(unique); i++ {
		if unique[i] != unique[i-1]-1 {
			return 0, false
		}
	}

	return unique[0], true
}

func rankGroups(cards []Card) []rankGroup {
	counts := rankCounts(cards)

	groups := make([]rankGroup, 0, len(counts))

	for rank, count := range counts {
		groups = append(groups, rankGroup{
			rank:  rank,
			count: count,
		})
	}

	sort.Slice(
		groups,
		func(i, j int) bool {
			if groups[i].count != groups[j].count {
				return groups[i].count >
					groups[j].count
			}

			return groups[i].rank >
				groups[j].rank
		},
	)

	return groups
}

func EvaluateFive(cards []Card) (HandValue, error) {
	if len(cards) != 5 {
		return HandValue{}, fmt.Errorf("EvaluateFive requires exactly 5 cards")
	}

	groups := rankGroups(cards)
	flush := isFlush(cards)

	if straightHigh, ok := straitghtHighCard(cards); ok {
		if flush {
			return HandValue{
				Rank: StraightFlush,
				Tiebreaks: []Rank{
					straightHigh,
				},
			}, nil
		}
	}

	// Four of a kind.
	if groups[0].count == 4 {
		return HandValue{
			Rank: FourOfAKind,
			Tiebreaks: []Rank{
				groups[0].rank,
				groups[1].rank,
			},
		}, nil
	}

	// Full house.
	if groups[0].count == 3 && groups[1].count == 2 {
		return HandValue{
			Rank: FullHouse,
			Tiebreaks: []Rank{
				groups[0].rank,
				groups[1].rank,
			},
		}, nil
	}

	// Flush.
	if flush {
		return HandValue{
			Rank:      Flush,
			Tiebreaks: sortedRanksDescending(cards),
		}, nil
	}

	// Straight.
	if straightHigh, ok := straitghtHighCard(cards); ok {
		return HandValue{
			Rank: Straight,
			Tiebreaks: []Rank{
				straightHigh,
			},
		}, nil
	}

	// Three of a kind.
	if groups[0].count == 3 {
		kickers := make([]Rank, 0, 2)

		for _, group := range groups[1:] {
			kickers = append(kickers, group.rank)
		}

		sort.Slice(kickers, func(i, j int) bool {
			return kickers[i] > kickers[j]
		},
		)

		return HandValue{
			Rank: ThreeOfAKind,
			Tiebreaks: append([]Rank{groups[0].rank},
				kickers...,
			),
		}, nil
	}

	// Two pair.
	if groups[0].count == 2 && groups[1].count == 2 {
		highPair := groups[0].rank
		lowPair := groups[1].rank

		kicker := groups[2].rank

		return HandValue{
			Rank: TwoPair,
			Tiebreaks: []Rank{
				highPair,
				lowPair,
				kicker,
			},
		}, nil
	}

	// One pair.
	if groups[0].count == 2 {
		kickers := make([]Rank, 0, 3)

		for _, group := range groups[1:] {
			kickers = append(kickers, group.rank)
		}

		sort.Slice(kickers, func(i int, j int) bool {
			return kickers[i] > kickers[j]
		},
		)

		return HandValue{
			Rank: Pair,
			Tiebreaks: append(
				[]Rank{groups[0].rank},
				kickers...,
			),
		}, nil
	}

	// High card.
	return HandValue{
		Rank:      HighCard,
		Tiebreaks: sortedRanksDescending(cards),
	}, nil
}

func EvaluateSeven(cards []Card) (HandValue, error) {
	if len(cards) != 7 {
		return HandValue{}, fmt.Errorf("EvaluateSeven requires exactly 7 cards")
	}

	var best HandValue
	first := true

	for a := 0; a < 3; a++ {
		for b := a + 1; b < 4; b++ {
			for c := b + 1; c < 5; c++ {
				for d := c + 1; d < 6; d++ {
					for e := d + 1; e < 7; e++ {
						hand, err := EvaluateFive([]Card{
							cards[a],
							cards[b],
							cards[c],
							cards[d],
							cards[e],
						})

						if err != nil {
							return HandValue{}, fmt.Errorf("evaluating five-card hand: %w", err)
						}

						if first ||
							CompareHands(hand, best) > 0 {
							best = hand
							first = false
						}
					}
				}
			}
		}
	}

	return best, nil
}

func (r HandRank) String() string {
	switch r {
	case HighCard:
		return "high card"
	case Pair:
		return "pair"
	case TwoPair:
		return "two pair"
	case ThreeOfAKind:
		return "three of a kind"
	case Straight:
		return "straight"
	case Flush:
		return "flush"
	case FullHouse:
		return "full house"
	case FourOfAKind:
		return "four of a kind"
	case StraightFlush:
		return "straight flush"
	default:
		return "unknown"
	}
}
