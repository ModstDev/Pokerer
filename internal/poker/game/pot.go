package game

import (
	"fmt"
	"sort"
)

type Pot struct {
	Amount          int64
	EligiblePlayers []int
}

type PotResult struct {
	Pot     Pot
	Winners []int
}

type PayoutResult struct {
	Pots []PotResult
}

func (g *Game) BuildPots() []Pot {
	levels := make([]int64, 0)

	for _, player := range g.Players {
		if player.TotalContribution <= 0 {
			continue
		}

		levels = append(levels, player.TotalContribution)
	}

	sort.Slice(levels, func(i, j int) bool {
		return levels[i] < levels[j]
	})

	levels = uniqueInt64(levels)

	pots := make([]Pot, 0)

	var previousLevel int64

	for _, level := range levels {
		contribution := level - previousLevel

		if contribution <= 0 {
			continue
		}
		contributors := 0
		eligible := make([]int, 0)

		for i, player := range g.Players {
			if player.TotalContribution >= level {
				contributors++
			}

			if player.TotalContribution >= level && !player.Folded {
				eligible = append(eligible, i)
			}
		}

		amount := contribution * int64(contributors)

		if amount > 0 && len(eligible) >= 2 {
			pots = append(pots, Pot{
				Amount:          amount,
				EligiblePlayers: eligible,
			})
		}

		previousLevel = level
	}

	return pots
}

func uniqueInt64(values []int64) []int64 {
	if len(values) == 0 {
		return nil
	}

	result := []int64{values[0]}

	for i := 1; i < len(values); i++ {
		if values[i] != values[i-1] {
			result = append(result, values[i])
		}
	}

	return result
}

func (g *Game) resolvePot(pot Pot, hands map[int]HandValue) (PotResult, error) {
	if len(pot.EligiblePlayers) == 0 {
		return PotResult{}, fmt.Errorf("pot has no eligible players")
	}

	var best HandValue
	winners := make([]int, 0)

	for _, playerIndex := range pot.EligiblePlayers {
		hand, ok := hands[playerIndex]
		if !ok {
			return PotResult{}, fmt.Errorf("missing hand for player %d", playerIndex)
		}

		if len(winners) == 0 {
			best = hand
			winners = append(winners, playerIndex)
			continue
		}

		comparison := CompareHands(hand, best)

		switch {
		case comparison > 0:
			best = hand
			winners = []int{playerIndex}

		case comparison == 0:
			winners = append(winners, playerIndex)
		}
	}

	return PotResult{
		Pot:     pot,
		Winners: winners,
	}, nil
}

func (g *Game) BuildPayout() (PayoutResult, error) {
	if g.State != StateShowdown {
		return PayoutResult{}, fmt.Errorf("game is not at showdown")
	}

	pots := g.BuildPots()

	if len(pots) == 0 {
		return PayoutResult{}, fmt.Errorf("no pots to distribute")
	}

	showdown, err := g.EvaluateShowdown()
	if err != nil {
		return PayoutResult{}, err
	}

	result := PayoutResult{
		Pots: make([]PotResult, 0, len(pots)),
	}

	for _, pot := range pots {
		potResult, err := g.resolvePot(pot, showdown.Hands)
		if err != nil {
			return PayoutResult{}, err
		}

		result.Pots = append(result.Pots, potResult)
	}

	return result, nil
}

func (g *Game) ApplyPayout(result PayoutResult) error {
	if g.State != StateShowdown {
		return fmt.Errorf("game is not at showdown")
	}

	for _, potResult := range result.Pots {
		if err := distributePot(&g.Players, potResult); err != nil {
			return err
		}
	}

	g.Pot = 0

	return nil
}

func distributePot(players *[]Player, potResult PotResult) error {
	if len(potResult.Winners) == 0 {
		return fmt.Errorf("pot has no winners")
	}

	share := potResult.Pot.Amount /
		int64(len(potResult.Winners))

	remainder := potResult.Pot.Amount %
		int64(len(potResult.Winners))

	for i, winnerIndex := range potResult.Winners {
		if winnerIndex < 0 ||
			winnerIndex >= len(*players) {
			return fmt.Errorf("invalid winner index %d", winnerIndex)
		}

		amount := share

		if int64(i) < remainder {
			amount++
		}

		(*players)[winnerIndex].Chips += amount
	}

	return nil
}
